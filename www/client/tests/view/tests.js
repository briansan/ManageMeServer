//
// View tests

QUnit.test("time2unix", function(assert) {
  assert.equal(0, time2unix("00:00"));
  assert.equal(60, time2unix("00:01"));
  assert.equal(600, time2unix("00:10"));
  assert.equal(3600, time2unix("01:00"));
  assert.equal(36000, time2unix("10:00"));
});

QUnit.test("unix2time", function(assert) {
  assert.equal("00:00", unix2time(0));
  assert.equal("00:00", unix2time(59));
  assert.equal("00:01", unix2time(60));
  assert.equal("00:10", unix2time(600));
  assert.equal("01:00", unix2time(3600));
  assert.equal("10:00", unix2time(36000));
});

QUnit.test("labelForInput", function(assert) {
  assert.equal($("#labelForInput")[0], getLabelForInput("testGetLabelForInput")[0]);
});

QUnit.test("validateEmail", function(assert) {
  assert.equal(true, validateEmail("foo@bar.baz"));
  assert.equal(false, validateEmail("foo@bar"));
  assert.equal(false, validateEmail("foo"));
});

QUnit.test("IntroView", function(assert) {
  // Test for the hiding of intro view based on localStorage setting
  localStorage.hideIntro = "false";
  IntroView.init();
  assert.equal(true, $("#introduction").is(":visible"));

  localStorage.hideIntro = "true";
  IntroView.init();
  assert.equal(false, $("#introduction").is(":visible"));

  // Test promptAndHide
  mockConfirm(true);
  IntroView.promptAndHide();
  assert.equal("true", localStorage.hideIntro);
  assert.equal(false, $("#introduction").is(":visible"));
});

QUnit.test("MainView", function(assert) {
  // Test toggleLoginRegister
  $("#registerForm").hide();
  $("#loginForm").hide();
  MainView.toggleLoginRegister(mockEvent);
  assert.equal(true, $("#registerForm").is(":visible"));
  assert.equal(true, $("#loginForm").is(":visible"));

  // Test showMenu
  MainView.showMenu();
  assert.equal(false, $("#registerForm").is(":visible"));
  assert.equal(false, $("#loginForm").is(":visible"));
  assert.equal(true, $("#menuView").is(":visible"));
  assert.equal(false, $("#profileView").is(":visible"));

  // Test showAuth
  MainView.showAuth();
  assert.equal(false, $("#registerForm").is(":visible"));
  assert.equal(true, $("#loginForm").is(":visible"));
  assert.equal(false, $("#menuView").is(":visible"));
  assert.equal(false, $("#profileView").is(":visible"));

  // Test showProfile
  MainView.showProfile();
  assert.equal(false, $("#registerForm").is(":visible"));
  assert.equal(false, $("#loginForm").is(":visible"));
  assert.equal(false, $("#menuView").is(":visible"));
  assert.equal(true, $("#profileView").is(":visible"));
});

QUnit.test("LoginView", function(assert) {
  // mock submit
  submitClicked = false;
  testSubmitHandler = function() {
    submitClicked = true;
  };

  // init
  LoginView.init(testSubmitHandler);

  // Test showRegister click event
  $("#registerForm").hide();
  $("#loginForm").show();
  $("#showRegister").click();
  assert.equal(true, $("#registerForm").is(":visible"));
  assert.equal(false, $("#loginForm").is(":visible"));

  // Test submit click event
  $("#loginSubmit").click();
  assert.equal(true, submitClicked);

  // Test validate (as well as username and password)
  assert.equal("Must specify a username", LoginView.validate());
  $("#loginUsername").val("foo");
  assert.equal("foo", LoginView.username());

  assert.equal("Must specify a password", LoginView.validate());
  $("#loginPassword").val("bar");
  assert.equal("bar", LoginView.password());
  assert.equal("", LoginView.validate());
});

