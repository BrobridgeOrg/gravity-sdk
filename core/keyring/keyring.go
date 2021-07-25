package keyring

import (
	"sync"

	"github.com/BrobridgeOrg/gravity-sdk/core/encryption"
)

type KeyInfo struct {
	encryption *encryption.Encryption
	permission *Permission
	ref        int
}

func (info *KeyInfo) Encryption() *encryption.Encryption {
	return info.encryption
}

func (info *KeyInfo) Permission() *Permission {
	return info.permission
}

type Keyring struct {
	apps sync.Map
}

func NewKeyring() *Keyring {
	return &Keyring{}
}

func (keyring *Keyring) Put(appID string, key string) *KeyInfo {
	en := encryption.NewEncryption()
	en.SetKey(key)

	info := &KeyInfo{
		encryption: en,
		permission: NewPermission(),
		ref:        1,
	}

	keyring.apps.Store(appID, info)

	return info
}

func (keyring *Keyring) Get(appID string) *KeyInfo {
	v, ok := keyring.apps.Load(appID)
	if !ok {
		return nil
	}

	return v.(*KeyInfo)
}

func (keyring *Keyring) Unref(appID string) {

	keyInfo := keyring.Get(appID)
	if keyInfo == nil {
		return
	}

	if keyInfo.ref == 1 {
		keyring.apps.Delete(appID)
	}
}

func (keyring *Keyring) Delete(appID string) {
	keyring.Delete(appID)
}
