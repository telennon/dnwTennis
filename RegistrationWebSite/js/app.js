'use strict';

/********************
// Declare app level module which depends on filters, and services
angular.module('dnwTennisApp', [
  'ngRoute',
  'dnwTennisApp.filters',
  'dnwTennisApp.services',
  'dnwTennisApp.directives',
  'dnwTennisApp.controllers'
]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view1', {templateUrl: 'partials/partial1.html', controller: 'MyCtrl1'});
  $routeProvider.when('/view2', {templateUrl: 'partials/partial2.html', controller: 'MyCtrl2'});
  $routeProvider.otherwise({redirectTo: '/view1'});
}]);

***********************/
var dnwTennisApp = angular.module('dnwTennisApp', ['ngRoute']);

// Configure Routes
dnwTennisApp.config(function($routeProvider) {
	$routeProvider

	// Home Page Route
	.when('/', {
		templateUrl : 'partials/partial1.php',
		controller : 'mainController'
	})

	// Form Page Rout
	.when('/view2', {
		templateUrl : 'partials/partial2.php',
		controller : 'formController'
	})
	.when('/thanksView', {
		templateUrl : 'partials/thanksView.html',
		controller : 'thanksController'
	})
});

dnwTennisApp.controller('mainController', function($scope) {
	$scope.message = "On the main Page"
});

dnwTennisApp.controller('formController', function($scope, $http, $location) {
	$scope.formData = {};
	$scope.processForm = function() {
		$http({
	        method  : 'POST',
	        url     : 'process.php',
	        data    : $.param($scope.formData),  // pass in data as strings
	        headers : { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }  // set the headers so angular passing info as form data (not request payload)
	    }).success(function(data) {
	            //console.log(data);

	            if (!data.success) {
	            	// if not successful, bind errors to error variables
	                $scope.errorlName = true;
	            } else {
	            	// if successful, bind success message to message
	                $scope.message = data.message;
	                //$location.hash('thanksView.html');
	                $location.path('/thanksView')
	            }
	        }).error(function(data) {
	        	$scope.message = data.message;
	        	console.log("Hitting error path");
	        });

		};
		$scope.submitForm = function(isValid) {
			if (isValid) {
				alert('Amazing');
			}
			$scope.submitted = true;
		};

	});

dnwTennisApp.controller('thanksController', function($scope) {
	$scope.message = "On the thanks Page";
});
