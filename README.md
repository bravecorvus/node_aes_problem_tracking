# Node AES GCM Tag Problem

I am unable to use `crypto` library in Node to extract the correct authentication tag when doing 256-bit AES GCM encryption in Node.js.

The repo contains 2 parts:

1. [go_aes_server](go_aes_server) a Go based server which hosts a web based interface which does AES encryption and decryption properly on files (hosted at [https://go-aes.voiceit.io](https://go-aes.voiceit.io)). And my attempts to write the encryption 
1. [node_aes](node_aes) My best attempt to implement a Node.js version of the 256-bit AES GCM encyption and decryption.

Because of how the higher level languages Go, and Java does the auth tag in their respective standard library AES GCM implementation, and Node needs to be able to decrypt data written in Go/Java and vice versa, I cannot change the following format of the encrypted form of the file:

```
| Nonce/IV (First 16 bytes) | Ciphertext | Authentication Tag (Last 16 bytes) |
```

Also, the tag coming at the end follows the standard specification of encrypted payloads specified in RFC 5116 section 2.1 ([source](https://crypto.stackexchange.com/questions/25249/where-is-the-authentication-tag-stored-in-file-encrypted-using-aes-gcm)). So I see it as fine for our purposes

Because I my company requires a common AES encryption standard in multiple languages, I have done the following implementations:

| *Language* | *Encryption* | *Decryption* |
| -- | -- | -- |
| Go | ✅ | ✅  |
| Java | ✅ | ✅  |
| Node.js | ❌ | ✅  |

Now the main reason why I know the problem is must be in Node encryption is because of the following set of circumstances:

Go and Java have bidirectional encryption/decryption capabilities (of course provided I am using the same encryption key). Furthermore, [node_aes/decrypt.js](node_aes/decrypt.js) is able to decrypt files that were encrypted in Go and Java.

To make it easier for people to test this themselves, I have written a Go based encryption/decryption web interface at [https://go-aes.voiceit.io](https://go-aes.voiceit.io). This will allow you to upload a file of your choice, and download either the encrypted form, or decrypted form of that same file. Please use this interface to encrypt a file, then decrypt that same file using [node_aes/decrypt.js](node_aes/decrypt.js):

```
cd node_aes
node decrypt.js ./go-encrypted-file
[will produce ./node-decrypted-go-encrypted-file]
```

---

Furthermore, a file encrypted by Node.js will fail to decrypt using Node.js (while Go and Java will also fail to decrypt that file):

```
cd node_aes
node encrypt.js ./file
node decrypt.js ./node-encrypted-file
[will fail]
```

File diffing a Go encrypted file against the same file encrypted in Node using [vbindiff](https://www.cjmweb.net/vbindiff/), the file is identical until the last 16 bytes where the auth tag gets written.

![vbindiff](https://drive.voiceit.io/files/vbindiff.png)

Any help would be much appreciated. Thank you.

## *WARNING*

  > If you find the code in my repo useful, and want to use this in your stuff, make sure you use a secure psudo-random character generator to produce a new IV/Nonce before putting these encryption methods in your production environments. I have hardcoded the nonce this in this repo to make it easier to run a binary file diffs. But in production, you need to generate a new, randomized nonce for every new thing you encrypt. All files encrypted using the same nonce breaks the security model of AES encryption because of [this](https://crypto.stackexchange.com/questions/26790/how-bad-it-is-using-the-same-iv-twice-with-aes-gcm).
