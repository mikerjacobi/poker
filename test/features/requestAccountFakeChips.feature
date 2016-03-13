Feature: Testing Request Fake Chips
    Scenario: Request Fake Chips
        Given user1 has 100 account balance
        And user1 logs in with cli1
        And user1 has a session cookie
        When user1 navigates to account
        And user1 requests 50 chips
        And accounts is screenshot
        Then user1 should have 150 chips in their account
