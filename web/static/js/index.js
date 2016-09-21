var app = angular.module('myApp', []);

app.controller('formCtrl',
    function($scope, $http) {
        $scope.signin = function() {
            $http.get("http://114.215.125.77:8090/slogin", {
                params: {
                    username: $scope.username,
                    password: $scope.password
                }
            }).success(function(response) {
                if (response.Status == "OK") {
                    //alert("Yes"); // TODO
                }
            });
        };
        $scope.signin();
    });

app.controller('customersCtrl',
    function($scope, $http) {
        var pathname = window.location.pathname.substring(1);
        var url = window.location.search;
        switch (pathname) {
            case "contests":
                var pid = url.substring(url.lastIndexOf("=") + 1, url.length);
                $http.get("http://114.215.125.77:8090/scontests", {
                    params: {
                        pid: pid,
                        ipaddr: "defalut"
                    }
                }).success(function(response) {
                    $scope.contestname = response.contestname;
                    $scope.starttime = response.starttime;
                    $scope.endtime = response.endtime;
                });
                break;
            case "contestinfo":
                var cid = url.substring(url.lastIndexOf("=") + 1, url.length);
                $http.get("http://114.215.125.77:8090/scontestinfo", {
                    params: {
                        cid: cid,
                        ipaddr: "defalut"
                    }
                }).success(function(response) {
                    $scope.problemlist = response;
                });
                break;
            case "problems":
                var page = url.substring(url.lastIndexOf("=") + 1, url.length);
                $http.get("http://114.215.125.77:8090/sproblems", {
                    params: {
                        page: page,
                        ipaddr: "defalut"
                    }
                }).success(function(response) {
                    $scope.problemlist = response;
                });
                break;
            case "problem":
                var pid = url.substring(url.lastIndexOf("=") + 1, url.length);
                $http.get("http://114.215.125.77:8090/sprobleminfo", {
                    params: {
                        pid: pid,
                        ipaddr: "defalut"
                    }
                }).success(function(response) {
                    $scope.PID = response.PID;
                    $scope.title = response.title;
                    $scope.description = response.description;
                    $scope.input = response.input;
                    $scope.output = response.output;
                    $scope.time = response.time;
                    $scope.memory = response.memory;
                    $scope.simpleinput = response.simpleinput;
                    $scope.simpleoutput = response.simpleoutput;
                });
                break;
        }
    });
