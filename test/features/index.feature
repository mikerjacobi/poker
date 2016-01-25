Feature: Testing the Index page
    Scenario: Load the Index page
        When "/" is loaded
        Then the element "root" has text "Index page2" 
        And pdiff the "/" page for deploy "deploy #2"
