package api

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"

	"github.com/briansan/ManageMeServer/errors"
	model "github.com/briansan/ManageMeServer/model/schema"
	"github.com/briansan/ManageMeServer/model/store"
)

// GetTasks retrieves all tasks
//   available to roles with ViewAllTasks
func GetTasks(c echo.Context) error {
	// Type assert user from context and authorize
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok {
		return echo.ErrUnauthorized
	}

	// Get userID from param and decide authorization
	userID := c.QueryParam("userID")
	if !allows(user.Role, model.PermissionViewAllTasks) {
		if len(userID) == 0 {
			userID = user.ID.Hex()
		} else if userID != user.ID.Hex() {
			return echo.ErrForbidden
		}
	}

	// Get db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	}
	defer db.Cleanup()

	// Construct query
	q, err := store.NewTaskQueryFromParams(
		userID, "",
		c.QueryParam("from"), c.QueryParam("to"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Fetch tasks
	tasks, err := db.GetAllTasks(q)
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	return c.JSON(http.StatusOK, tasks)
}

func PostTasks(c echo.Context) error {
	t := model.Task{}
	c.Bind(&t)

	// Validate
	if err := t.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// If not admin or task is for user, 401
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok {
		return echo.ErrUnauthorized
	}
	ownerIsPosting := user.ID.Hex() == t.UserID.Hex()
	isAdmin := allows(user.Role, model.PermissionModifyAllTasks)
	if !ownerIsPosting && !isAdmin {
		return echo.ErrUnauthorized
	}

	// Get db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	defer db.Cleanup()

	// Try to add task
	if err := db.CreateTask(&t); err != nil {
		return errors.MongoErrorResponse(err)
	}

	return c.JSON(http.StatusCreated, t)
}

func GetTaskByID(c echo.Context) error {
	taskID := c.Param("taskID")

	// Type assert user from context and 401 if no user
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok {
		return echo.ErrForbidden
	}

	// Establish db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	defer db.Cleanup()

	// Construct query: no permission to view all tasks
	//   restricts query to include userID
	userID := ""
	if !allows(user.Role, model.PermissionViewAllTasks) {
		userID = user.ID.Hex()
	}
	q, _ := store.NewTaskQueryFromParams(userID, taskID, "", "")

	// Fetch task
	t, err := db.GetTask(q)
	if err != nil {
		return errors.MongoErrorResponse(err)
	}

	return c.JSON(http.StatusOK, t)
}

func PatchTask(c echo.Context) error {
	taskID := c.Param("taskID")

	// Type assert user from context and authorize
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok {
		return echo.ErrUnauthorized
	}

	// Get user patch doc
	taskPatch := &model.TaskPatch{}
	c.Bind(taskPatch)

	// Return if empty patch
	if taskPatch.NoChange() {
		return echo.NewHTTPError(http.StatusBadRequest, "body is empty")
	}

	// Establish db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	defer db.Cleanup()

	// Construct query: no permission to view all tasks
	//   restricts query to include userID
	userID := ""
	if !allows(user.Role, model.PermissionModifyAllTasks) {
		userID = user.ID.Hex()
	}
	q, _ := store.NewTaskQueryFromParams(userID, taskID, "", "")

	// Fetch task
	t, err := db.GetTask(q)
	if err != nil {
		return errors.MongoErrorResponse(err)
	}

	// Try to update task
	if t, err = db.UpdateTask(t.ID.Hex(), taskPatch); err != nil {
		return errors.MongoErrorResponse(err)
	}
	return c.JSON(http.StatusOK, t)
}

func DeleteTask(c echo.Context) error {
	taskID := c.Param("taskID")

	// Type assert user from context and authorize
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok {
		return echo.ErrUnauthorized
	}

	// Establish db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	defer db.Cleanup()

	// Construct query: no permission to view all tasks
	//   restricts query to include userID
	userID := ""
	if !allows(user.Role, model.PermissionModifyAllTasks) {
		userID = user.ID.Hex()
	}
	q, _ := store.NewTaskQueryFromParams(userID, taskID, "", "")

	// Fetch task
	t, err := db.GetTask(q)
	if err != nil {
		return errors.MongoErrorResponse(err)
	}

	// Try to delete task
	t, err = db.DeleteTask(taskID)
	if err != nil {
		return errors.MongoErrorResponse(err)
	}

	return c.JSON(http.StatusOK, t)
}

func GetUserTasks(c echo.Context) error {
	userID := c.Param("userID")

	// Type assert user from context and authorize
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok || (user.ID.Hex() != userID && !allows(user.Role, model.PermissionViewAllTasks)) {
		return echo.ErrUnauthorized
	}

	// Establish db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	defer db.Cleanup()

	// Construct query
	q, err := store.NewTaskQueryFromParams(
		userID, "", c.QueryParam("from"), c.QueryParam("to"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Fetch task (include userID in query if no permissions to modify all
	tasks, err := db.GetAllTasks(q)
	if err != nil {
		return errors.MongoErrorResponse(err)
	}

	return c.JSON(http.StatusOK, tasks)
}

func PostUserTasks(c echo.Context) error {
	userID := c.Param("userID")

	t := model.Task{}
	c.Bind(&t)

	// set user id
	uid := bson.ObjectIdHex(userID)
	t.UserID = &uid

	// Validate
	if err := t.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// If not admin or task is for user, 401
	user, ok := c.Get("user").(*model.UserSecure)
	if !ok {
		return echo.ErrUnauthorized
	}
	ownerIsPosting := user.ID.Hex() == t.UserID.Hex()
	isAdmin := allows(user.Role, model.PermissionModifyAllTasks)
	if !ownerIsPosting && !isAdmin {
		return echo.ErrUnauthorized
	}

	// Get db connection
	db, err := store.NewMongoStore()
	if err != nil {
		return errors.MongoErrorResponse(err)
	}
	defer db.Cleanup()

	// Try to add task
	if err := db.CreateTask(&t); err != nil {
		return errors.MongoErrorResponse(err)
	}

	return c.JSON(http.StatusCreated, t)
}

func initTasks(api *echo.Group) {
	api.GET("/tasks", GetTasks, DoJWTAuth)
	api.POST("/tasks", PostTasks, DoJWTAuth)
	api.GET("/tasks/:taskID", GetTaskByID, DoJWTAuth)
	api.PATCH("/tasks/:taskID", PatchTask, DoJWTAuth)
	api.DELETE("/tasks/:taskID", DeleteTask, DoJWTAuth)
	api.GET("/users/:userID/tasks", GetUserTasks, DoJWTAuth)
	api.POST("/users/:userID/tasks", PostUserTasks, DoJWTAuth)
}
