import os
import base64

key = os.urandom(32)
print(base64.b64encode(key).decode('utf-8'))