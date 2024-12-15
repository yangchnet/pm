package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/yangchnet/pm/config"
)

var (
	ErrPassword = errors.New("密码错误")
)

func createHash(key string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func Encrypt(key string, text []byte) ([]byte, error) {
	keyHash := createHash(key)

	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(key string, text []byte) (rawPasswd []byte, err error) {
	defer func() {
		userKeyPath := config.GetString("user_key_path")
		if err != nil {
			os.Remove(userKeyPath)
		}
	}()
	keyHash := createHash(key)

	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, ErrPassword
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, ErrPassword
	}
	return data, nil
}
