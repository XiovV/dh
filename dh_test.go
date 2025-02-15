package dh

import (
	"math/big"
	"reflect"
	"testing"
)

func computeDh(aliceGroup, bobGroup int) ([]byte, []byte) {
	alice := NewKeyPair(aliceGroup)
	bob := NewKeyPair(bobGroup)
	aliceSecret := alice.ComputeSecret(bob.PublicKey)
	bobSecret := bob.ComputeSecret(alice.PublicKey)
	return aliceSecret, bobSecret
}

func TestNew(t *testing.T) {
	alice := NewKeyPair()
	bob := NewKeyPair()

	aliceSecret := alice.ComputeSecret(bob.PublicKey)
	bobSecret := bob.ComputeSecret(alice.PublicKey)

	if !reflect.DeepEqual(bobSecret, aliceSecret) {
		t.Errorf("Secrets do not match")
	}

	alice = New(alice.PublicKey, alice.PrivateKey)
	bob = New(bob.PublicKey, bob.PrivateKey)

	aliceSecret = alice.ComputeSecret(bob.PublicKey)
	bobSecret = bob.ComputeSecret(alice.PublicKey)

	if !reflect.DeepEqual(bobSecret, aliceSecret) {
		t.Errorf("Secrets do not match")
	}
}

func TestSameDefaultGroupDh(t *testing.T) {
	aliceSecret, bobSecret := computeDh(14, 14)
	if !reflect.DeepEqual(bobSecret, aliceSecret) {
		t.Errorf("Secrets do not match")
	}
}

func TestSameGroupDh(t *testing.T) {
	aliceSecret, bobSecret := computeDh(16, 16)
	if !reflect.DeepEqual(bobSecret, aliceSecret) {
		t.Errorf("Secrets do not match")
	}
}

func TestDifferentGroupDh(t *testing.T) {
	aliceSecret, bobSecret := computeDh(14, 15)
	if reflect.DeepEqual(bobSecret, aliceSecret) {
		t.Errorf("Secrets can not match")
	}
}

func TestSubsequentDh(t *testing.T) {
	alice := NewKeyPair()
	bob := NewKeyPair()
	eve := NewKeyPair()

	if !reflect.DeepEqual(alice.ComputeSecret(bob.PublicKey), bob.ComputeSecret(alice.PublicKey)) {
		t.Errorf("Secrets do not match for alice-bob")
	}

	if !reflect.DeepEqual(alice.ComputeSecret(eve.PublicKey), eve.ComputeSecret(alice.PublicKey)) {
		t.Errorf("Secrets do not match for alice-eve")
	}

	if !reflect.DeepEqual(bob.ComputeSecret(eve.PublicKey), eve.ComputeSecret(bob.PublicKey)) {
		t.Errorf("Secrets do not match for bob-eve")
	}
}

func TestDiffieHellman_ComputeSecret(t *testing.T) {
	privateKey := big.NewInt(2)
	publicKey := new(big.Int).Exp(g, privateKey, modp2048pInt)

	dh := NewKeyPair() //default group with modp2048pInt
	secret1 := dh.ComputeSecret(publicKey)
	secret2 := new(big.Int).Exp(dh.PublicKey, privateKey, modp2048pInt).Bytes()

	if !reflect.DeepEqual(secret1, secret2) {
		t.Errorf("Expected secret key to be (%d), got (%d)", secret2, secret1)
	}
}

func TestModpLengths(t *testing.T) {
	modpLengths := make(map[int][]byte)
	modpLengths[192] = modp1536pBytes[:]
	modpLengths[256] = modp2048pBytes[:]
	modpLengths[384] = modp3072pBytes[:]
	modpLengths[512] = modp4096pBytes[:]
	modpLengths[768] = modp6144pBytes[:]
	modpLengths[1024] = modp8192pBytes[:]

	for length, modp := range modpLengths {
		if len(modp) != length {
			t.Errorf("Expected modp length to be (%d), got (%d)", length, len(modp))
		}
	}
}

func TestInvalidGroup(t *testing.T) {
	defer func() {
		if i := recover(); i == nil {
			t.Errorf("Invalid group not causing error")
		}
	}()

	NewKeyPair(1)
	NewKeyPair(9)
	NewKeyPair(89)
	NewKeyPair(1004)
}
