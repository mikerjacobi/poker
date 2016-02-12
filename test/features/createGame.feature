Feature: Testing Create Game
  Scenario: Successfully Create Game
    Given there are no games
    And user1 logs in with cli1
    And user1 has a session cookie
    When user1 navigates to lobby
    And user1 creates highcard game
    And lobby is screenshot

