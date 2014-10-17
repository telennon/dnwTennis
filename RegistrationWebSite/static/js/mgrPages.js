  'use strict';

angular.module('mgrPages', ['trNgGrid']);
function MainCtrl($scope) {
	$scope.myItems = [{name: "Maroni", age: 50},
		{name: "Tianicum", age: 43},
		{name: "Jacob", age: 27},
		{name: "Nephi", age: 29},
		{name: "Enos", age: 99}];
};