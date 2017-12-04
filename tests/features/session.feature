# file: session.feature
Feature: session check

#  Scenario: cleaning of the base
#    When database is clean

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

  Scenario: unsuccessful authorization, incorrect login
    When I authorization with user "test1@e154.ru" and password "testtest"
    Then the response code should be 401
    And the response should match json:
      """
      {
        "message": "Пользователь не найден",
        "status": "error"
      }
      """

  Scenario: unsuccessful authorization, incorrect password
    When I authorization with user "test@e154.ru" and password "testtest123"
    Then the response code should be 403
    And the response should match json:
      """
      {
        "message": "Не верный пароль",
        "status": "error"
      }
      """

  Scenario: successful authorization
    When I authorization with user "test@e154.ru" and password "testtest"
    Then the response code should be 201

  Scenario: does not allow GET method
    When I send "GET" request to "/api/v1/access_list"
    Then the response code should be 200

  Scenario: user finished the session
    When I finishing the session
    Then the response code should be 201

  Scenario: does not allow GET method
    When I send "GET" request to "/api/v1/access_list"
    Then the response code should be 401
