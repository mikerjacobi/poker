Feature: Testing Replay HighCard Game
  Scenario: Replay HighCard
    Given there are no games
    And there is a highcard game
    And user1 logs in with cli1
    And user1 has a session cookie
    And user2 logs in with cli2
    And user2 has a session cookie
    When user1 navigates to lobby
    When user2 navigates to lobby
    And user1 joins game
    And user2 joins game
    And user1 replays game
    And user1 screenshots game
    And user2 screenshots game
    Then game screenshots should be equal
