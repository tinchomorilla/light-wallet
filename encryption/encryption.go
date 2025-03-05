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

// DeriveKey generates a key using PBKDF2 with the given password.
// salt_for_decrypting must be nil when we are encrypting.
// If we are encrypting, it generates a random salt.
// If we are decrypting, it uses the provided salt.
func deriveKey(password string, salt_for_decrypting []byte) ([]byte, []byte, error) {

	if salt_for_decrypting == nil {
		salt := make([]byte, 16)  
		_, err := rand.Read(salt) 
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate random salt: %v", err)
		}
		key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
		return key, salt, nil
	}

	key := pbkdf2.Key([]byte(password), salt_for_decrypting, 100000, 32, sha256.New)

	return key, nil, nil
}

// EncryptPrivateKey encrypts the private key using AES and a password-derived key.
func EncryptPrivateKey(privateKey *ecdsa.PrivateKey, password string) ([]byte, []byte, error) {
	
	key, salt, err := deriveKey(password, nil)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive key: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(privateKeyBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], privateKeyBytes)

	return ciphertext, salt, nil
}

// DecryptPrivateKey decrypts the private key using AES 
// It receives the encrypted private key, the password, and the salt used to derive the key
func DecryptPrivateKey(encryptedPrivateKey []byte, password string, salt []byte) (*ecdsa.PrivateKey, error) {
	// Derive the key using the password and salt
	key, _, err := deriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %v", err)
	}

	if len(encryptedPrivateKey) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract the IV from the encrypted private key (first BlockSize bytes)
	iv := encryptedPrivateKey[:aes.BlockSize]
	encryptedPrivateKey = encryptedPrivateKey[aes.BlockSize:]

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Initialize the AES decryption stream
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encryptedPrivateKey, encryptedPrivateKey)

	// Convert the decrypted bytes back to an ECDSA private key
	privateKey, err := crypto.ToECDSA(encryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to private key: %v", err)
	}

	return privateKey, nil
}

