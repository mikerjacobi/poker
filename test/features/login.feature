Feature: Testing Login
  Scenario: Successfully Login
    Given user navigates to login
    And login is screenshot
    And user1 logs in
    When we wait .5 seconds
    Then user has a session cookie


