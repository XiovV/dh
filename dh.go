package dh

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

var (
	defaultGroup = 14
	g            = big.NewInt(2)

	modp1536pInt = new(big.Int).SetBytes(modp1536pBytes)
	modp2048pInt = new(big.Int).SetBytes(modp2048pBytes)
	modp3072pInt = new(big.Int).SetBytes(modp3072pBytes)
	modp4096pInt = new(big.Int).SetBytes(modp4096pBytes)
	modp6144pInt = new(big.Int).SetBytes(modp6144pBytes)
	modp8192pInt = new(big.Int).SetBytes(modp8192pBytes)

	modp1536pKeyLen = len(modp1536pBytes)
	modp2048pKeyLen = len(modp2048pBytes)
	modp3072pKeyLen = len(modp3072pBytes)
	modp4096pKeyLen = len(modp4096pBytes)
	modp6144pKeyLen = len(modp6144pBytes)
	modp8192pKeyLen = len(modp8192pBytes)
)

func getGroupParams(id int) (*big.Int, int) {
	switch id {
	case 5:
		return modp1536pInt, modp1536pKeyLen
	case 14:
		return modp2048pInt, modp2048pKeyLen
	case 15:
		return modp3072pInt, modp3072pKeyLen
	case 16:
		return modp4096pInt, modp4096pKeyLen
	case 17:
		return modp6144pInt, modp6144pKeyLen
	case 18:
		return modp8192pInt, modp8192pKeyLen

	default:
		panic(fmt.Sprintf("Invalid group: %d. Available groups: 5, 14, 15, 16, 17, 18.\n", id))
	}
}

type DiffieHellman struct {
	PublicKey, PrivateKey, p *big.Int
}

func (dh DiffieHellman) ComputeSecret(public *big.Int) []byte {
	return new(big.Int).Exp(public, dh.PrivateKey, dh.p).Bytes()
}

func generatePrivateKey(keyLen int) *big.Int {
	buff := make([]byte, keyLen)
	_, err := io.ReadFull(rand.Reader, buff)
	if err != nil {
		panic(err)
	}
	return new(big.Int).SetBytes(buff)
}

func dhFromGroup(groupId int) DiffieHellman {
	dh := DiffieHellman{}
	p, keyLen := getGroupParams(groupId)
	dh.PrivateKey = generatePrivateKey(keyLen)
	dh.PublicKey = new(big.Int).Exp(g, dh.PrivateKey, p)
	dh.p = p
	return dh
}

func New(public, private *big.Int) DiffieHellman {
	dh := DiffieHellman{}
	p, _ := getGroupParams(defaultGroup)
	dh.PublicKey = public
	dh.PrivateKey = private
	dh.p = p
	return dh
}

func NewKeyPair(group ...int) DiffieHellman {
	if len(group) >= 1 {
		return dhFromGroup(group[0])
	}
	return dhFromGroup(defaultGroup)
}
