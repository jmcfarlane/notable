# Python imports
import os
import pdb
import sys
import time
import unittest
import uuid

# Third party imports
from selenium import webdriver

# Backport for python2.6
if not hasattr(unittest, 'skipIf'):
    import unittest2 as unittest

def webDriverNotSupported():
    "Determine if webdriver is supported (at runtime)"
    def reasons():
        yield sys.version_info >= (3, 0)
        yield os.environ.get('TRAVIS')
    return [r for r in reasons() if r]

class TestWebApp(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        "Create a web driver browser"
        if not webDriverNotSupported():
            self.b = webdriver.Firefox()

    @classmethod
    def tearDownClass(self):
        "Close the browser when the class test is over"
        if not webDriverNotSupported():
            self.b.close()

    @unittest.skipIf(webDriverNotSupported(), '')
    def setUp(self):
        "Reload the page at the start of every test"
        self._reload()

    def _create_note(self, encrypted=False):
        "Create a note, using encryption if desired"
        self.token = uuid.uuid4().hex
        self.b.find_element_by_class_name('create').click()
        time.sleep(1)
        for field, text in self._note_form_fields(encrypted=encrypted):
            field.send_keys(text + self.token)
        content = 'Hello world: %s' % self.token
        self.b.execute_script("window._editor.setValue('%s')" % content)
        self.b.find_element_by_css_selector('button.save').click()
        time.sleep(0.2)
        saved = self.b.find_element_by_css_selector('.alert.saved')
        self.assertTrue(saved.is_displayed())
        time.sleep(1)

    def _created_note(self):
        "Return a row element of the last created note"
        css = '//td[text()="subject %s"]' % self.token
        return self.b.find_element_by_xpath(css)

    def _note_form_fields(self, encrypted=False):
        "Return a generator of all note form fields"
        inputs = [('subject', 'subject '), ('tags', 'tag ')]
        inputs = inputs + [('password', 'password ')] if encrypted else inputs
        for placeholder, text in inputs:
            css = '.%s input' % placeholder
            yield self.b.find_element_by_css_selector(css), text

    def _reload(self):
        "Reload the page"
        self.b.get('http://localhost:8083')
        time.sleep(0.5)

    @unittest.skipIf(webDriverNotSupported(), '')
    def test_001_create_note(self):
        "Create a new note and validate it's there after page reload"
        self._create_note()
        self._reload()
        self._created_note().click()
        time.sleep(1)
        for field, text in self._note_form_fields():
            self.assertEquals(field.get_attribute('value'), text + self.token)

    @unittest.skipIf(webDriverNotSupported(), '')
    def test_005_create_encrypted_note(self):
        "Create an encrypted note and make sure it can be opened"
        self._create_note(encrypted=True)
        self._reload()
        self._created_note().click()
        time.sleep(1)
        passwd = self.b.find_element_by_id('password')
        passwd.send_keys('password %s' % self.token)
        self.b.find_element_by_link_text('Decrypt').click()
        time.sleep(1)
        for field, text in self._note_form_fields():
            self.assertEquals(field.get_attribute('value'), text + self.token)

if __name__ == '__main__':
    unittest.main()
