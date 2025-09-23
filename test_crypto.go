package main

import (
	"fmt"
	"log"

	"github.com/hantdev/chorus-controller/internal/crypto"
)

func main() {
	// Test encryption/decryption
	key := "test-encryption-key-32-chars-long"
	crypto, err := crypto.New(key)
	if err != nil {
		log.Fatal("Failed to create crypto:", err)
	}

	// Test data
	originalSecret := "AKIAIOSFODNN7EXAMPLE"

	fmt.Println("Original Secret:", originalSecret)

	// Encrypt
	encrypted, err := crypto.Encrypt(originalSecret)
	if err != nil {
		log.Fatal("Failed to encrypt:", err)
	}
	fmt.Println("Encrypted:", encrypted)

	// Decrypt
	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		log.Fatal("Failed to decrypt:", err)
	}
	fmt.Println("Decrypted:", decrypted)

	// Verify
	if originalSecret == decrypted {
		fmt.Println("✅ Encryption/Decryption test passed!")
	} else {
		fmt.Println("❌ Encryption/Decryption test failed!")
	}
}
