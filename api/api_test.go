package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"

	model "github.com/briansan/ManageMeServer/model/schema"
	"github.com/briansan/ManageMeServer/model/store"
)

type APITestSuite struct {
	suite.Suite
	e *echo.Echo
}

func (suite *APITestSuite) SetupTest() {
	os.Setenv("MANAGEME_MONGO_DATABASE", "test")
	os.Setenv("MANAGEME_SECRET", "test_secret")
	os.Setenv("TESTING", "true")

	store.InitMongoSession()
	store.Nuke()
	suite.e = New()
}

func (suite *APITestSuite) Test001_NormalUsage() {
	// 0a. GET /api/service/ping
	code, pong := suite.request("GET", "/api/service/ping", "", nil, nil)
	suite.Equal(http.StatusOK, code)
	suite.Equal("pong", pong)

	// 0b. GET /api/login (as admin)
	var token map[string]string
	code, _ = suite.request(
		"GET", "/api/login",
		basicAuthString("boss", "test_secret"),
		nil, &token,
	)
	suite.Equal(http.StatusOK, code)

	adminSession, ok := token["session"]
	suite.True(ok)

	// 1. POST /api/users
	username, password, email := "foo", "bar", "foo@bar.com"
	user := &model.User{
		Username: &username,
		Password: &password,
		Email:    &email,
	}

	secureUser := &model.UserSecure{}
	code, _ = suite.request("POST", "/api/users", "", user, secureUser)
	suite.Equal(http.StatusCreated, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)
	suite.Nil(user.PreferredHours)

	// save user id
	uid := secureUser.ID.Hex()

	// 2. GET /api/login
	code, _ = suite.request(
		"GET", "/api/login",
		basicAuthString(username, password),
		nil, &token,
	)
	suite.Equal(http.StatusOK, code)

	session, ok := token["session"]
	suite.True(ok)
	jwtAuth := jwtAuthString(session)

	// 2a. GET /api/users/{userID} (as user)
	secureUser = &model.UserSecure{}
	code, _ = suite.request("GET", "/api/users/"+uid, jwtAuth, nil, secureUser)
	suite.Equal(http.StatusOK, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)
	suite.Nil(user.PreferredHours)

	// 2b. GET /api/users/{username} (as user)
	code, _ = suite.request("GET", "/api/users/"+username, jwtAuth, nil, secureUser)
	suite.Equal(http.StatusOK, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)
	suite.Nil(user.PreferredHours)

	// 2c. GET /api/users/{userID} (as admin)
	code, _ = suite.request("GET", "/api/users/"+uid, jwtAuthString(adminSession), nil, secureUser)
	suite.Equal(http.StatusOK, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)
	suite.Nil(user.PreferredHours)

	// 2d. GET /api/users/{username} (as admin)
	code, _ = suite.request("GET", "/api/users/"+username, jwtAuthString(adminSession), nil, secureUser)
	suite.Equal(http.StatusOK, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)
	suite.Nil(user.PreferredHours)

	// 3a. GET /api/users (as admin)
	users := []*model.UserSecure{}
	code, _ = suite.request("GET", "/api/users", jwtAuthString(adminSession), nil, &users)
	suite.Equal(http.StatusOK, code)
	suite.Equal(2, len(users))

	// 3b. GET /api/users (fails as user)
	secureUser = &model.UserSecure{}
	code, _ = suite.request("GET", "/api/users", jwtAuth, nil, secureUser)
	suite.Equal(http.StatusForbidden, code)

	// 4. PATCH /api/users/{userID}
	user = &model.User{PreferredHours: model.NewTimeRange(1, 2)}
	secureUser = &model.UserSecure{}
	code, _ = suite.request("PATCH", "/api/users/"+uid, jwtAuth, user, secureUser)
	suite.Equal(http.StatusOK, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)

	suite.NotNil(secureUser.PreferredHours)
	suite.Equal(1, int(*secureUser.PreferredHours.Start))
	suite.Equal(2, int(*secureUser.PreferredHours.Finish))

	// 4b. PATCH /api/users/{userID}.PreferredHours (fails)
	user = &model.User{PreferredHours: model.NewTimeRange(2, 1)}
	secureUser = &model.UserSecure{}
	code, _ = suite.request("PATCH", "/api/users/"+uid, jwtAuth, user, secureUser)
	suite.Equal(http.StatusBadRequest, code)

	// 4c. PATCH /api/users/{userID}.Role (fails)
	user = &model.User{Role: &model.RoleAdmin}
	secureUser = &model.UserSecure{}
	code, _ = suite.request("PATCH", "/api/users/"+uid, jwtAuth, user, secureUser)
	suite.Equal(http.StatusForbidden, code)

	// 5. DELETE /api/users/{userID}
	secureUser = &model.UserSecure{}
	code, _ = suite.request("DELETE", "/api/users/"+uid, jwtAuth, user, secureUser)
	suite.Equal(http.StatusOK, code)
	suite.Equal(username, secureUser.Username)
	suite.Equal(email, secureUser.Email)

	suite.NotNil(secureUser.PreferredHours)
	suite.Equal(1, int(*secureUser.PreferredHours.Start))
	suite.Equal(2, int(*secureUser.PreferredHours.Finish))
}

