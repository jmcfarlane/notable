# Python imports
import hashlib
import random

# Third party imports
from Crypto.Cipher import AES

BLOCKS = 16
MODE = AES.MODE_CBC
PAD = chr(0)

def pad(string):
    return string + (BLOCKS - len(string) % BLOCKS) * PAD

def encrypt(string, pwd):
    iv = ''.join(chr(random.randint(0, 255)) for i in range(BLOCKS))
    return iv + AES.new(key(pwd), MODE, iv).encrypt(pad(string))

def key(pwd):
    return hashlib.sha256(pwd).digest()

def decrypt(cipher, pwd):
    iv = cipher[:BLOCKS]
    return AES.new(key(pwd), MODE, iv).decrypt(cipher[BLOCKS:]).rstrip(PAD)

def main():
    pwd = 'my secret password'
    s = 'I love }apples{'
    encrypted = encrypt(s, pwd)
    print decrypt(encrypted, pwd)

if __name__ == '__main__':
    main()
