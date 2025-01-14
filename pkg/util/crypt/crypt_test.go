package crypt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/keyderivation"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/prf"
)

func TestDerivedKeyRotation(t *testing.T) {
	message, adata := []byte("secret"), []byte(t.Name())
	u, _ := uuid.NewV7()
	salt := u[:]
	var chipertext, plaintext []byte

	template, err := keyderivation.CreatePRFBasedKeyTemplate(prf.HKDFSHA256PRFKeyTemplate(), aead.AES128GCMKeyTemplate())
	require.NoError(t, err, "should create prf based key template")
	mgr := keyset.NewManager()
	rotateKey := func() *keyset.Handle {
		id, err := mgr.Add(template)
		require.NoError(t, err, "should add template")

		mgr.SetPrimary(id)
		t.Logf("new handle with id '%d' is set as primary", id)

		h, err := mgr.Handle()
		require.NoError(t, err, "should return handle")

		return h
	}

	t.Run("encrypt", func(t *testing.T) {
		m := DerivableKeyset[PrimitiveAEAD]{
			master:      rotateKey(),
			constructur: NewPrimitiveAEAD,
		}
		aead, err := m.GetPrimitive(salt)
		require.NoError(t, err, "should return aead primitive")

		chipertext, err = aead.Encrypt(message, adata)
		require.NoError(t, err, "should encrypt original message")
	})

	t.Run("decrypt", func(t *testing.T) {
		m := DerivableKeyset[PrimitiveAEAD]{
			master:      rotateKey(),
			constructur: NewPrimitiveAEAD,
		}
		aead, err := m.GetPrimitive(salt)
		require.NoError(t, err, "should return aead primitive")

		plaintext, err = aead.Decrypt(chipertext, adata)
		require.NoError(t, err, "should decrypt chipertext")

	})

	assert.Equal(t, message, plaintext, "decrypted message should be equal to original message")
}
