package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func AESDecryptBytes(payload []byte) ([]byte, error) {
	iv, cipherText := payload[:16], payload[16:]

	block, err := aes.NewCipher([]byte("fx6v22kwCjm9oasmMnymhpVJa6H4Xpkc"))
	if err != nil {
		return []byte{}, err
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return []byte{}, err
	}

	plaintext, err := aesgcm.Open(nil, iv, cipherText, nil)
	if err != nil {
		return []byte{}, err
	}

	return plaintext, nil
}

func AESEncryptBytes(payload []byte) ([]byte, error) {

	iv := []byte("BfVsfgErXsbfiA00")

	block, err := aes.NewCipher([]byte("fx6v22kwCjm9oasmMnymhpVJa6H4Xpkc"))
	if err != nil {
		return []byte{}, err
	}
	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return []byte{}, err
	}
	ciphertext := aesgcm.Seal(nil, iv, payload, nil)
	return append(iv, ciphertext...), nil
}
