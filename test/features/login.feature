Feature: Testing Login
  Scenario: Successfully Login
    Given user1 logs in
    #And login is screenshot
    When we wait 1 seconds
    Then user has a session cookie


