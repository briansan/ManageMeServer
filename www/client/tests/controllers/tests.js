//
// Controller tests

QUnit.test("AuthController", function(assert) {
  // Test getUserID
  localStorage.user = '{"id": "foo"}';
  assert.equal("foo", AuthController.getUserID())

  // Test load
  AuthController.load()
  assert.equal(AuthController.login, LoginView.submitHandler)
  assert.equal(AuthController.register, RegisterView.submitHandler)
  assert.equal(true, MainView.showAuthCalled)

  // Test register
  AuthController.register(mockEvent)
  assert.equal(RegisterView.user(), ManageMeAPI.registerUser)
  assert.equal(AuthController.registerSuccessLogin, ManageMeAPI.registerSuccess)
  assert.equal(alertError, ManageMeAPI.registerFailure)

  // Test login
  AuthController.login(mockEvent)
  assert.equal(LoginView.username(), ManageMeAPI.loginUsername)
  assert.equal(LoginView.password(), ManageMeAPI.loginPassword)
  assert.equal(AuthController.loginSuccess, ManageMeAPI.loginSuccess)
  assert.equal(alertError, ManageMeAPI.loginFailure)

  // Test logout
  localStorage.jwt = "foo";
  localStorage.user = "bar";
  MainView.showAuthCalled = false;
  AuthController.logout(null);
  assert.equal(undefined, localStorage.jwt);
  assert.equal(undefined, localStorage.user);
  assert.equal(true, MainView.showAuthCalled);

  // Test loginSuccess
  AuthController.loginSuccess({session: "e30K.eyJhdWQiOiJmb28ifQo=.e30K"})
  assert.equal("foo", parseJwt(localStorage.jwt).aud)
  assert.equal(true, LoginView.clearCalled)
  assert.equal(true, RegisterView.clearCalled)

  // Test registerSucessLogin
  delete ManageMeAPI.loginUsername;
  delete ManageMeAPI.loginPassword;
  delete ManageMeAPI.loginSuccess;
  delete ManageMeAPI.loginFailuer;
  RegisterView.userReturn = {username: "foo", password: "bar"};
  AuthController.registerSuccessLogin(null);
  assert.equal("foo", ManageMeAPI.loginUsername);
  assert.equal("bar", ManageMeAPI.loginPassword);
  assert.equal(AuthController.loginSuccess, ManageMeAPI.loginSuccess);
  assert.equal(alertError, ManageMeAPI.loginFailure);
});

QUnit.test("ProfileController", function(assert) {
  // Test load
  localStorage.user = '"foo"'
  ProfileController.load(null);
  assert.equal("foo", ProfileView.initUser);
  assert.equal(MenuController.load, ProfileView.cancelHandler)
  assert.equal(ProfileController.updateUser, ProfileView.updateHandler)
  assert.equal(ProfileController.deleteAccount, ProfileView.deleteHandler)
  assert.equal(true, MainView.showProfileCalled);

  // Test updateUser
  localStorage.user = '{"id": "foo"}';
  ProfileController.updateUser(mockEvent);
  assert.equal("foo", ManageMeAPI.patchID)
  assert.equal(ProfileView.userReturn, ManageMeAPI.patchUser)
  assert.equal(alertError, ManageMeAPI.patchFailure)

  // Test deleteUser
  localStorage.user = '{"username": "foo"}';
  promptReturn = "bar";
  ProfileController.deleteAccount(mockEvent);
  assert.equal("foo", ManageMeAPI.loginUsername);
  assert.equal("bar", ManageMeAPI.loginPassword);
});
