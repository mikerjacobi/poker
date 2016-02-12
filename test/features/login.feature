Feature: Testing Login
  Scenario: Successfully Login
    Given user1 navigates to login
    And login is screenshot
    And user1 logs in with cli1
    Then user1 has a session cookie


