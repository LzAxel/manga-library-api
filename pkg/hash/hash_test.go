package hash_test

import (
	"manga-library/pkg/hash"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash_HashingAndComparing(t *testing.T) {
	t.Run("compare same passwords", func(t *testing.T) {
		salt := hash.GenerateSalt()
		passwordHash := hash.HashPassword(salt, "testPassword")
		require.True(t, hash.ComparePassword(passwordHash, "testPassword"))
	})

	t.Run("compare different passwords", func(t *testing.T) {
		salt := hash.GenerateSalt()
		passwordHash := hash.HashPassword(salt, "testPassword")
		require.False(t, hash.ComparePassword(passwordHash, "fakePassword"))
	})
}
