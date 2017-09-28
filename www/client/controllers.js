"use strict";

// alertError will pop up an alert when the API returns an error response
function alertError(resp) {
  alert(resp.responseJSON.message);
};

// parseJwt returns the claims section of the jwt
function parseJwt(token) {
  var base64Url = token.split('.')[1];
  var base64 = base64Url.replace('-', '+').replace('_', '/');
  return JSON.parse(window.atob(base64));
};

// AuthController defines the login/register functionality as well as
//   storage of jwt/user object
var AuthController = {
  // getUser returns the parsed user obj in local storage
  getUser: function() {
    return JSON.parse(localStorage.user);
  },

  // getUserID returns the logged in user's id
  getUserID: function() {
    return AuthController.getUser().id;
  },

  // load sets up the login and register views
  load: function(event) {
    LoginView.init(AuthController.login);
    RegisterView.init(AuthController.register);
    MainView.showLogin();
  },

  // register will register the user when the register.submit button is pressed
  register: function(event) {
    event.preventDefault();
  
    // validate
    var err = RegisterView.validate();
    if (err.length > 0) {
      alert(err);
      return;
    }
  
    // register user to api
    ManageMeAPI.register(RegisterView.user(), AuthController.registerSuccessLogin, alertError);
  },
  
  // login will login the user when the login.submit button is pressed
  login: function(event) {
    event.preventDefault();
  
    // validate
    var err = LoginView.validate();
    if (err.length > 0) {
      alert(err);
      return;
    }
  
    // login
    ManageMeAPI.login(LoginView.username(), LoginView.password(),
      AuthController.loginSuccess, alertError);
  },

  // logout clears the local storage and returns to the auth view
  logout: function(event) {
    delete localStorage.jwt;
    delete localStorage.user;
    AuthController.load();
  },

  // loginSuccess is the handler for a successful login
  loginSuccess: function(resp) {
    // save the auth data in the local storage
    var jwt = resp.session;
    localStorage.jwt = jwt;
  
    // load menu
    MenuController.load();

    // clear out the auth views
    LoginView.clear();
    RegisterView.clear();
  
  },
  
  // registerSuccessLogin will login the user from a successful registration
  registerSuccessLogin: function(resp) {
    var user = RegisterView.user();
    alert("Welcome to ManageMe!");
    ManageMeAPI.login(user.username, user.password, AuthController.loginSuccess, alertError);
  }
  
};

var ProfileController = {
  init: function(user) {
    ProfileController.user = user;
  },

  load: function() {
    ProfileView.init(
      ProfileController.user, 
      userHasPermission(AuthController.getUser(), permissionModifyAllUsers),
      MenuController.load, 
      ProfileController.updateUser, 
      ProfileController.confirmDelete);
    MainView.showProfile();
  },

  updateSuccess: function(user) {
      alert("Update successful!");
      if (AuthController.getUserID() == user.id) {
        localStorage.user = JSON.stringify(user);
        MenuController.load();
      } else {
        UsersListController.load();
      }
  },

  updateUser: function(event) {
    event.preventDefault();
  
    // validate
    var err = ProfileView.validate()
    if (err.length > 0) {
      alert(err)
      return;
    }
  
    // call patch
    var user = ProfileView.user()
    ManageMeAPI.patchUser(ProfileController.user.id, user, ProfileController.updateSuccess, alertError);
  },

  deleteAccount: function(msg) {
    // Say thank you and logout
    ManageMeAPI.deleteUser(ProfileController.user.id, function(msg) {
      if (AuthController.getUserID() == msg.id) {
        alert("Thank you for using ManageMe");
        AuthController.logout();
      } else {
        MainView.showMenu();
      }
    }, alertError);
  },

  confirmDelete: function(event) {
    event.preventDefault();
  
    // Make sure they want to delete
    var pw = prompt("If you are sure you wish to permanently delete your account and all associated tasks, type in your password to confirm");
    if (pw == null) {
      return
    }

    // Try to login to confirm ownership and delete if success
    ManageMeAPI.login(AuthController.getUser().username, pw, ProfileController.deleteAccount, alertError);
  }
};

var UsersListController = {
  load: function() {
    ManageMeAPI.getUsers(UsersListController.loadView, alertError);
  },

  userClicked: function(userID) {
    ManageMeAPI.getUser(userID, function(user) {
      ProfileController.init(user);
      ProfileController.load();
    }, alertError);
  },

  loadView: function(users) {
    UsersListView.init(users, MainView.showMenu, UsersListController.userClicked);
    MainView.showUsersList();
  }

};

var TaskViewController = {
  init: function(task) {
    TaskViewController.task = task;
  },

  loadCreate: function() {
    // Some magic to convert the current time from localtime to UTC
    var now = Date.now();
    now -= (new Date(now)).getTimezoneOffset()*60000;
    var before = now - 3600000;

    TaskViewController.task = {
      id: null,
      title: "",
      description: "",
      userID: AuthController.getUserID(),
      start: before/1000,
      finish: now/1000
    };
    TaskView.init(
      TaskViewController.task,
      userHasPermission(AuthController.getUser(), permissionModifyAllTasks),
      TaskViewController.backClicked,
      TaskViewController.addTask,
      null
    );
    MainView.showTaskView();
  },

  load: function() {
    TaskView.init(
      TaskViewController.task,
      userHasPermission(AuthController.getUser(), permissionModifyAllTasks),
      TaskViewController.backClicked,
      TaskViewController.updateClicked,
      TaskViewController.deleteClicked
    );
    MainView.showTaskView();
  },

  backClicked: function(event) {
    event.preventDefault();
    TaskView.clear();
    MainView.showTasksList();
  },

  updateSuccess: function(task) {
    alert("Update successful!");
    TasksListController.load();
  },

  addTask: function(event) {
    event.preventDefault();
    ManageMeAPI.postTasks(
      TaskView.task(),
      TaskViewController.updateSuccess,
      alertError
    );
  },

  updateClicked: function(event) {
    event.preventDefault();
    ManageMeAPI.patchTask(
      TaskViewController.task.id,
      TaskView.task(),
      TaskViewController.updateSuccess,
      alertError
    );
  },

  deleteSuccess: function(task) {
    alert("Task successfully deleted");
    TasksListController.load();
  },

  deleteClicked: function(event) {
    event.preventDefault();
    var conf = prompt("If you are sure you wish to permanently delete this task, type in the current title of this task");
    if (conf != TaskViewController.task.title) {
      alert("Wrong title, please try again");
      return;
    }

    // Try to delete
    ManageMeAPI.deleteTask(TaskViewController.task.id, TaskViewController.deleteSuccess, alertError);
  }
}

