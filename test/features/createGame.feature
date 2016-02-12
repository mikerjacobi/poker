Feature: Testing Create Game
  Scenario: Successfully Create Game
    Given there are no games
    And user1 logs in
    And we wait .5 seconds
    And user has a session cookie
    When user navigates to lobby
    And user creates highcard game
    And we wait .5 seconds
    And lobby is screenshot

