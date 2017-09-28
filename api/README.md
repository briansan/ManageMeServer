# ManageMe RESTful HTTP service

## Models

### TimeRange
start    int (unix timestamp)
finish   int (unix timestamp)

### User
```
id             bson.ObjectID
username       string
password       string
email          string
role           int
preferred_time TimeRange
```

### Task
```
id           bson.ObjectID
user_id      bson.ObjectID
title        string
description  string
start        int (unix timestamp)
finish       int (unix timestamp)
```

## Permissions
```
CreateUser:
  can create user
ModifySelfTasks:
  can view/modify self
  can CRUD tasks where task.user_id = self
ModifyAllUsers: 
  can CRUD all users
ModifyAllUsersRestricted:
  like ModifyAllUsers except:
    cannot modify users where role = Admin
    cannot modify user.role
ViewAllTasks:
  can read all tasks
ModifyAllTasks: 
  can CRUD all tasks
```

## Roles
```
Anon: CreateUser
User: ModifySelfTasks
Manager: User + ModifyAllUserRestricted + ViewAllTasks
Admin: Manager + ModifyAllUsers + ModifyAllTasks
```

## API
all routes mounted on `/api`

### GET /service/ping
- allows: All
- details: healthcheck endpoint reporting version

### GET /login
- allows: All
- details: presents authenticated user with 1 hr jwt session
- requires: BasicAuth

### GET /users
- allows: Manager, Admin
- details: retrieves all users
- requires: Bearer JWT Auth

### POST /users
- allows: Anon, Manager, Admin
- details: creates a user
- requires: Bearer JWT Auth

### GET /users/:userID
- allows: User*, Manager, Admin
- details: retrieves a user by id
- requires: Bearer JWT Auth

### PATCH /users/:userID
- allows: User*, Manager, Admin
- details: updates a user by field
- requires: Bearer JWT Auth

### DELETE /users/:userID
- allows: User*, Manager, Admin
- details: deletes a user and all associated tasks
- requires: Bearer JWT Auth

### GET /users/:userID/tasks
- allows: User*, Manager, Admin
- details: retrieves all tasks for user
- requires: Bearer JWT Auth

### POST /users/:userID/tasks
- allows: User*, Manager*, Admin*
- details: creates a task for user
- requires: Bearer JWT Auth

### GET /tasks
- allows: Manager, Admin
- details: retrieves all tasks
- requires: Bearer JWT Auth

### POST /tasks
- allows: User*, Manager*, Admin*
- details: creates a task
- requires: Bearer JWT Auth

### GET /tasks/:id
- allows: User*, Manager, Admin
- details: retrieves a task
- requires: Bearer JWT Auth

### PATCH /tasks/:id
- allows: User*, Manager*, Admin
- details: updates a task by field
- requires: Bearer JWT Auth

### DELETE /tasks/:id
- allows: User*, Manager*, Admin
- details: deletes a task
- requires: Bearer JWT Auth

[^*]: only allowed for resources owned by that role's user
