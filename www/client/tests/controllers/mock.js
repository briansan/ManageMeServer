var LoginView = {
  submitHandler: null,
  init: function(submitHandler) {
    this.submitHandler = submitHandler;
  },

  validateReturn: "",
  validate: function() {
    return this.validateReturn;
  },

  username: function() {
    return "foo";
  },

  password: function() {
    return "bar";
  },

  clearCalled: false,
  clear: function() {
    this.clearCalled = true;
  }
};

var RegisterView = {
  submitHandler: null,
  init: function(submitHandler) {
    this.submitHandler = submitHandler;
  },

  validateReturn: "",
  validate: function() {
    return this.validateReturn;
  },

  userReturn: "foo",
  user: function() {
    return this.userReturn;
  },

  clearCalled: false,
  clear: function() {
    this.clearCalled = true;
  }
};


var MainView = {
  showAuthCalled: false,
  showAuth: function() {
    this.showAuthCalled = true;
  },

  showProfileCalled: false,
  showProfile: function() {
    this.showProfileCalled = true;
  },

  reset: function() {
    this.showAuthCalled = false;
  }
};

var ProfileView = {
  // init
  initUser: null,
  cancelHandler: null,
  updateHandler: null,
  deleteHandler: null,
  init: function(user, cancelHandler, updateHandler, deleteHandler) {
    this.initUser = user
    this.cancelHandler = cancelHandler;
    this.updateHandler = updateHandler;
    this.deleteHandler = deleteHandler;
  },

  // validate
  validateReturn: "",
  validate: function() {
    return this.validateReturn;
  },

  // user
  userReturn: {
    username: "foo",
    password: "bar",
    email: "foo@bar.baz",
  },
  user: function() {
    return this.userReturn;
  }
};

var ManageMeAPI = {
  getUser: function(id, success, failure) { },

  patchID: null,
  patchUser: null,
  patchSuccess: null,
  patchFailure: null,
  patchUser: function(id, userPatch, success, failure) { 
    this.patchID = id;
    this.patchUser = userPatch;
    this.patchSuccess = success;
    this.patchFailure = failure;
  },

  registerUser: null,
  registerSuccess: null,
  registerFailure: null,
  register: function(user, success, failure) {
    this.registerUser = user;
    this.registerSuccess = success;
    this.registerFailure = failure;
  },

  loginUsername: null,
  loginPassword: null,
  loginSuccess: null,
  loginFailure: null,
  login: function(username, password, success, failure) {
    this.loginUsername = username;
    this.loginPassword = password;
    this.loginSuccess = success;
    this.loginFailure = failure;
  }
};
