var app = angular.module('BreathApp', ['ui.slider']);

var REDRAW_INTERVAL_MS = 50;

app.controller('BreathCtrl', function($scope, $interval) {
  $scope.name = "syam";
  $scope.perMinute = 6;
  // Is +ve when active, zero when inactive.
  $scope.startTimeMs = 0;
  $scope.chestBox = document.getElementById('chest');
  $scope.lungBox = document.getElementById('lung');
  $scope.lungBoxWidth = 500;
  $scope.elapsedMinutes = 0;
  $scope.elapsedSeconds = 0;
  $scope.breatheInSlider = "3000";
  $scope.breatheOutSlider = "3500";

  $scope.nowMs = function() { return (new Date()).getTime(); }

  $scope.cycleTimeMs = function() {
    console.log($scope.breatheInTimeMs());
    console.log($scope.breatheOutTimeMs());
    return $scope.breatheInTimeMs() + $scope.breatheOutTimeMs();
  }

  $scope.breatheInTimeMs = function() {
    return Math.floor($scope.breatheInSlider);
  }

  $scope.breatheOutTimeMs = function() {
    return Math.floor($scope.breatheOutSlider);
  }
  
  $scope.startDrawing = function() {
    if ($scope.startTimeMs > 0) {
      console.log('Internal error! breathing already in progress!');
    }
    $scope.startTimeMs = $scope.nowMs();
    $scope.redrawLoop = $interval($scope.redraw, REDRAW_INTERVAL_MS, 0);
    $scope.startBreatheOut = Math.floor($scope.cycleTimeMs() * $scope.breatheInFraction());
  } 

  $scope.breatheInFraction = function() {
    return $scope.breatheInTimeMs() * 1.0 / ($scope.breatheInTimeMs() + $scope.breatheOutTimeMs());
  }

  $scope.stopDrawing = function() {
    // This should be safe in any case.
    $scope.startTimeMs = 0;
    $interval.cancel($scope.redrawLoop);
  }

  $scope.px = function(width, total, part) {
    return Math.floor(width * part / total) + 'px';
  }

  $scope.lungWidth = function() {
    // set the width of the box, box.style.width to be portion of time done in breathein
    // and breathe out.
    if ($scope.timeInCycleMs > $scope.startBreatheOut) {
      // basically, set the width based on time spent in breathOut. At the beginning,
      // when timeInCycle == startBreatheOut, then lung is 100% full. As it reaches cycleTimeMs(),
      // it becomes 0. Hence the following forumla.
      var total = $scope.cycleTimeMs() - $scope.startBreatheOut;
      var part = $scope.cycleTimeMs() - $scope.timeInCycleMs;
      return $scope.px($scope.lungBoxWidth, total, part);
    } else {
      $scope.lungBox.style.width = $scope.px($scope.lungBoxWidth, $scope.startBreatheOut, $scope.timeInCycleMs);
    }
  }

  $scope.redraw = function() {
    if ($scope.startTimeMs == 0) {
      console.log('Should have stopped!');
      return;
    }
    var elapsedTimeMs = $scope.nowMs() - $scope.startTimeMs;
    $scope.timeInCycleMs = elapsedTimeMs % $scope.cycleTimeMs();
    $scope.elapsedMinutes = Math.floor(elapsedTimeMs / 1000 / 60);
    $scope.elapsedSeconds = Math.floor(elapsedTimeMs / 1000) % 60;
  }
});

var ZEROS = '0000000000000000000000000000';

app.filter('paddingLeft', function() {
  return function(n, len) {
    var nStr = n + '';
    if (nStr.length < len) {
      nStr = ZEROS.substr(0, len - nStr.length) + nStr;
    }
    return nStr;
  }
});
