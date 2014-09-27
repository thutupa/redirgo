var actionsApp = angular.module('ActionsApp', []);

actionsApp.controller('ActionsCtrl', function($scope, $timeout) {
  $scope.startActionAdd = function() {
    gapi.client.action.add({
        redirect: $scope.redirect,
        words: $scope.words
    }).execute($scope.endActionAdd);
    $scope.formDisabled = true;
    $scope.message = '';
    $scope.error = '';
    $scope.searchPhrase = '';
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

  $scope.fetchItems = function() {
    gapi.client.action.list({phrase: $scope.searchPhrase}).execute(function(resp) {
      $scope.items = resp.items;
      $scope.$apply();
    });
  }
  $scope.items = [];
  $scope.fetchItems();
});
