package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// 解密
func Decrypt(ciphertext string, key, iv []byte) (string, error) {
	// base64解码
	src, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// CBC解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(src))
	blockMode.CryptBlocks(plaintext, src)

	// Unpadding
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	if unpadding > length {
		return "", fmt.Errorf("decrypt error. unpadding(%d) > length(%d)", unpadding, length)
	}
	return string(plaintext[:(length - unpadding)]), nil
}

// 加密
func Encrypt(plaintext, key, iv []byte) (string, error) {
	// 创建block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Padding
	src := []byte(plaintext)
	blockSize := block.BlockSize()
	padding := blockSize - len(src)%blockSize
	padtext := append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padtext))
	blockMode.CryptBlocks(ciphertext, padtext)

	// base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
