package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

// EncryptUserKey 使用主密钥对用户密钥进行加密
func EncryptUserKey(key []byte) error {
	masterKey, err := hex.DecodeString(os.Getenv("MASTER_KEY"))
	if err != nil {
		return err
	}

	if len(key) != 16 {
		return errors.New(fmt.Sprintf("invalid key size %d", len(key)))
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return err
	}

	block.Encrypt(key, key)
	return nil
}

// DecryptUserKey 使用主密钥对加密的用户密钥进行解密
func DecryptUserKey(key []byte) error {
	masterKey, err := hex.DecodeString(os.Getenv("MASTER_KEY"))
	if err != nil {
		return err
	}

	if len(key) != 16 {
		return errors.New(fmt.Sprintf("invalid key size %d", len(key)))
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return err
	}

	block.Decrypt(key, key)
	return nil
}

// GenerateRandomBytes 生成指定长度的密钥
func GenerateRandomBytes(length int) (key []byte) {
	key = make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		Log().Error("generate key error", err)
	}

	return key
}

func MD5(data []byte) string {
	digest := md5.Sum(data)
	return hex.EncodeToString(digest[:])
}

// EncryptFile 用给定的密钥对文佳内容进行加密
func EncryptFile(content []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 生成iv
	iv := GenerateRandomBytes(block.BlockSize())
	// 加密运算
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(content, content)
	// 返回 iv + ciphertext
	return append(iv, content...), nil
}

// DecryptFile 用给定的密钥对文佳内容进行加密
func DecryptFile(content []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	size := block.BlockSize()
	if len(content) > size {
		// 提取 iv, ciphertext
		iv := content[:size]
		ciphertext := content[size:]
		// 解密运算
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(ciphertext, ciphertext)
		return ciphertext, nil
	}
	return nil, errors.New("decryption failed")
}
