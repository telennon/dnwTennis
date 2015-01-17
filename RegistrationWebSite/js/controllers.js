'use strict';

/* Controllers */

angular.module('dnwTennisApp.controllers', []).
  controller('MyCtrl1', [function($scope) {

  }])
  .controller('MyCtrl2', [function($scope, $http) {
			// create a blank object to hold our form information
			// $scope will allow this to pass between controller and view
			$scope.formData = {};

			// process the form
			$scope.processForm = function() {
				$http({
			        method  : 'POST',
			        url     : 'process.php',
			        data    : $.param($scope.formData),  // pass in data as strings
			        headers : { 'Content-Type': 'application/x-www-form-urlencoded' }  // set the headers so angular passing info as form data (not request payload)
			    })
			        .success(function(data) {
			            console.log(data);

			            if (!data.success) {
			            	// if not successful, bind errors to error variables
			                $scope.errorName = data.errors.name;
			                $scope.errorSuperhero = data.errors.superheroAlias;
			            } else {
			            	// if successful, bind success message to message
			                $scope.message = data.message;
			            }
			        });

			};

  }]);
  /***********************
  var dnwTennisApp = angular.module('dnwTennisApp', []);

  function MyCtrl1 ($scope) {

  };

  function MyCtrl2($scope, $http) {
  	$scope.formData = {};

	// process the form
	$scope.processForm = function() {
		$http({
        	method  : 'POST',
        	url     : 'process.php',
        	data    : $.param($scope.formData),  // pass in data as strings
        	headers : { 'Content-Type': 'application/x-www-form-urlencoded' }  // set the headers so angular passing info as form data (not request payload)
    	})
        .success(function(data) {
            console.log(data);

            if (!data.success) {
            	// if not successful, bind errors to error variables
                $scope.errorName = data.errors.name;
                $scope.errorSuperhero = data.errors.superheroAlias;
            } else {
            	// if successful, bind success message to message
                $scope.message = data.message;
            }
        });
};
};

*************/