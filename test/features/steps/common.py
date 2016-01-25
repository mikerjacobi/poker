from behave import *
from selenium import webdriver
import dpxdt
from dpxdt.tools import diff_my_images as pdiff
from dpxdt import gflags

use_step_matcher("re")

base_url = "http://jacobra.com:8004"
dpxdt_url = "http://jacobra.com:8015"

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

@step("pdiff the \"(.*)\" page for deploy \"(.*)\"")
def pdiff_page(context, route, deploy_name):
   #../../../bslatin/dpxdt/dpxdt/tools/diff_my_images.py  --upload_build_id=1 --release_server_prefix=http://jacobra.com:8015/api --upload_release_name="My release name" --release_cut_url=http://jacobra.com:8004/ --tests_json_path=features/steps/tests.json 
    context.driver.save_screenshot('/tmp/screenshot.png');
    upload_release_name = deploy_name 
    gflags.FLAGS.verbose = False
    gflags.FLAGS.release_cut_url = base_url+route
    gflags.FLAGS.release_server_prefix = dpxdt_url + "/api"
    gflags.FLAGS.tests_json_path = "features/steps/tests.json"
    gflags.FLAGS.upload_build_id = "2"
    pdiff.real_main(
        release_url=gflags.FLAGS.release_cut_url,
        tests_json_path=gflags.FLAGS.tests_json_path,
        upload_build_id=gflags.FLAGS.upload_build_id,
        upload_release_name=upload_release_name
    )

