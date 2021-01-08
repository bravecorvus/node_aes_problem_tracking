const fs = require('fs');
const crypto = require('crypto');
const path = require('path');

let AESDecrypt = (data, key) => {
  const decoded = Buffer.from(data, 'binary');

  const nonce = decoded.slice(0, 16);
  const ciphertext = decoded.slice(16, decoded.length - 16);
  const tag = decoded.slice(decoded.length - 16);

  let decipher = crypto.createDecipheriv('aes-256-gcm', key, nonce);
  decipher.setAuthTag(tag)
  decipher.setAutoPadding(false);
  try {
    let plaintext = decipher.update(ciphertext, 'binary', 'binary');
    plaintext += decipher.final('binary');
    return Buffer.from(plaintext, 'binary');
  } catch (ex) {
    console.log('AES Decrypt Failed. Exception: ', ex);
    throw ex;
  }
}


if (process.argv.length < 3) {
  console.log('Please pass file path to decrypt');
  process.exit(1);
}
let original = fs.readFileSync(process.argv[2]);
let decrypted = AESDecrypt(original, 'fx6v22kwCjm9oasmMnymhpVJa6H4Xpkc');
fs.writeFileSync('./node-decrypted-' + path.basename(process.argv[2]), decrypted);
