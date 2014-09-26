var phonecatApp = angular.module('ActionsApp', []);

phonecatApp.controller('ActionsCtrl', function($scope, $timeout) {
    $scope.startActionAdd = function() {
        gapi.client.action.add({
            redirect: $scope.redirect,
            words: $scope.words
        }).execute($scope.endActionAdd);
	$scope.formDisabled = true;
	$scope.message = '';
	$scope.error = '';
	$scope.searchPhrase = '';
	$scope.items = [];
    }

    $scope.endActionAdd = function(resp) {
	$scope.formDisabled = false;
	if (resp.error) {
	  $scope.error = resp.error.message;
	} else {
	  $scope.redirect = '';
	  $scope.words = '';
	  $scope.setMessage('Added');
	  $scope.fetchItems();
	}
	$scope.$apply();
    }

    $scope.clearMessage = function() {
      $scope.message = '';
    }

    $scope.setMessage = function(message) {
      $scope.message = message;
      $timeout($scope.clearMessage, 1000);
    }

    var POLL_INTERVAL = 1000;

    $scope.fetchItems = function() {
      if (!gapi || !gapi.client || !gapi.client.action || !gapi.client.action.list) {
	$timeout($scope.fetchItems, POLL_INTERVAL);
	return;
      }

      gapi.client.action.list({phrase: $scope.searchPhrase}).execute(function(resp) {
	$scope.items = resp.items;
	$scope.$apply();
      });
    }
    $scope.fetchItems();
});
