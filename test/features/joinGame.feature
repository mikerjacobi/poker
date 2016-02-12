Feature: Testing Join Game
  Scenario: Successfully Join Game
    Given user1 logs in
    And there are no games
    And there is a highcard game
    And we wait .5 seconds
    And user has a session cookie
    When user navigates to lobby
    And user joins game
    And we wait .5 seconds
    And highcard is screenshot


