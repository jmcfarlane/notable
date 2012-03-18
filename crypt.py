# Python imports
import hashlib
import random

# Third party imports
from Crypto.Cipher import AES

BLOCK_SIZE = 32
MODE = AES.MODE_CBC
PADDING = '{'

iv = ''.join(chr(random.randint(0, 0xFF)) for i in range(16))
iv = ''.join(chr(0) for i in range(16))
pad = lambda s: s + (BLOCK_SIZE - len(s) % BLOCK_SIZE) * PADDING

def encrypt(string, password):
    aes = AES.new(key(password), MODE, iv)
    return aes.encrypt(pad(string))

def key(password):
    return hashlib.sha256(password).digest()

def decrypt(ciphertext, password):
    aes = AES.new(key(password), MODE, iv)
    return aes.decrypt(ciphertext).rstrip(PADDING)

def main():
    password = 'my secret password'
    s = 'I love }apples{'
    encrypted = encrypt(s, password)
    print decrypt(encrypted, password)

if __name__ == '__main__':
    main()
