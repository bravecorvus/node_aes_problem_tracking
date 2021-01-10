const fs = require('fs');
const crypto = require('crypto');
const path = require('path');

let AESEncrypt = (data, key) => {

  const nonce = 'BfVsfgErXsbfiA00';
  const encoded = Buffer.from(data, 'binary');

  const cipher = crypto.createCipheriv('aes-256-gcm', key, nonce);
  try {
    let encrypted = nonce;
    encrypted += cipher.update(encoded, 'binary', 'binary')
    encrypted += cipher.final('binary');
    const tag = cipher.getAuthTag().toString('binary');
    encrypted += tag;
    return Buffer.from(encrypted, 'binary');
  } catch (ex) {
    console.log('AES Encrypt Failed. Exception: ', ex);
    throw ex;
  }
}


if (process.argv.length < 3) {
  console.log('Please pass file path to encrypt');
  process.exit(1);
}
let original = fs.readFileSync(process.argv[2]);
let encrypted = AESEncrypt(original, 'fx6v22kwCjm9oasmMnymhpVJa6H4Xpkc');
fs.writeFileSync('./node-encrypted-' + path.basename(process.argv[2]), encrypted);
