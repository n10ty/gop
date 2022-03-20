package main

import (
	"encoding/base32"
	"os"
	"testing"
)

func Test_readWriteKeychain(t *testing.T) {
	var path = ".test.keychain"
	defer func() {
		os.Remove(path)
	}()
	keychain := &Keychain{
		keys: map[string]string{
			"name1": "key1",
			"name2": "key2",
		},
	}
	err := writeKeychain(keychain, path)
	if err != nil {
		t.Errorf("Error during write keychain: %s", err)
		return
	}

	actualKeychain, err := readKeychain(path)
	if err != nil {
		t.Errorf("Error during read keychain: %s", err)
	}

	if len(actualKeychain.keys) != len(keychain.keys) {
		t.Errorf("Error during test keychain read-write: expected %d, actual %d", len(keychain.keys), len(actualKeychain.keys))
		return
	}

	if actualKeychain.keys["name1"] != keychain.keys["name1"] {
		t.Errorf("Error during test keychain read-write: expected %s, actual %s", keychain.keys["name1"], actualKeychain.keys["name1"])
		return
	}

	if actualKeychain.keys["name2"] != keychain.keys["name2"] {
		t.Errorf("Error during test keychain read-write: expected %s, actual %s", keychain.keys["name2"], actualKeychain.keys["name2"])
		return
	}
}

func Test_addGetKeyFromKeychain(t *testing.T) {
	keychain := Keychain{keys: map[string]string{}}
	key := "L7TXVT75PU52SE3G"
	name := "name"
	err := keychain.add(name, key)
	if err != nil {
		t.Errorf("Error during add key to keychain: %s", err)
		return
	}

	actualKey, err := keychain.get(name)
	if err != nil {
		t.Errorf("Error during get value from keychain: %s", err)
		return
	}

	if key != base32.StdEncoding.EncodeToString(actualKey) {
		t.Errorf("Error durin add-get key from keychain: expected: %s, actual: %s", key, base32.StdEncoding.EncodeToString(actualKey))
	}
}

func Test_hotp(t *testing.T) {
	key, _ := base32.StdEncoding.DecodeString("L7TXVT75PU52SE3G")
	code := hotp(key, 1333589, 6)
	expected := 304098
	if code != expected {
		t.Errorf("Error during hotp: expected %d, actual %d", expected, code)
	}
}
