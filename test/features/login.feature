Feature: Testing Login
  Scenario: Successfully Login
    Given user1 navigates to login
    And login is screenshot
    When user logs in
    Then user has a session cookie


