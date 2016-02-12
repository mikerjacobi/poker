Feature: Testing Join Game
  Scenario: Successfully Join Game
    Given there are no games
    And user1 logs in with cli1
    And there is a highcard game
    And user1 has a session cookie
    When user1 navigates to lobby
    And user1 joins game
    And highcard is screenshot


