function loadComplete() {
  signInStateChange();
}

function signInStateChange() {
  isSignedIn(function(signedIn) {
    document.querySelector('#preAuthSpinner').style.display = 'none';
    if (!signedIn) {
      document.querySelector('#signinButton').style.display = 'block';
      document.querySelector('#signoutButton').style.display = 'none';
    } else {
      document.querySelector('#signinButton').style.display = 'none';
      document.querySelector('#signoutButton').style.display = 'block';
      var $injector = angular.bootstrap(document, [window.AngularApp]);
    }
  });
}

function init() {
  var ROOT = '//' + window.location.host + '/_ah/api';
  loadLibraries([{name: 'action', version: 'v1', root: ROOT}, 
	         {name: 'oauth2', version: 'v2'}],
      loadComplete, ROOT);
}

//------------- Generic Functions to help wih init.
function loadLibraries(libraries, callback) {
  // libNames keeps track of all libraries which are not executed yet.
  var libNames = [];

  // complete() generates a function that when invokes, markes library 'name'
  // as been loaded. That is, the returned function removes it from the
  // list of names.
  function complete(name) {
    return function() {
      var index = libNames.indexOf(name);
      libNames.splice(index, 1);
      if (libNames.length == 0) {
	callback();
      }
    }
  }

  libraries.forEach(function(e) {
    libNames.push(e.name);
    gapi.client.load(e.name, e.version, complete(e.name), e.root);
  });
}


// ---------------- Generic function to help with auth.
// The cb should take a single argument that is a boolean.
// Set to true if user is signed in, false if not.
function isSignedIn(cb) {
  gapi.auth.authorize({
      client_id: CLIENT_ID,
      immediate: true,
      scope: ['https://www.googleapis.com/auth/userinfo.email']
    }, function(response) {
      cb(!response.error);
    });
}

function signInWithPopup(postSignIn) {
  gapi.auth.authorize({
      client_id: CLIENT_ID,
      immediate: false,
      scope: ['https://www.googleapis.com/auth/userinfo.email']
  }, postSignIn);
}
