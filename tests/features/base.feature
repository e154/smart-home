# file: version.feature
Feature: get version
  In order to know smart-home version
  I need to be able to request version

  Scenario: cleaning of the base
    When database is clean

  Scenario: does not allow GET method
    When I send "GET" request to "/"
    Then the response code should be 401

  Scenario: does not allow GET method
    When I send "GET" request to "/root"
    Then the response code should be 401

  Scenario: does not allow GET method
    When I send "GET" request to "/api"
    Then the response code should be 401

  Scenario: does not allow GET method
    When I send "GET" request to "/api/v1"
    Then the response code should be 401

  Scenario: does not allow GET method
    When I send "GET" request to "/api/v1/"
    Then the response code should be 401

