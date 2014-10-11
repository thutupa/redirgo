var actionsApp = angular.module('ActionsApp', []);

var REFILTER_DELAY_MS = 250;

actionsApp.controller('ActionsCtrl', function($scope, $timeout) {
  $scope.startActionAdd = function() {
    var toAdd = {
        redirectLink: $scope.addRedirectLink,
        actionWords: $scope.addActionWords,
        inMemory: true
    };
    $scope.items.push(toAdd);

    gapi.client.action.add(toAdd).execute($scope.endActionAdd);
    $scope.formDisabled = true;
    $scope.message = '';
    $scope.error = '';
    $scope.filterPhrase = '';
  }

  $scope.endActionAdd = function(resp) {
    $scope.formDisabled = false;
    if (resp.error) {
      $scope.error = resp.error.message;
    } else {
      $scope.addRedirectLink = '';
      $scope.addActionWords = '';
      $scope.setMessage('Added');
      $scope.fetchItems();
    }
  }

  $scope.initActionEdit = function(index) {
    var editedItem = $scope.items[index];
    editedItem.editRedirectLink = editedItem.redirectLink;
    editedItem.editActionWords = editedItem.actionWords.join(' ');
    editedItem.inEdit = true;
  }

  // When submit button is clicked during the "edit".
  $scope.startActionEdit = function(index) {
    var editedItem = $scope.items[index];
    editedItem.error = '';
    editedItem.editDisabled = true;
    gapi.client.action.edit({
      id: editedItem.id,
      redirectLink: editedItem.editRedirectLink,
      actionWords: editedItem.editActionWords
    }).execute(function(resp) {
      if (resp.error) {
	editedItem.error = resp.error;
	editDisabled = false;
      } else {
	$scope.endActionEdit(index);
      }
    });
  }

  // When cancel button is clicked during the "edit".
  $scope.cancelActionEdit = function(index) {
    var editedItem = $scope.items[index];
    editedItem.inEdit = false;
  };

  $scope.endActionEdit = function(index) {
    var editedItem = $scope.items[index];
    editedItem.error = '';
    editedItem.inEdit = false;
    editedItem.editDisabled = false;
    editedItem.redirectLink = editedItem.editRedirectLink;
    editedItem.actionWords = editedItem.editActionWords.split(' ');
    $scope.$apply();
  };

  $scope.clearMessage = function() {
    $scope.message = '';
  };

  $scope.setMessage = function(message) {
    $scope.message = message;
    $timeout($scope.clearMessage, 1000);
  };

  $scope.fetchItems = function() {
    $scope.fetchItemsInProgress = true;
    gapi.client.action.list({phrase: $scope.filterPhrase}).execute(function(resp) {
      $scope.fetchItemsInProgress = false;
      $scope.items = resp.items || [];
      $scope.$apply();
    });
  };

  $scope.delayedInitRefilter = function() {
    // Cancel any existing promise and ignore any
    // errors in doing so.
    if ($scope.refilterPromise) {
      try {
	$scope.refilterPromise.cancel();
	console.log('Cancelled existing promise');
      } catch(e) {
	console.log('Cancelling existing promise raised exception ' + e);
      }
    }

    $scope.refilterPromise = $timeout(function() {
      $scope.refilterPromise = undefined;
      $scope.fetchItems();
    }, REFILTER_DELAY_MS, false /* don't call under $apply -- fetchItems already does it */);
  };
  $scope.items = [];
  $scope.fetchItems();
});
