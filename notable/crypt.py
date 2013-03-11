# Python imports
import base64
import hashlib
import random
import sys

# Stub in the event pycrypto isn't available
class NoEncryption(object):
    MODE_CBC = None

# Third party imports
try:
    from Crypto.Cipher import AES
    from Crypto import Random
except (ImportError, NameError):
    AES = NoEncryption()

BLOCKS = 16 # 16bit IV and 32bit AES padding
MODE = AES.MODE_CBC
PAD = chr(0)

def pad(string):
    _b32 = BLOCKS * 2
    return string + (_b32 - len(string) % _b32) * PAD

def encrypt(string, pwd):
    iv = Random.new().read(BLOCKS)
    cipher = iv + AES.new(key(pwd), MODE, iv).encrypt(pad(string))
    return base64.b64encode(cipher)

def key(pwd):
    return hashlib.sha256(pwd.encode('utf-8')).digest()

def decrypt(cipher, pwd):
    cipher = cipher.encode() if hasattr(cipher, 'encode') else cipher
    cipher = base64.b64decode(cipher)
    iv = cipher[:BLOCKS]
    decrypted = AES.new(key(pwd), MODE, iv).decrypt(cipher[BLOCKS:])
    if sys.version_info >= (3, 0):
        return decrypted.decode('utf-8', errors='ignore').rstrip(PAD)
    return decrypted.rstrip(PAD)

def main():
    pwd = 'my secret password'
    s = 'I love }apples{'
    encrypted = encrypt(s, pwd)
    print(decrypt(encrypted, pwd))

if __name__ == '__main__':
    main()
