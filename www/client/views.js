"use strict";

// dt2unix converts a date time string to unix timestamp
function dt2unix(t) {
  return Date.parse(t)/1000;
}

// time2unix converts the start/finish times to unix timestamps
function time2unix(t) {
  return dt2unix('Thu, 01 Jan 1970 ' + t + ':00 GMT');
};

// unix2time converts the unix timestamps to HH:mm format
function unix2time(t) {
  return unix2dt(t).substr(-5);
};

function unix2str(t) {
  var date = new Date(t*1000);
  return date.toUTCString().slice(0, -7);
}

function unix2dt(t) {
  var date = new Date(t*1000);
  var year = date.getUTCFullYear();
  var month = ("0" + (date.getUTCMonth()+1)).substr(-2);
  var day = ("0" + date.getUTCDate()).substr(-2);
  var hours = ("0" + date.getUTCHours()).substr(-2);
  var minutes = ("0" + date.getUTCMinutes()).substr(-2);
  return year + "-" + month + "-" + day + "T" + hours + ":" + minutes
}

// getLabelForInput retrieves the label that corresponds to the given input id
function getLabelForInput(id) {
  return $("label[for="+id+"]");
};

// validateEmail from https://stackoverflow.com/questions/46155/how-to-validate-email-address-in-javascript
function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(email);
};

// IntroView defines the div#introduction view's interface
var IntroView = {
  init: function() {
    // Hide intro if preference is set
    if (localStorage.hideIntro == "true") {
      $("#introduction").hide();
    }
    
    // Setup trigger to prompt and hide onclick
    $("#hideIntro").unbind("click");
    $("#hideIntro").click(IntroView.promptAndHide);
  },

  // promptAndHide asks the user if she would like to hide 
  // the intro forever and then hides it
  promptAndHide: function(event) {
    localStorage.hideIntro = confirm("Would you like to hide this message permanently?");
    $("#introduction").hide();
  }
};

// MainView defines the interface to toggle subviews
var MainView = {
  // toggleLoginRegister switches the visibility states of login and register
  toggleLoginRegister: function(event) {
    $("#registerForm").toggle();
    $("#loginForm").toggle();
    event.preventDefault();
  },

  // hideAll hides all the views calling this should be followed by showing a view
  hideAll: function() {
    $("#registerForm").hide();
    $("#loginForm").hide();
    $("#menuView").hide();
    $("#profileView").hide();
    $("#usersListView").hide();
    $("#tasksListView").hide();
    $("#taskView").hide();
  },
  
  // showMenu shows menu view and hides all others 
  showMenu: function() {
    MainView.hideAll();
    $("#menuView").show();
  },
  
  // showLogin shows login view and hides all others
  showLogin: function() {
    MainView.hideAll();
    $("#loginForm").show();
  },

  // showProfile shows profile view and hides all others
  showProfile: function() {
    MainView.hideAll();
    $("#profileView").show();
  },

  // showUsersList shows users list view and hides all others 
  showUsersList: function() {
    MainView.hideAll();
    $("#usersListView").show();
  },

  // showTaskView shows task view and hides all others 
  showTaskView: function() {
    MainView.hideAll();
    $("#taskView").show();
  },

  // showTasksList shows tasks list view and hides all others 
  showTasksList: function() {
    MainView.hideAll();
    $("#tasksListView").show();
  }
};

// LoginView defines the form#loginForm interface
var LoginView = {
  clear: function() {
    $("#loginUsername").val("");
    $("#loginPassword").val("");
  },
  init: function(submitHandler) {
    // Setup trigger to toggle with register
    $("#showRegister").unbind("click");
    $("#showRegister").click(MainView.toggleLoginRegister);

    // Setup trigger to handle form submission
    $("#loginSubmit").unbind("click");
    $("#loginSubmit").click(submitHandler);
  },

  validate: function() {
    if (LoginView.username().length == 0) {
      return "Must specify a username"
    }
    if (LoginView.password().length == 0) {
      return "Must specify a password"
    }
    return "";
  },

  username: function() {
    return $("#loginUsername").val();
  },

  password: function() {
    return $("#loginPassword").val();
  }
};

