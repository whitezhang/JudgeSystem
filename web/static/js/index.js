var app = angular.module('myApp', []);
app.controller('customersCtrl',
    function($scope, $http) {
        $http.get("http://localhost:8090/sContest?cid=1&&ipaddr=default").success(function(response) {
            $scope.problemlist = response.problemlist;
        });
    });

app.controller('formCtrl',
    function($scope) {

    });