QUnit.test("RegisterView", function(assert) {
  // mock submit
  submitClicked = false;
  testSubmitHandler = function() {
    submitClicked = true;
  };

  // init
  RegisterView.init(testSubmitHandler);

  // Test showLogin click event
  $("#registerForm").show();
  $("#loginForm").hide();
  $("#showLogin").click();
  assert.equal(false, $("#registerForm").is(":visible"));
  assert.equal(true, $("#loginForm").is(":visible"));

  // Test submit
  $("#registerSubmit").click();
  assert.equal(true, submitClicked);
  
  // Test validate
  assert.equal("Must specify a username", RegisterView.validate());
  $("#registerUsername").val("foo");
  assert.equal("Invalid email", RegisterView.validate());
  $("#registerEmail").val("foo@bar.baz");
  assert.equal("Must specify a password", RegisterView.validate());
  $("#registerPassword").val("bar");
  assert.equal("Passwords do not match", RegisterView.validate());
  $("#registerConfirmPassword").val("bar");

  // Test start/finish preferred times
  $("#registerPreferredHoursStart").val("10:00");
  $("#registerPreferredHoursFinish").val("09:00");
  assert.equal("Preferred Hours start time must be less than finish", RegisterView.validate());
  $("#registerPreferredHoursFinish").val("11:00");
  assert.equal("", RegisterView.validate());

  // Test user
  user = RegisterView.user();
  assert.equal("foo", user.username);
  assert.equal("bar", user.password);
  assert.equal("foo@bar.baz", user.email);

  $("#registerPreferredHoursStart").val("");
  $("#registerPreferredHoursFinish").val("");
  user = RegisterView.user();
  assert.equal(null, user.preferredHours);

  $("#registerPreferredHoursStart").val("10:00");
  user = RegisterView.user();
  assert.equal(null, user.preferredHours);

  $("#registerPreferredHoursFinish").val("11:00");
  user = RegisterView.user();
  assert.equal(36000, user.preferredHours.start);
  assert.equal(39600, user.preferredHours.finish);
});

QUnit.test("MenuView", function(assert) {
  assert.equal(
    '<a href="#" id="bar"><li class="list-group-item">foo</li></a>',
    MenuView.newItem("foo", "bar")
  );
});

QUnit.test("ProfileView", function(assert) {
  var user = {
    username: "foo",
    email: "foo@bar.baz",
    preferredHours: {
      start: 3600,
      finish: 36000
    }
  };

  // cancelHandler
  var cancelClicked = false;
  cancelHandler = function() {
    cancelClicked = true;
  };

  // updateHandler
  var updateClicked = false;
  updateHandler = function() {
    updateClicked = true;
  };

  // deleteHandler
  var deleteClicked = false;
  deleteHandler = function() {
    deleteClicked = true;
  };

  // Test loading the profile
  ProfileView.init(user, cancelHandler, updateHandler, deleteHandler);
  assert.equal("foo", $("#profileUsername").val());
  assert.equal("foo@bar.baz", $("#profileEmail").val());
  assert.equal("01:00", $("#profilePreferredHoursStart").val());
  assert.equal("10:00", $("#profilePreferredHoursFinish").val());

  // Test the button handlers
  $("#profileCancel").click();
  assert.equal(true, cancelClicked);
  $("#profileSubmit").click();
  assert.equal(true, updateClicked);
  $("#profileDelete").click();
  assert.equal(true, deleteClicked);

  // Test validate password mismatch
  $("#profileNewPassword").val("foo");
  $("#profileConfirmPassword").val("bar");
  assert.equal("Passwords do not match", ProfileView.validate());

  // Test validate preferred hours out of order
  $("#profileConfirmPassword").val("foo");
  $("#profilePreferredHoursStart").val("10:00");
  $("#profilePreferredHoursFinish").val("01:00");
  assert.equal("Preferred Hours start time must be less than finish", ProfileView.validate());

  // Test validate OK
  $("#profilePreferredHoursFinish").val("11:00");
  assert.equal("", ProfileView.validate());

  // Test user
  var profileUser = ProfileView.user();
  assert.equal(user.username, profileUser.username);
  assert.equal(user.email, profileUser.email);
  assert.equal(36000, profileUser.preferredHours.start);
  assert.equal(39600, profileUser.preferredHours.finish);

  // Test user no preferredHours
  $("#profilePreferredHoursStart").val("");
  $("#profilePreferredHoursFinish").val("");
  profileUser = ProfileView.user();
  assert.equal(undefined, profileUser.preferredHours);

  $("#profilePreferredHoursStart").val("10:00");
  profileUser = ProfileView.user();
  assert.equal(undefined, profileUser.preferredHours);
  
  $("#profilePreferredHoursFinish").val("11:00");
  profileUser = ProfileView.user();
  assert.equal(36000, profileUser.preferredHours.start);
  assert.equal(39600, profileUser.preferredHours.finish);
});
