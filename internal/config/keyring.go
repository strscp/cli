package config

import (
	"github.com/zalando/go-keyring"
)

const keyringService = "starscope-cli"

// Keyring manages token storage in the OS keychain.
type Keyring struct{}

// NewKeyring creates a new Keyring instance.
func NewKeyring() *Keyring {
	return &Keyring{}
}

// Get retrieves a token from the OS keychain.
func (k *Keyring) Get(profile string) (string, error) {
	return keyring.Get(keyringService, profile)
}

// Set stores a token in the OS keychain.
func (k *Keyring) Set(profile, token string) error {
	return keyring.Set(keyringService, profile, token)
}

// Delete removes a token from the OS keychain.
func (k *Keyring) Delete(profile string) error {
	return keyring.Delete(keyringService, profile)
}
