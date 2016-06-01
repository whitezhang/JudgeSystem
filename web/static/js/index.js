var app = angular.module('myApp', []);
app.controller('customersCtrl',
    function($scope, $http) {
        var pathname = window.location.pathname.substring(1);
        var url = window.location.search;

        switch (pathname) {
            case "contest":
                var cid = url.substring(url.lastIndexOf("=") + 1, url.length);
                $http.get("http://localhost:8090/scontest", {
                    params: {
                        cid: cid,
                        ipaddr: "defalut"
                    }
                }).success(function(response) {
                    $scope.problemlist = response.problemlist;
                });
                break;
            case "problem":
                var pid = url.substring(url.lastIndexOf("=") + 1, url.length);
                $http.get("http://localhost:8090/sproblem", {
                    params: {
                        pid: pid,
                        ipaddr: "defalut"
                    }
                }).success(function(response) {
                    // $scope.problemlist = response.problemlist;
                });
                break;
        }
    });