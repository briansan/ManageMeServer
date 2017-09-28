package schema

import (
	"encoding/json"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test001_Roles(t *testing.T) {
	// Test anon
	assert.True(t, RoleHasPermission(RoleAnon, PermissionCreateUser))
	assert.False(t, RoleHasPermission(RoleAnon, PermissionModifySelfTasks))
	assert.False(t, RoleHasPermission(RoleAnon, PermissionModifySelfTasks))
	assert.False(t, RoleHasPermission(RoleAnon, PermissionModifyAllUsers))
	assert.False(t, RoleHasPermission(RoleAnon, PermissionModifyAllUsersRestricted))
	assert.False(t, RoleHasPermission(RoleAnon, PermissionViewAllTasks))
	assert.False(t, RoleHasPermission(RoleAnon, PermissionModifyAllTasks))

	// Test User
	assert.False(t, RoleHasPermission(RoleUser, PermissionCreateUser))
	assert.True(t, RoleHasPermission(RoleUser, PermissionModifySelfTasks))
	assert.False(t, RoleHasPermission(RoleUser, PermissionModifyAllUsers))
	assert.False(t, RoleHasPermission(RoleUser, PermissionModifyAllUsersRestricted))
	assert.False(t, RoleHasPermission(RoleUser, PermissionViewAllTasks))
	assert.False(t, RoleHasPermission(RoleUser, PermissionModifyAllTasks))

	// Test Manager
	assert.False(t, RoleHasPermission(RoleManager, PermissionCreateUser))
	assert.True(t, RoleHasPermission(RoleManager, PermissionModifySelfTasks))
	assert.False(t, RoleHasPermission(RoleManager, PermissionModifyAllUsers))
	assert.True(t, RoleHasPermission(RoleManager, PermissionModifyAllUsersRestricted))
	assert.True(t, RoleHasPermission(RoleManager, PermissionViewAllTasks))
	assert.False(t, RoleHasPermission(RoleManager, PermissionModifyAllTasks))

	// Test Admin
	assert.False(t, RoleHasPermission(RoleAdmin, PermissionCreateUser))
	assert.True(t, RoleHasPermission(RoleAdmin, PermissionModifySelfTasks))
	assert.True(t, RoleHasPermission(RoleAdmin, PermissionModifyAllUsers))
	assert.True(t, RoleHasPermission(RoleAdmin, PermissionModifyAllUsersRestricted))
	assert.True(t, RoleHasPermission(RoleAdmin, PermissionViewAllTasks))
	assert.True(t, RoleHasPermission(RoleAdmin, PermissionModifyAllTasks))
}

func Test002_User(t *testing.T) {
	// Test that preferred times is nil on init
	var err error
	u := User{}

	// Test email field
	err = u.Validate()
	assert.Equal(t, "email field is required as string", err.Error())

	// Test username field
	json.Unmarshal([]byte(`{"email": "foo"}`), &u)
	err = u.Validate()
	assert.Equal(t, "username field is required as string", err.Error())
	assert.Equal(t, "foo", *u.Email)

	// Test password field
	json.Unmarshal([]byte(`{"email": "foo", "username": "bar"}`), &u)
	err = u.Validate()
	assert.Equal(t, "password field is required as string", err.Error())
	assert.Equal(t, "bar", *u.Username)

	// Test good
	json.Unmarshal([]byte(`{"email": "foo", "username": "bar", "password": "baz"}`), &u)
	err = u.Validate()
	assert.Nil(t, err)
	assert.Equal(t, "baz", *u.Password)

	// Test preferredHours.finish as wrong type
	json.Unmarshal([]byte(`{"email": "foo", "username": "bar", "password": "baz", "preferredHours": {"start": 100, "finish": "200"}}`), &u)
	err = u.Validate()
	assert.Equal(t, "finish field is required as unix timestamp (int)", err.Error())

	// Test preferredHours.start < preferredHours.finish
	json.Unmarshal([]byte(`{"email": "foo", "username": "bar", "password": "baz", "preferredHours": {"start": 100, "finish": 50}}`), &u)
	err = u.Validate()
	assert.Equal(t, "start field is required as less than finish", err.Error())

	// Test all good
	json.Unmarshal([]byte(`{"email": "foo", "username": "bar", "password": "baz", "preferredHours": {"start": 100, "finish": 200}}`), &u)
	err = u.Validate()
	assert.Equal(t, 100, *u.PreferredHours.Start)
	assert.Equal(t, 200, *u.PreferredHours.Finish)

	// Test oldPassword
	json.Unmarshal([]byte(`{"email": "foo", "username": "bar", "password": "baz", "oldPassword": "foobar"}`), &u)
	assert.Equal(t, "foobar", *u.OldPassword)
}

func Test003_Task(t *testing.T) {
	// Test time range
	tr := &TimeRange{}
	err := tr.Validate()
	assert.Equal(t, "start field is required as unix timestamp (int)", err.Error())

	// Test user field
	task := &Task{}
	err = task.Validate()
	assert.Equal(t, "user field is required as string", err.Error())

	// Test title field
	json.Unmarshal([]byte(`{"user": "foo"}`), task)
	err = task.Validate()
	assert.Equal(t, "title field is required as string", err.Error())

	// Test start field
	json.Unmarshal([]byte(`{"user": "foo", "title": "bar"}`), task)
	err = task.Validate()
	assert.Equal(t, "start field is required as unix timestamp (int)", err.Error())

	// Test finish field
	json.Unmarshal([]byte(`{"user": "foo", "title": "bar", "start": 1}`), task)
	err = task.Validate()
	assert.Equal(t, "finish field is required as unix timestamp (int)", err.Error())

	// Test good
	json.Unmarshal([]byte(`{"user": "foo", "title": "bar", "start": 1, "finish": 1}`), task)
	err = task.Validate()
	assert.Nil(t, err)
}
