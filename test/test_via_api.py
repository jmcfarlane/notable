# Python imports
import json
import os
import unittest
import uuid

# Third party imports
import requests

class TestWebAPI(unittest.TestCase):

    def _url(self, path):
        "Return an api path"
        return os.path.join('http://localhost:8083/', path.lstrip('/'))

    def _create_note(self, encrypted=False):
        self.token = uuid.uuid4().hex
        note = dict(subject='Subject %s' % self.token,
                    tags='tags %s' % self.token,
                    content='Content %s' % self.token)
        if encrypted:
            note['password'] = 'password %s' % self.token
        r = requests.post(self._url('/api/note/create'), json.dumps(note))
        created = r.json if isinstance(r.json, dict) else r.json()
        self.assertEquals(r.status_code, 200)
        self.assertEquals(created.get('subject'), note.get('subject'))
        self.assertEquals(created.get('tags'), note.get('tags'))
        self.assertEqual(len(created.get('uid')), len(self.token))
        return note, created

    def _get_content(self, created, password):
        url = '/api/note/content/%s' % created.get('uid')
        form = dict(password=password)
        return requests.post(self._url(url), data=form)

    def test_001_create_and_fetch_note(self):
        "api: create and fetch note"
        note, created = self._create_note()
        r = self._get_content(created, '')
        self.assertEquals(r.status_code, 200)
        self.assertEquals(r.text, note.get('content'))

    def test_005_create_and_fetch_encrypted_note(self):
        "api: create and fetch encrypted note"
        note, created = self._create_note(encrypted=True)
        r = self._get_content(created, 'password %s' % self.token)
        self.assertEquals(r.status_code, 200)
        self.assertEquals(r.text, note.get('content'))

    def test_010_create_and_fetch_encrypted_note_with_wrong_passwd(self):
        "api: create and fetch encrypted note with wrong password"
        note, created = self._create_note(encrypted=True)
        r = self._get_content(created, 'This is the wrong password!')
        self.assertEquals(r.status_code, 403)
        self.assertEquals(r.text, 'Nope, try again')

    def test_015_note_listing_structure(self):
        "api: note listing (validating the data schema)"
        keys = ['created', 'encrypted', 'subject', 'tags', 'uid', 'updated']
        r = requests.get(self._url('/api/notes/list'))
        for note in (r.json if isinstance(r.json, list) else r.json()):
            self.assertEquals(type(note), type(dict()))
            self.assertEquals(sorted(note.keys()), keys)

if __name__ == '__main__':
    unittest.main()
