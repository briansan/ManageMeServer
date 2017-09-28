"use strict";

var baseURL = "http://localhost:8888";

function log(msg) {
  console.log(msg);
}

function getAPIURL() {
  $.ajax({url: "/api"}).done(function(url) {
    log("using api url: "+url);
    baseURL = url;
  }).fail(function(resp) {
    alert("failed to get api url");
  });
}

getAPIURL();

function request(method, path, data, successHandler, failureHandler) {
  var req = {
    url: baseURL + path,
    method: method,
    data: JSON.stringify(data),
    headers: {
      "Content-type": "application/json; charset=utf-8"
    }
  }
  if (localStorage.jwt) {
    req.headers["Authorization"] = "Bearer " + localStorage.jwt;
  }
  return $.ajax(req).done(successHandler).fail(failureHandler);
};

var ManageMeAPI = {
  login: function(username, password, successHandler, failureHandler) {
    $.ajax({
      url: baseURL + "/api/login",
      beforeSend: function (xhr) {
        xhr.setRequestHeader ("Authorization", "Basic " + btoa(username + ":" + password));
      }
    }).done(successHandler).fail(failureHandler);
  },
  register: function(user, successHandler, failureHandler) {
    return request("POST", "/api/users", user, successHandler, failureHandler);
  },
  getMappedUsers: function(successHandler, failureHandler) {
    return request("GET", "/api/users?mapped=true", "", successHandler, failureHandler);
  },
  getUsers: function(successHandler, failureHandler) {
    return request("GET", "/api/users", "", successHandler, failureHandler);
  },
  getUser: function(userID, successHandler, failureHandler) {
    return request("GET", "/api/users/"+userID, "", successHandler, failureHandler);
  },
  patchUser: function(userID, user, successHandler, failureHandler) {
    return request("PATCH", "/api/users/"+userID, user, successHandler, failureHandler);
  },
  deleteUser: function(userID, successHandler, failureHandler) {
    return request("DELETE", "/api/users/"+userID, "", successHandler, failureHandler);
  },
  getTasks: function(userID, from, to, successHandler, failureHandler) {
    var url = "/api/tasks?"
    if (userID) {
      url += "userID="+userID+"&"
    }
    if (from) {
      url += "from="+from+"&"
    }
    if (to) {
      url += "to="+to
    }
    return request("GET", url, "", successHandler, failureHandler);
  },
  postTasks: function(task, successHandler, failureHandler) {
    return request("POST", "/api/tasks", task, successHandler, failureHandler);
  },
  getTask: function(taskID, successHandler, failureHandler) {
    return request("GET", "/api/tasks/"+taskID, "", successHandler, failureHandler);
  },
  patchTask: function(taskID, task, successHandler, failureHandler) {
    return request("PATCH", "/api/tasks/"+taskID, task, successHandler, failureHandler);
  },
  deleteTask: function(taskID, successHandler, failureHandler) {
    return request("DELETE", "/api/tasks/"+taskID, "", successHandler, failureHandler);
  }
};

var permissionCreateUser               = 1 << 0;
var permissionModifySelfTasks          = 1 << 1;
var permissionModifyAllUsers           = 1 << 2;
var permissionModifyAllUsersRestricted = 1 << 3;
var permissionViewAllTasks             = 1 << 4;
var permissionModifyAllTasks           = 1 << 5;

var roleUser = permissionModifySelfTasks;
var roleManager = roleUser | permissionModifyAllUsersRestricted | permissionViewAllTasks;
var roleAdmin = roleManager | permissionModifyAllUsers | permissionModifyAllTasks;

function userHasPermission(user, perm) {
  return (user.role & perm) > 0;
};

var roleStrToInt = {
  User: roleUser,
  Manager: roleManager,
  Admin: roleAdmin
}

var roleIntToStr = {}
roleIntToStr[roleUser] = "User";
roleIntToStr[roleManager] = "Manager";
roleIntToStr[roleAdmin] = "Admin";
