# Python imports
from distutils.core import setup

# Project imports
from notable import app

# Attributes
AUTHOR = 'John McFarlane'
DESCRIPTION = 'A very simple note taking application'
EMAIL = 'john.mcfarlane@rockfloat.com'
NAME = 'Notable'
PYPI = 'http://pypi.python.org/packages/source/N/Notable'
URL = 'https://github.com/jmcfarlane/Notable'
CLASSIFIERS = """
Development Status :: 2 - Pre-Alpha
Intended Audience :: Developers
License :: OSI Approved :: MIT License
Operating System :: OS Independent
Programming Language :: Python
Topic :: Internet :: WWW/HTTP
Intended Audience :: End Users/Desktop
Topic :: Office/Business :: News/Diary
Topic :: Security :: Cryptography
Topic :: Utilities
"""

setup(
    author = AUTHOR,
    author_email = EMAIL,
    classifiers = [c for c in CLASSIFIERS.split('\n') if c],
    description = DESCRIPTION,
    download_url = '%s/Notable-%s.tar.gz' % (PYPI, app.version),
    name = NAME,
    packages = ['notable'],
    package_data = {'notable': ['static/*', 'views/*']},
    scripts = ['scripts/notable'],
    url = URL,
    version = app.version
)
