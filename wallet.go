package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/golangcrypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type StorageWallet struct {
	PrivateDer []byte `json:"priv_der"`
	PubKeyByte []byte `json:"pubKey"`
	Curve      string `json:"curve"`
	Address    string `json:"address"`
}

const version = 0
const walletFile = "wallets.dat"

func NewWallet() *Wallet {
	prvKey, pubKey := GenerateKeyPair()

	return &Wallet{
		PrivateKey: prvKey,
		PublicKey:  pubKey,
	}
}

func GenerateKeyPair() (ecdsa.PrivateKey, []byte) {
	var pubKey []byte

	ec := elliptic.P256()
	prvKey, err := ecdsa.GenerateKey(ec, rand.Reader)
	if err != nil {
		log.Panicln("Error generating private key")
	}

	pubKey = append(prvKey.X.Bytes(), prvKey.Y.Bytes()...)
	return *prvKey, pubKey
}

func (wa *Wallet) GenerateAddress() []byte {
	pubKeyHash := HashPubKey(wa.PublicKey)

	versionPayload := append([]byte{version}, pubKeyHash...)

	hash := sha256.Sum256(versionPayload)
	resultHash := sha256.Sum256(hash[:])
	addrCheckSum := resultHash[:4]

	binaryAddr := append(versionPayload, addrCheckSum...)
	base58Addr := base58.Encode(binaryAddr)

	return []byte(base58Addr)
}

func HashPubKey(pubKey []byte) []byte {
	shaHasher := sha256.New()
	_, err := shaHasher.Write(pubKey)
	if err != nil {
		log.Panicln("Error writing in sha function: ", err)
	}
	data := shaHasher.Sum(nil)

	ripemdHasher := ripemd160.New()
	_, err = ripemdHasher.Write(data)
	if err != nil {
		log.Panicln("Error writing in ripemd function: ", err)
	}
	pubKeyHash := ripemdHasher.Sum(nil)

	return pubKeyHash
}
