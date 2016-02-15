Feature: Testing Login
  Scenario: Successfully Login
    Given user1 logs in
    When we wait .5 seconds
    And login is screenshot
    Then user has a session cookie


