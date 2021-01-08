# Node AES GCM Tag Problem

I am unable to use `crypto` library in Node to extract the correct authentication tag when doing AES encryption in Node.js.

The repo contains 2 parts, a Go server which serves a web based interface which does AES encryption and decryption properly on files (hosted at [https://aes.voiceit.io](https://aes.voiceit.io)). And my attempts to write the encryption 
