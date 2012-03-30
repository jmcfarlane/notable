# Python imports
import base64
import hashlib
import random

# Stub in the event pycrypto isn't available
class NoEncryption(object):
    MODE_CBC = None

# Third party imports
try:
    from Crypto.Cipher import AES
except (ImportError, NameError):
    AES = NoEncryption()

BLOCKS = 16 # 16bit IV and 32bit AES padding
MODE = AES.MODE_CBC
PAD = chr(0)

def pad(string):
    _b32 = BLOCKS * 2
    return string + (_b32 - len(string) % _b32) * PAD

def encrypt(string, pwd):
    iv = ''.join(chr(random.randint(0, 255)) for i in range(BLOCKS))
    cipher = iv + AES.new(key(pwd), MODE, iv).encrypt(pad(string))
    return base64.b64encode(cipher)

def key(pwd):
    return hashlib.sha256(pwd).digest()

def decrypt(cipher, pwd):
    cipher = base64.b64decode(cipher)
    iv = cipher[:BLOCKS]
    return AES.new(key(pwd), MODE, iv).decrypt(cipher[BLOCKS:]).rstrip(PAD)

def main():
    pwd = 'my secret password'
    s = 'I love }apples{'
    encrypted = encrypt(s, pwd)
    print(decrypt(encrypted, pwd))

if __name__ == '__main__':
    main()
