Feature: Testing the Index page
    Scenario: Load the Index page
        Given "/" is loaded
        Then the element "root" has text "Index page" 
