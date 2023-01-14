package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

var Iv = []byte("532b6195636c6127")[:aes.BlockSize]
var Key = []byte("532b6195636c61279a010000")

func EncryptAES(dst, src, key, iv []byte) error {
	aesBlockEncryptor, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncryptor, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

func DecryptAES(dst, src, key, iv []byte) error {
	aesBlockEncryptor, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesDecrypter := cipher.NewCFBDecrypter(aesBlockEncryptor, iv)
	aesDecrypter.XORKeyStream(src, dst)
	return nil
}