var TasksListController = {
  init: function(all) {
    TasksListController.all = all;
  },
  load: function() {
    ManageMeAPI.getTasks(
      TasksListController.all ? "" : AuthController.getUserID(),
      null,
      null,
      TasksListController.loadView,
      alertError
    );
  },

  taskClicked: function(taskID) {
    ManageMeAPI.getTask(taskID, function(task) {
      TaskViewController.init(task);
      TaskViewController.load();
    }, alertError);
  },

  loadView: function(tasks) {
    if (TasksListController.all) {
      ManageMeAPI.getMappedUsers(function(users) {
        for (var i = 0; i < tasks.length; i++) {
          var task = tasks[i]
          var user = users[task.userID];
          if (!user) continue;
          tasks[i].user = user.username;
          var userHrs = user.preferredHours;
          if (!userHrs) continue;
          var start = task.start % 86400;
          var finish = task.finish % 86400;
          tasks[i].conflict = (start < userHrs.finish) && (finish > userHrs.start);
        }
        TasksListView.init(
          tasks,
          TasksListController.filterClicked,
          TasksListController.exportClicked,
          TasksListController.taskClicked,
          MainView.showMenu,
          TaskViewController.loadCreate
        );
        MainView.showTasksList();
      }, alertError)
      return
    } else {
      for (var i = 0; i < tasks.length; i++) {
        var task = tasks[i]
        tasks[i].user = AuthController.getUser().username;
        var userHrs = AuthController.getUser().preferredHours;
        if (!userHrs) continue;
        var start = task.start % 86400;
        var finish = task.finish % 86400;
        tasks[i].conflict = (start < userHrs.finish) && (finish > userHrs.start);
      }
      TasksListView.init(
        tasks,
        TasksListController.filterClicked,
        TasksListController.exportClicked,
        TasksListController.taskClicked,
        MainView.showMenu,
        TaskViewController.loadCreate
      );
      MainView.showTasksList();
    }
  },

  filterClicked: function(from, to) {
    ManageMeAPI.getTasks(TasksListController.all ? "" : AuthController.getUserID(),
      from, to, TasksListController.loadView, alertError);
  },

  exportClicked: function(from, to) {
    ManageMeAPI.getTasks(TasksListController.all ? "" : AuthController.getUserID(),
      from, to, function(tasks) {
        TasksListController.exportToHTML(from, to, tasks);
      }, alertError);
  },

  exportToHTML: function(from, to, tasks) {
    // initial html
    var start = new Date(from);
    var html = "<p><b>Date: </b>"+start.getUTCDate()+"."+(start.getUTCMonth()+1)+"</p>";
    var td = to-from;
    html += "<p><b>Total Time: </b>"+Math.floor(td/3600)+"h"+Math.floor((td%3600)/60)+"m</p>";
    html += "<b>Notes:</b><ul>"
    for (var i = 0; i < tasks.length; i++) {
      var task = tasks[i];
      html += "<li><b>"+task.title+": </b>"+task.description+"</li>"
    }
    html += "</ul>"
    
    // export it
    var link = document.createElement("a");
    link.setAttribute("download", "report.html");
    link.setAttribute("href", "data:text/html;charset=utf8,"+encodeURIComponent(html));
    link.click();
  }
}

function loadSelfTasks(event) {
  TasksListController.init(false);
  TasksListController.load();
}

function loadAllTasks(event) {
  TasksListController.init(true);
  TasksListController.load();
}

function loadProfile(event) {
  ProfileController.init(AuthController.getUser());
  ProfileController.load();
}

var MenuController = {
  menuMap: {
    "#viewProfile": loadProfile,
    "#viewTasks": loadSelfTasks,
    "#viewAllTasks": loadAllTasks,
    "#viewUsers": UsersListController.load,
    "#logout": AuthController.logout
  },

  load: function() {
    // refuse to load menu if no jwt in local storage
    if (localStorage.jwt == undefined) {
      return;
    }
  
    // fetch the user object and save into local storage
    var claims = parseJwt(localStorage.jwt);
    ManageMeAPI.getUser(claims.aud, function(msg) {
      localStorage.user = JSON.stringify(msg);
  
      // switch to menu view
      MenuView.init(msg, MenuController.menuMap);
      MainView.showMenu();
    }, function(resp) {
      // failure to fetch user object means session is expired, show login page
      alert("Bad authentication for session, logging out");
      AuthController.logout();
    });
  }
};

// Main
function main() {
  IntroView.init();
  AuthController.load();
  if (localStorage.jwt) {
    MenuController.load();
  }
};
