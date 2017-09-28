// Mock confirm
function mockConfirm(resp) {
  window.confirm = function(msg) {
    return resp
  }
}

// Mock alert
window.alert = function(msg) {return msg;}

// Mock prompt
var promptReturn = "";
window.prompt = function(msg) {return promptReturn;}

// Mock event
mockEvent = {
  preventDefault: function() {}
}
