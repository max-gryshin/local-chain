package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenerateKeyEllipticP256() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

func PrivateKeyToBytes(privateKey *ecdsa.PrivateKey) ([]byte, []byte) {
	privDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privDER}), pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
}

func PublicKeyToBytes(public *ecdsa.PublicKey) []byte {
	pubDER, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		panic(err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
}

func PublicKeyFromBytes(pubKey []byte) (*ecdsa.PublicKey, error) {
	outputBlock, _ := pem.Decode(pubKey)
	if outputBlock == nil || outputBlock.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key PEM, receiver: %s, publicBlock: %v", pubKey, outputBlock)
	}
	outputPub, err := x509.ParsePKIXPublicKey(outputBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outputPub: %v", err)
	}
	outputPubKey, ok := outputPub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("outputPubKey key is not ECDSA")
	}
	return outputPubKey, nil
}

func PrivateKeyFromBytes(privateKey []byte) (*ecdsa.PrivateKey, error) {
	privateBlock, _ := pem.Decode(privateKey)
	if privateBlock == nil || privateBlock.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key PEM")
	}
	ecdsaPrivateKey, err := x509.ParseECPrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	return ecdsaPrivateKey, nil
}