// RegisterView defines the div#registerForm interface
var RegisterView = {
  clear: function() {
    $("#registerUsername").val(""),
    $("#registerEmail").val(""),
    $("#registerPassword").val("")
    $("#registerConfirmPassword").val("")
    $("#registerPreferredHoursStart").val("");
    $("#registerPreferredHoursFinish").val("");
  },
  init: function(submitHandler) {
    RegisterView.clear()

    // Setup trigger to toggle with login
    $("#showLogin").unbind("click");
    $("#showLogin").click(MainView.toggleLoginRegister);

    // Setup trigger to handle form submission
    $("#registerSubmit").unbind("click");
    $("#registerSubmit").click(submitHandler);
  },

  // validate ensures that the modified input fields 
  //   are in a proper state for submission
  validate: function() {
    // validate username
    if (!$("#registerUsername").val().length) {
      return "Must specify a username";
    }
    // validate email
    var email = $("#registerEmail").val();
    if (!email.length || !validateEmail(email)) {
      return "Invalid email";
    }
    // confirm pw
    var pw = $("#registerPassword").val();
    if (!pw.length) {
      return "Must specify a password";
    }
    // confirm that the password was input correctly a second time
    if (pw != $("#registerConfirmPassword").val()) {
      return "Passwords do not match";
    }
    // confirm that the start and finish times, if specified are sequential
    var startTime = $("#registerPreferredHoursStart").val()
    var finishTime = $("#registerPreferredHoursFinish").val()
    if (startTime.length && finishTime.length &&
        (time2unix(startTime) > time2unix(finishTime))) {
      return "Preferred Hours start time must be less than finish"
    }
    return "";
  },

  // user retrieves the values from the register input fields
  //   and returns a constructed user object
  user: function() {
    // text fields are easy
    var user = { 
      username: $("#registerUsername").val(),
      email: $("#registerEmail").val(),
      password: $("#registerPassword").val()
    };   

    // preferred hours are optional
    var start = $("#registerPreferredHoursStart").val();
    if (start.length == 0) {
      return user;
    }
    var finish = $("#registerPreferredHoursFinish").val();
    if (finish.length == 0) {
      return user;
    }
    user.preferredHours = {start: time2unix(start), finish: time2unix(finish)};
    return user;
  }
}

// MenuView defines the interface to the main menu
var MenuView = {
  // newItem creates a new menu list item with the given name and id
  newItem: function(name, id) {
    return '<a href="#" id="'+id+'"><li class="list-group-item">'+name+'</li></a>';
  },
  
  // init resets the menu view and 
  //   setups the items based on the input user's permissions
  init: function(user, callbacks) {
    $("#menuList").empty();

    // Create menu items
    $("#menuList").append(MenuView.newItem("View your profile", "viewProfile"));

    if (userHasPermission(user, permissionModifySelfTasks)) {
      $("#menuList").append(MenuView.newItem("View your tasks", "viewTasks"));
    }
    if (userHasPermission(user, permissionViewAllTasks)) {
      $("#menuList").append(MenuView.newItem("View all tasks", "viewAllTasks"));
    }
    if (userHasPermission(user, permissionModifyAllUsersRestricted)){
      $("#menuList").append(MenuView.newItem("View all users", "viewUsers"));
    }

    $("#menuList").append(MenuView.newItem("Logout", "logout"));

    // Hook up the menu to callbacks
    for (var id in callbacks) {
      $(id).click(callbacks[id]);
    }
  }
};

// ProfileView defines the interface for the profile menu
var ProfileView = {
  clear: function() {
    $("#profileUsername").val("");
    $("#profileEmail").val("");
    $("#profileRoleGroup").hide();
    $("#profileRole").val("");
    $("#profilePreferredHoursStart").val("");
    $("#profilePreferredHoursFinish").val("");
    $("#profileOldPassword").val("");
    $("#profileNewPassword").val("");
    $("#profileConfirmPassword").val("");
    $("#profileCancel").unbind("click");
    $("#profileSubmit").unbind("click");
    $("#profileDelete").unbind("click");
  },

  // init clears out the previous profile view and fills it with the input user
  init: function(user, showRole, cancelHandler, updateHandler, deleteHandler, viewTasksHandler) {
    ProfileView.clear();

    $("#profileUsername").val(user.username);
    $("#profileEmail").val(user.email);

    if (showRole) {
      $("#profileRoleGroup").show();
      $("#profileRole").val(roleIntToStr[user.role]);
    }

    if (user.preferredHours != null) {
      $("#profilePreferredHoursStart").val(unix2time(user.preferredHours.start));
      $("#profilePreferredHoursFinish").val(unix2time(user.preferredHours.finish));
    }

    $("#profileCancel").click(cancelHandler);
    $("#profileSubmit").click(updateHandler);
    $("#profileDelete").click(deleteHandler);
  },

  // validate ensures that the modified input fields 
  //   are in a proper state for submission
  validate: function() {
    // if a new password is provided, 
    // confirm that it's been input correctly a second time
    var newPw = $("#profileNewPassword").val();
    if (newPw.length && (newPw != $("#profileConfirmPassword").val())) {
      return "Passwords do not match";
    }
    // if preferred hours are specified, make sure they are sequential
    var startTime = $("#profilePreferredHoursStart").val()
    var finishTime = $("#profilePreferredHoursFinish").val()
    if (startTime.length && finishTime.length &&
        (time2unix(startTime) > time2unix(finishTime))) {
      return "Preferred Hours start time must be less than finish"
    }
    return "";
  },

  user: function() {
    // text fields are easy
    var user = { 
      username: $("#profileUsername").val(),
      email: $("#profileEmail").val(),
    };   

    // only put in password if specified
    var password = $("#profileNewPassword").val();
    if (password.length > 0) {
      user.oldPassword = $("#profileOldPassword").val();
      user.password = password;
    }

    // role
    var role = $("#profileRole").val();
    if (role) {
      user.role = roleStrToInt[role];
    }

    // preferred hours are optional
    var start = $("#profilePreferredHoursStart").val();
    if (start.length == 0) {
      return user;
    }
    var finish = $("#profilePreferredHoursFinish").val();
    if (finish.length == 0) {
      return user;
    }
    user.preferredHours = {start: time2unix(start), finish: time2unix(finish)};
    return user;
  }
};

