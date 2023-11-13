import os
import base64
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import serialization

class RpcEncryptor:
    def __init__(self, symmetric_key):
        self.symmetric_key = symmetric_key

    @classmethod
    def new_encryptor(cls):
        key = os.urandom(32)  # AES-256 bit key
        return cls(key)

    def key_string(self):
        return base64.b64encode(self.symmetric_key).decode('utf-8')

    @classmethod
    def from_key_string(cls, key_string):
        key = base64.b64decode(key_string)
        return cls(key)

    def encrypt(self, data):
        nonce = os.urandom(12)
        cipher = Cipher(algorithms.AES(self.symmetric_key), modes.GCM(nonce), backend=default_backend())
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(data.encode('utf-8')) + encryptor.finalize()
        return base64.b64encode(nonce + ciphertext).decode('utf-8')

    def decrypt(self, encoded_data):
        data = base64.b64decode(encoded_data)
        nonce, ciphertext = data[:12], data[12:]
        cipher = Cipher(algorithms.AES(self.symmetric_key), modes.GCM(nonce), backend=default_backend())
        decryptor = cipher.decryptor()
        plaintext = decryptor.update(ciphertext) + decryptor.finalize()
        return plaintext.decode('utf-8')