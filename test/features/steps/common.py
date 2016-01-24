from behave import *
from selenium import webdriver

use_step_matcher("re")

base_url = "http://jacobra.com:8004"

@step("\"(.*)\" is loaded")
def load_page(context, route):
    url = base_url + route
    context.driver = webdriver.PhantomJS()
    context.driver.get(url)
    assert url in context.driver.current_url

@step("the element \"(.*)\" has text \"(.*)\"")
def check_element_text(context, element, text):
    elem = context.driver.find_element_by_id(element)
    assert text in elem.get_attribute("innerHTML")
