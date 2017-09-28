package store

import (
	"fmt"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/briansan/ManageMeServer/errors"
	"github.com/briansan/ManageMeServer/model/schema"
)

const (
	tasksCollectionName = "tasks"
)

func newTaskQueryByID(id string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id)}
}

// NewTaskQueryFromParams constructs a mongo query from string inputs
//   error can be safely ignored if from and to are passed a empty strings
func NewTaskQueryFromParams(userID, taskID, from, to string) (bson.M, error) {
	q := bson.M{}

	// userID
	if len(userID) > 0 {
		q["userID"] = bson.ObjectIdHex(userID)
	}

	// taskID
	if len(taskID) > 0 {
		q["_id"] = bson.ObjectIdHex(taskID)
	}

	// from
	if len(from) > 0 {
		if fromI, err := strconv.Atoi(from); err != nil {
			return nil, errors.NewValidationError("from", "int")
		} else {
			q["finish"] = bson.M{"$gte": fromI}
		}
	}

	// to
	if len(to) > 0 {
		if toI, err := strconv.Atoi(to); err != nil {
			return nil, errors.NewValidationError("to", "int")
		} else {
			q["start"] = bson.M{"$lte": toI}
		}
	}
	return q, nil
}

// GetTasksCollection returns an mgo instance to the tasks collection
func (m *MongoStore) GetTasksCollection() *mgo.Collection {
	return m.GetDatabase().C(tasksCollectionName)
}

// CreateTask inserts task object into db
// error is 500 if mongo fails, 409 if user exists, else nil
func (m *MongoStore) CreateTask(task *schema.Task) error {
	// UserID is a required field
	if task.UserID == nil {
		return fmt.Errorf("task must contain userID")
	}
	// Try to insert and return error
	task.ID = bson.NewObjectId()
	if err := m.GetTasksCollection().Insert(task); err != nil {
		return err
	}
	return nil
}

// GetAllTasks retrieves all tasks for option userID
func (m *MongoStore) GetAllTasks(q bson.M) ([]*schema.Task, error) {
	// Fetch the tasks
	tasks := []*schema.Task{}
	err := m.GetTasksCollection().Find(q).All(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask looks up task in db with given query for entire object
// error is 500 if mongo fails, else nil
func (m *MongoStore) GetTask(q bson.M) (*schema.Task, error) {
	task := schema.Task{}
	err := m.GetTasksCollection().Find(q).One(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask ...
func (m *MongoStore) UpdateTask(taskID string, taskPatch *schema.TaskPatch) (*schema.Task, error) {
	// Try to update the user
	q := newTaskQueryByID(taskID)
	changeInfo := mgo.Change{
		Update:    bson.M{"$set": taskPatch},
		Upsert:    false,
		ReturnNew: true,
	}
	task := schema.Task{}
	_, err := m.GetTasksCollection().Find(q).Apply(changeInfo, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTask removes task from db with given taskID
// error is 500 if mongo fails, else nil
func (m *MongoStore) DeleteTask(taskID string) (*schema.Task, error) {
	task, err := m.GetTask(newTaskQueryByID(taskID))
	if err != nil {
		return nil, err
	}

	if err = m.GetTasksCollection().Remove(newTaskQueryByID(taskID)); err != nil {
		return task, err
	}

	return task, nil
}

// DeleteTasks removes tasks for given userID
// error is 500 if mongo fails, else nil
func (m *MongoStore) DeleteTasksForUser(userID string) error {
	_, err := m.GetTasksCollection().RemoveAll(bson.M{"userID": bson.ObjectIdHex(userID)})
	return err
}