// UserListView ...
var UsersListView = {
  clear: function() {
    $("#usersList").find("tr:gt(0)").remove();
    $("#usersListBack").unbind("click");
  },

  newItem: function(user) {
    var el = '<tr id="'+user.id+'">' +
      '<td>'+user.id+'</td>' +
      '<td>'+user.username+'</td>' +
      '<td>'+user.email+'</td>' +
      '<td>' +
      (user.preferredHours ? unix2time(user.preferredHours.start) : "None") +
      '</td>' +
      '<td>' +
      (user.preferredHours ? unix2time(user.preferredHours.finish) : "None") +
      '</td></tr>'
    return el;
  },

  init: function(users, backHandler, showUserHandler) {
    UsersListView.clear();

    for (var i in users) {
      var user = users[i];
      $("#usersList").append(UsersListView.newItem(user));
      $("#"+user.id).click(function(event) {
        showUserHandler(this.id);
      });
    }
    $("#usersListBack").click(backHandler);
  }
}

// TaskView
var TaskView = {
  clear: function() {
    $("#taskTitle").val("");
    $("#taskDescription").val("");
    $("#taskUser").empty();
    $("#taskStartTime").val("");
    $("#taskFinishTime").val("");
    $("#taskStartTime").val("");
    $("#taskBack").unbind("click");
    $("#taskUpdate").unbind("click");
    $("#taskDelete").unbind("click");
  },
  init: function(task, canEdit, backHandler, updateHandler, deleteHandler) {
    TaskView.clear();

    $("#taskTitle").val(task.title);
    $("#taskDescription").val(task.description);

    $("#taskUser").append("<option>"+task.userID+"</option>");
    $("#taskUser").val(task.userID);

    if (task.start) {
      $("#taskStartTime").val(unix2dt(task.start));
    }
    if (task.finish) {
      $("#taskFinishTime").val(unix2dt(task.finish));
    }

    $("#taskBack").click(backHandler);
    $("#taskUpdate").click(updateHandler);

    if (deleteHandler == null) {
      $("#taskDelete").hide();
    } else {
      $("#taskDelete").show();
      $("#taskDelete").click(deleteHandler);
    }
  },

  task: function() {
    return {
      title: $("#taskTitle").val(),
      description: $("#taskDescription").val(),
      userID: $("#taskUser").val(),
      start: dt2unix($("#taskStartTime").val()),
      finish: dt2unix($("#taskFinishTime").val())
    }
  }
}

// TasksListView ...
var TasksListView = {
  clear: function() {
    $("#tasksList").find("tr:gt(0)").remove();
    $("#tasksListFilter").unbind("click");
    $("#tasksListBack").unbind("click");
    $("#tasksListAdd").unbind("click");
  },
  newItem: function(task) {
    return '<tr style=' +
      '"background:'+(task.conflict ? "red" : "green")+';'+
      'color:white" id="'+task.id+'">' +
    '<td>'+(task.user ? task.user : "")+'</td>' +
    '<td>'+task.title+'</td>' +
    '<td>'+task.description.slice(0, 32)+'</td>' +
    '<td>'+unix2str(task.start)+'</td>' +
    '<td>'+unix2str(task.finish)+'</td></tr>'
  },
  getFilterParams: function() {
    var from = $("#tasksListFilterFrom").val();
    if (from) {
      from = dt2unix(from) * 1000;
    }
    var to = $("#tasksListFilterTo").val();
    if (to) {
      to = dt2unix(to) * 1000;
    }
    return [from, to];
  },
  init: function(tasks, filterHandler, exportHandler, showTaskHandler, backHandler, addHandler) {
    TasksListView.clear()

    for (var i in tasks) {
      var task = tasks[i];
      $("#tasksList").append(TasksListView.newItem(task));
      $("#"+task.id).click(function(event) {
        showTaskHandler(this.id);
      });
    }
    $("#tasksListFilter").click(function(event) {
      event.preventDefault();
      var filters = TasksListView.getFilterParams();
      filterHandler(filters[0], filters[1]);
    });
    $("#tasksListExport").click(function(event) {
      event.preventDefault();
      var filters = TasksListView.getFilterParams();
      exportHandler(filters[0], filters[1]);
    });
    $("#tasksListBack").click(backHandler);
    $("#tasksListAdd").click(addHandler);
  }
};
