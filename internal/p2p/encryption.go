package p2p

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/sirupsen/logrus"
)

type Encryption interface {
	GenerateECDHKeyPair() error
}

type KeyData struct {
	PrivateKey *ecdh.PrivateKey
	PublicKey  *ecdh.PublicKey
}

type DefaultEncryption struct{}

func GenerateECDHKeyPair() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
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

func GenerateIdentityKey() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	// ed25519 key generation
	idPubKey, idPrivKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}

	return idPrivKey, idPubKey, nil
}

func GenerateSharedSecret(selfPrivKey *ecdh.PrivateKey, peerPubKey *ecdh.PublicKey) ([]byte, error) {
	sharedSecret, err := selfPrivKey.ECDH(peerPubKey)
	if err != nil {
		logrus.Errorf("error generating shared secret %s", err.Error())
		return nil, err
	}

	return sharedSecret, nil
}

func GenerateSymmetricKeyFromSharedSecret(sharedSecret []byte) [32]byte {
	symmetricKey := sha256.Sum256(sharedSecret)
	return symmetricKey
}

func EncryptMessage(message string, symmetricKey [32]byte) ([]byte, []byte, error) {
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

func DecryptMessage(ciphertext []byte, symmetricKey [32]byte, nonce []byte) (string, error) {
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
