package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/mac"
)

func TestBlindIndexes(t *testing.T) {
	mgr := keyset.NewManager()
	hid, err := mgr.Add(mac.HMACSHA256Tag256KeyTemplate())
	require.NoError(t, err, "should add mac handle")
	err = mgr.SetPrimary(hid)
	require.NoError(t, err, "should set primary handle")
	handle, err := mgr.Handle()
	require.NoError(t, err, "should obtain mac handle")
	m, err := mac.New(handle)
	require.NoError(t, err, "should create mac primitive")

	data := []byte("asdasjdiu9lksdlfkjasopfijaposdpasi09ie283u023hj02i0t83089tu045jt054050j")
	v, err := m.ComputeMAC(data[:])
	require.NoError(t, err, "should compute mac")

	hid, err = mgr.Add(mac.HMACSHA256Tag128KeyTemplate())
	require.NoError(t, err, "should add new mac handle")
	err = mgr.SetPrimary(hid)
	require.NoError(t, err, "should set new primary handle")
	handle, err = mgr.Handle()
	require.NoError(t, err, "should obtain new mac handle")

	vs, err := GetBlindIdxs(handle, data[:], len(v))
	require.NoError(t, err, "should compute multiple mac")

	assert.Contains(t, vs, v, "should contain previous mac")
}
