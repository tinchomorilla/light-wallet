package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/pbkdf2"
)

// DeriveKey generates an encryption key from the password using PBKDF2 with a random salt.
func DeriveKey(password string) ([]byte, error) {

	salt := make([]byte, 16)  // Create a byte slice of length 16
	_, err := rand.Read(salt) // Fill the slice with random data
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %v", err)
	}
	
	fmt.Println("salt:", salt)
	fmt.Println("password:", password)

	// We should store the salt along with the encrypted private key
	// so that we can derive the same key later to decrypt the private key.
	// For now, we are not storing the salt.
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	return key, nil

}

// EncryptPrivateKey encrypts the private key using AES and a password-derived key.
func EncryptPrivateKey(privateKey *ecdsa.PrivateKey, password string) ([]byte, error) {
	key, err := DeriveKey(password)

	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(privateKeyBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], privateKeyBytes)

	return ciphertext, nil
}

// DecryptPrivateKey decrypts the private key using AES and the password-derived key.
func DecryptPrivateKey(encryptedPrivateKey []byte, password string) (*ecdsa.PrivateKey, error) {
	key, err := DeriveKey(password)

	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %v", err)
	}

	if len(encryptedPrivateKey) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := encryptedPrivateKey[:aes.BlockSize]
	encryptedPrivateKey = encryptedPrivateKey[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encryptedPrivateKey, encryptedPrivateKey)

	privateKey, err := crypto.ToECDSA(encryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to private key: %v", err)
	}

	return privateKey, nil
}
