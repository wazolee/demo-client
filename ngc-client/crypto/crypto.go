// crypto/crypto.go
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// DecryptFile AES-256-GCM-mel visszafejt egy fájlt hex kulccsal
func DecryptFile(filepath, hexKey string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("érvénytelen kulcs formátum: %v", err)
	}
	if len(key) != 32 {
		return "", fmt.Errorf("a kulcsnak 32 bájtónak kell lennie (64 hex karakter)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("túl rövid titkosított adat")
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("visszafejtési hiba (rossz kulcs?): %v", err)
	}

	return string(plaintext), nil
}

// EncryptFile – ha egyszer titkosítani akarsz (teszthez)
func EncryptFile(plaintext, hexKey string) (string, error) {
	key, _ := hex.DecodeString(hexKey)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}