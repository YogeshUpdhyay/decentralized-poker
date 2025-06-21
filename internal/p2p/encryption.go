package p2p

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

type Encryption interface {
	GenerateEphermalKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey, error)
	GenerateSharedSecret(*ecdh.PrivateKey, *ecdh.PublicKey) ([]byte, error)
	GenerateSymmetricKeyFromSharedSecret(sharedSecret []byte) [32]byte
	EncryptMessage(message string, symmetricKey [32]byte) ([]byte, []byte, error)
	DecryptMessage(ciphertext []byte, symmetricKey [32]byte, nonce []byte) (string, error)
}

type DefaultEncryption struct{}

func (d *DefaultEncryption) GenerateEphermalKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
	// generating curve
	curve := ecdh.P256()

	// creating key pair
	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}

	pubKey := privKey.PublicKey()

	logrus.Info("ecdh key pair generated")

	return privKey, pubKey, nil
}

func (d *DefaultEncryption) GenerateSharedSecret(selfPrivKey *ecdh.PrivateKey, peerPubKey *ecdh.PublicKey) ([]byte, error) {
	sharedSecret, err := selfPrivKey.ECDH(peerPubKey)
	if err != nil {
		logrus.Errorf("error generating shared secret %s", err.Error())
		return nil, err
	}

	return sharedSecret, nil
}

func (d *DefaultEncryption) GenerateSymmetricKeyFromSharedSecret(sharedSecret []byte) [32]byte {
	symmetricKey := sha256.Sum256(sharedSecret)
	return symmetricKey
}

func (d *DefaultEncryption) EncryptMessage(message string, symmetricKey [32]byte) ([]byte, []byte, error) {
	plaintext := []byte(message)
	block, err := aes.NewCipher(symmetricKey[:])
	if err != nil {
		logrus.Errorf("error generating block while encrypting %s", err.Error())
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logrus.Errorf("error generating aesgcm %s", err.Error())
		return nil, nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logrus.Errorf("error generating nonce while encrypting %s", err.Error())
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func (d *DefaultEncryption) DecryptMessage(ciphertext []byte, symmetricKey [32]byte, nonce []byte) (string, error) {
	block, err := aes.NewCipher(symmetricKey[:])
	if err != nil {
		logrus.Errorf("error generating block while decrypting %s", err.Error())
		return constants.Empty, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logrus.Errorf("error generating aesgcm %s", err.Error())
		return constants.Empty, err
	}

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logrus.Errorf("error decrypting %s", err.Error())
		return constants.Empty, err
	}

	return string(decrypted), nil
}

func (d *DefaultEncryption) GenerateIdentityKey() (ed25519.PrivateKey, error) {
	// ed25519 key generation
	_, idPrivKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return idPrivKey, nil
}

func (d *DefaultEncryption) EncryptAndSaveIdentityKey(passphrase string, priv ed25519.PrivateKey, filePath string) error {
	// generating salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	// generate key from salt and passphrase
	key, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32) // N=32k, r=8, p=1
	if err != nil {
		return err
	}

	// encrypt the private key using aes gcm
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, priv, nil)

	// store (salt + nonce + ciphertext), all base64-encoded
	combined := append(salt, append(nonce, ciphertext...)...)
	encoded := base64.StdEncoding.EncodeToString(combined)

	return os.WriteFile(filePath, []byte(encoded), 0600)
}

func (d *DefaultEncryption) LoadAndDecryptKey(passphrase, filePath string) (ed25519.PrivateKey, error) {
	// read base64-encoded file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	// extract salt, nonce, ciphertext
	salt := decoded[:16]
	nonceSize := 12 // AES-GCM default
	nonce := decoded[16 : 16+nonceSize]
	ciphertext := decoded[16+nonceSize:]

	// derive key
	key, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// decrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return ed25519.PrivateKey(plaintext), nil
}
