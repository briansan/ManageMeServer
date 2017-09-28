package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"

	"github.com/briansan/ManageMeServer/model/schema"
)

type StoreTestSuite struct {
	suite.Suite
	store *MongoStore
}

func (suite *StoreTestSuite) SetupTest() {
	// Use test database and reestablish session
	os.Setenv(envDatabaseName, "test")
	InitMongoSession()

	var err error
	suite.store, err = NewMongoStore()
	suite.Nil(err)

	suite.store.GetUsersCollection().RemoveAll(nil)
	suite.store.GetTasksCollection().RemoveAll(nil)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

// Test001_User asserts proper CRUD functionality of user object with mongo
func (suite *StoreTestSuite) Test001_User() {
	username := "foo"
	email := "bar"
	pw := "baz"
	role := 1337

	// Test CreateUser
	newUser := &schema.User{
		Username: &username,
		Email:    &email,
		Password: &pw,
		Role:     &role,
	}
	err := suite.store.CreateUser(newUser)
	suite.Nil(err)

	// Test GetUserByUsername
	user, err := suite.store.GetUserByUsername(username)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(email, user.Email)

	id := user.ID.Hex()

	// Test GetUserByEmail
	user, err = suite.store.GetUserByEmail(email)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(email, user.Email)

	// Test GetUserByPassword
	user, err = suite.store.GetUserByCreds(username, pw)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(email, user.Email)

	// Test CreateUser with conflict
	err = suite.store.CreateUser(newUser)
	suite.NotNil(err)
	suite.Equal("user with username as foo already exists", err.Error())

	// Test UpdateUser
	newUsername := "foobar"
	userPatch := &schema.User{Username: &newUsername}
	user, err = suite.store.UpdateUser(id, userPatch)
	suite.Nil(err)
	suite.Equal(newUsername, user.Username)
	suite.Equal(email, user.Email)
	suite.Equal(role, user.Role)

	// Add second user
	err = suite.store.CreateUser(newUser)
	suite.Nil(err)

	// Try to update second user
	u, err := suite.store.GetUserByUsername(*newUser.Username)
	suite.Nil(err)

	user, err = suite.store.UpdateUser(u.ID.Hex(), userPatch)
	suite.Nil(user)
	suite.True(mgo.IsDup(err))

	// Test GetAllUsers
	users, err := suite.store.GetAllUsers()
	suite.Nil(err)
	suite.Equal(len(users), 2)

	// Test DeleteUser
	user, err = suite.store.DeleteUser(id)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(newUsername, user.Username)
	suite.Equal(email, user.Email)
}

// Test002_Admin asserts proper admin insertion
func (suite *StoreTestSuite) Test002_Admin() {
	// Create admin
	err := suite.store.AdminExistsOrCreate("test_secret")
	suite.Nil(err)

	// Fetch admin
	u, err := suite.store.GetUserByCreds("boss", "test_secret")
	suite.Nil(err)
	suite.Equal(schema.RoleAdmin, u.Role)
}

// Test003_Task asserts proper task usage
func (suite *StoreTestSuite) Test003_Task() {
	username := "foo"
	email := "bar"
	pw := "baz"
	role := 1337

	// Test CreateUser
	newUser := &schema.User{
		Username: &username,
		Email:    &email,
		Password: &pw,
		Role:     &role,
	}
	err := suite.store.CreateUser(newUser)
	suite.Nil(err)

	// Test GetUserByUsername
	user, err := suite.store.GetUserByUsername(username)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(email, user.Email)

	id := user.ID.Hex()

	// Test CreateTask
	tr := schema.NewTimeRange(100, 1000)
	task := schema.Task{
		TimeRange:   *tr,
		Title:       "foo",
		Description: "bar",
	}

	// Will fail since no userID
	err = suite.store.CreateTask(&task)
	suite.Equal("task must contain userID", err.Error())

	// Make it pass
	task.UserID = &user.ID
	err = suite.store.CreateTask(&task)
	suite.NoError(err)

	// Test GetAllTasksByUserID
	q, err := NewTaskQueryFromParams(user.ID.Hex(), "", "", "")
	suite.NoError(err)
	tasks1, err := suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(1, len(tasks1))

	t := tasks1[0]
	getTask, err := suite.store.GetTask(newTaskQueryByID(t.ID.Hex()))
	suite.NoError(err)
	suite.Equal(t, getTask)

	// Test GetAllTasksByTimeRange
	q, err = NewTaskQueryFromParams("", "", "1", "50")
	suite.NoError(err)
	tasks2, err := suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(0, len(tasks2))

	q, err = NewTaskQueryFromParams("", "", "1", "100")
	suite.NoError(err)
	tasks2, err = suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(1, len(tasks2))

	q, err = NewTaskQueryFromParams("", "", "95", "105")
	suite.NoError(err)
	tasks2, err = suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(1, len(tasks2))

	q, err = NewTaskQueryFromParams("", "", "900", "1000")
	suite.NoError(err)
	tasks2, err = suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(1, len(tasks2))

	q, err = NewTaskQueryFromParams("", "", "1001", "1002")
	suite.NoError(err)
	tasks2, err = suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(0, len(tasks2))

	// Test GetAllTasksByUserIDTimeRange
	q, err = NewTaskQueryFromParams(user.ID.Hex(), "", "1", "50")
	suite.NoError(err)
	tasks3, err := suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(0, len(tasks3))

	q, err = NewTaskQueryFromParams(user.ID.Hex(), "", "100", "101")
	suite.NoError(err)
	tasks3, err = suite.store.GetAllTasks(q)
	suite.NoError(err)
	suite.Equal(1, len(tasks3))

	// Test UpdateTask
	newTitle := "fooo"
	newDesc := "descr"
	newTr := schema.NewTimeRange(10, 100)
	patchTask := schema.TaskPatch{
		TimeRange:   *newTr,
		Title:       &newTitle,
		Description: &newDesc,
	}

	newT, err := suite.store.UpdateTask(t.ID.Hex(), &patchTask)
	suite.Equal(newTitle, newT.Title)
	suite.Equal(newDesc, newT.Description)
	suite.Equal(*newTr.Start, *newT.TimeRange.Start)
	suite.Equal(*newTr.Finish, *newT.TimeRange.Finish)

	// Test UpdateTask without time rnage
	patchTask.TimeRange.Start = nil
	patchTask.TimeRange.Finish = nil
	newTitle = "foo"
	newDesc = "desc"

	newT, err = suite.store.UpdateTask(t.ID.Hex(), &patchTask)
	suite.Equal(newTitle, newT.Title)
	suite.Equal(newDesc, newT.Description)
	suite.Equal(*newTr.Start, *newT.TimeRange.Start)
	suite.Equal(*newTr.Finish, *newT.TimeRange.Finish)

	// Test DeleteTask
	delT, err := suite.store.DeleteTask(newT.ID.Hex())
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(delT, newT)

	// Test DeleteUser
	user, err = suite.store.DeleteUser(id)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(email, user.Email)
}