func (suite *APITestSuite) Test002_TaskUsage() {
	// 0. GET /api/login (as admin)
	var token map[string]string
	code, _ := suite.request(
		"GET", "/api/login",
		basicAuthString("boss", "test_secret"),
		nil, &token,
	)
	suite.Equal(http.StatusOK, code)

	adminSession, ok := token["session"]

	// 1a. POST /api/users (foo: user)
	username := "foo"
	password := "bar"
	email := "foo@bar.baz"
	user := &model.User{
		Username: &username,
		Password: &password,
		Email:    &email,
	}

	var postUser model.UserSecure
	code, _ = suite.request("POST", "/api/users", "", user, &postUser)
	suite.Equal(http.StatusCreated, code)
	suite.Equal(*user.Username, postUser.Username)

	// 1b. POST /api/users (oof: manager)
	managerUsername := "oof"
	managerPassword := "rab"
	managerEmail := "oof@rab.zab"
	managerUser := &model.User{
		Username: &managerUsername,
		Password: &managerPassword,
		Email:    &managerEmail,
		Role:     &model.RoleManager,
	}

	var postManager model.UserSecure
	code, _ = suite.request("POST", "/api/users", jwtAuthString(adminSession), managerUser, &postManager)
	suite.Equal(http.StatusCreated, code)
	suite.Equal(*managerUser.Username, postManager.Username)

	// 2a. GET /api/login (as foo)
	code, _ = suite.request(
		"GET", "/api/login",
		basicAuthString(username, password),
		nil, &token,
	)
	suite.Equal(http.StatusOK, code)

	session, ok := token["session"]
	suite.True(ok)

	// 2b. GET /api/login (as oof)
	code, _ = suite.request(
		"GET", "/api/login",
		basicAuthString(managerUsername, managerPassword),
		nil, &token,
	)
	suite.Equal(http.StatusOK, code)

	managerSession, ok := token["session"]

	// 3a. POST /api/tasks
	tr1 := model.NewTimeRange(1, 100)
	task1 := &model.Task{
		Title:       "foo",
		Description: "bar",
	}

	// 400 bad request
	code, resp := suite.request("POST", "/api/tasks", jwtAuthString(session), task1, nil)
	suite.Equal(http.StatusBadRequest, code)
	suite.Contains(resp, "userID")

	// 401 unauthorized (as oof)
	task1.UserID = &postUser.ID
	task1.TimeRange = *tr1
	code, _ = suite.request("POST", "/api/tasks", jwtAuthString(managerSession), task1, nil)
	suite.Equal(http.StatusUnauthorized, code)

	// 201 created
	task1.UserID = &postUser.ID
	var postTask1 model.Task
	code, _ = suite.request("POST", "/api/tasks", jwtAuthString(session), task1, &postTask1)
	suite.Equal(task1.Title, postTask1.Title)
	suite.Equal(task1.Description, postTask1.Description)
	suite.Equal(task1.TimeRange.Start, tr1.Start)
	suite.Equal(task1.TimeRange.Finish, tr1.Finish)

	// 3b. POST /api/user/foo/tasks (a second one)
	tr2 := model.NewTimeRange(100, 200)
	task2 := &model.Task{
		UserID:      &postUser.ID,
		Title:       "foo",
		Description: "bar",
		TimeRange:   *tr2,
	}

	// 401 unauthorized (as oof)
	url := fmt.Sprintf("/api/users/%s/tasks", postUser.ID.Hex())
	code, _ = suite.request("POST", url, jwtAuthString(managerSession), task1, nil)
	suite.Equal(http.StatusUnauthorized, code)

	// 201 created (as foo)
	var postTask2 model.Task
	code, _ = suite.request("POST", url, jwtAuthString(session), task2, &postTask2)
	suite.Equal(task2.Title, postTask2.Title)
	suite.Equal(task2.Description, postTask2.Description)
	suite.Equal(task2.TimeRange.Start, tr2.Start)
	suite.Equal(task2.TimeRange.Finish, tr2.Finish)

	// 4a. GET /api/tasks (as foo)
	code, _ = suite.request("GET", "/api/tasks", jwtAuthString(session), nil, nil)
	suite.Equal(http.StatusForbidden, code)

	// 4b. GET /api/tasks (as manager)
	var getTasks []*model.Task
	fmt.Println(postManager.Role, model.PermissionViewAllTasks)
	code, _ = suite.request("GET", "/api/tasks", jwtAuthString(managerSession), nil, &getTasks)
	suite.Equal(http.StatusOK, code)
	suite.Equal(2, len(getTasks))

	// 5a. GET /api/tasks/:taskID (as foo)
	var getTask model.Task
	url = fmt.Sprintf("/api/tasks/%s", postTask1.ID.Hex())
	code, _ = suite.request("GET", url, jwtAuthString(session), nil, &getTask)
	suite.Equal(http.StatusOK, code)
	suite.Equal(postTask1, getTask)

	// 5b. GET /api/tasks/:taskID (as manager)
	code, _ = suite.request("GET", url, jwtAuthString(managerSession), nil, &getTask)
	suite.Equal(http.StatusOK, code)
	suite.Equal(postTask1, getTask)

	// 5c. GET /api/tasks/:taskID (as admin)
	code, _ = suite.request("GET", url, jwtAuthString(adminSession), nil, &getTask)
	suite.Equal(http.StatusOK, code)
	suite.Equal(postTask1, getTask)

	// 6a. GET /api/users/foo/tasks (as foo)
	var getUserTasks []*model.Task
	url = fmt.Sprintf("/api/users/%s/tasks", postUser.ID.Hex())
	code, _ = suite.request("GET", url, jwtAuthString(session), nil, &getUserTasks)
	suite.Equal(http.StatusOK, code)
	suite.Equal(2, len(getUserTasks))

	// 6b. GET /api/users/foo/tasks (as oof)
	url = fmt.Sprintf("/api/users/%s/tasks", postUser.ID.Hex())
	code, _ = suite.request("GET", url, jwtAuthString(managerSession), nil, &getUserTasks)
	suite.Equal(http.StatusOK, code)
	suite.Equal(2, len(getUserTasks))

	// 6c. GET /api/users/foo/tasks (as admin)
	url = fmt.Sprintf("/api/users/%s/tasks", postUser.ID.Hex())
	code, _ = suite.request("GET", url, jwtAuthString(adminSession), nil, &getUserTasks)
	suite.Equal(http.StatusOK, code)
	suite.Equal(2, len(getUserTasks))

	// 7a. PATCH /api/tasks/:taskID (as foo)
	newTitle1 := "lorem"
	taskPatch1 := &model.TaskPatch{
		Title: &newTitle1,
	}
	url = fmt.Sprintf("/api/tasks/%s", postTask1.ID.Hex())

	var patchTask model.Task
	code, _ = suite.request("PATCH", url, jwtAuthString(session), taskPatch1, &patchTask)
	suite.Equal(http.StatusOK, code)
	suite.Equal(newTitle1, patchTask.Title)

	// 7b. PATCH /api/tasks/:taskID (as oof)
	newTitle1 = "ipsum"
	code, _ = suite.request("PATCH", url, jwtAuthString(managerSession), taskPatch1, nil)
	suite.Equal(http.StatusNotFound, code)

	// 7c. PATCH /api/tasks/:taskID (as foo)
	code, _ = suite.request("PATCH", url, jwtAuthString(adminSession), taskPatch1, &patchTask)
	suite.Equal(http.StatusOK, code)
	suite.Equal(newTitle1, patchTask.Title)

	// 8a. DELETE /api/tasks/:taskID1 (as foo)
	var deleteTask model.Task
	code, _ = suite.request("DELETE", url, jwtAuthString(session), nil, &deleteTask)
	suite.Equal(http.StatusOK, code)

	// 8b. DELETE /api/tasks/:taskID2 (as oof)
	url = fmt.Sprintf("/api/tasks/%s", postTask2.ID.Hex())
	code, _ = suite.request("DELETE", url, jwtAuthString(managerSession), nil, nil)
	suite.Equal(http.StatusNotFound, code)

	// 8c. DELETE /api/tasks/:taskID2 (as admin)
	code, _ = suite.request("DELETE", url, jwtAuthString(adminSession), nil, &deleteTask)
	suite.Equal(http.StatusOK, code)
}

func (suite *APITestSuite) request(method, path, auth string, body, response interface{}) (int, string) {
	var req *http.Request
	var err error

	if body != nil {
		// interface to json string
		buf, err := json.Marshal(body)
		suite.Nil(err)

		// create request
		req, err = http.NewRequest(method, path, bytes.NewReader(buf))
		suite.Nil(err)

		// set content-type
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		// create request with no body
		req, err = http.NewRequest(method, path, nil)
		suite.Nil(err)
	}
	// set auth
	if len(auth) > 0 {
		req.Header.Set(echo.HeaderAuthorization, auth)
	}

	// record response
	rec := httptest.NewRecorder()
	suite.e.ServeHTTP(rec, req)
	resp, _ := ioutil.ReadAll(rec.Body)

	// json string to interface if response
	if resp != nil {
		json.Unmarshal(resp, response)
	}

	return rec.Code, string(resp)
}

func basicAuthString(user, pass string) string {
	b64auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))
	return fmt.Sprintf("Basic %s", b64auth)
}

func jwtAuthString(jwt string) string {
	return fmt.Sprintf("Bearer %s", jwt)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
