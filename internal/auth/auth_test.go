package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "123"

	for _, hasher := range []Hasher{
		NewBcryptHasher(4),
		NewArgon2Hasher(),
	} {
		hash, err := hasher.Hash(password)
		if err != nil {
			t.Fatalf("Hash: %v", err)
		}

		got, err := hasher.Check(password, hash)
		if err != nil {
			t.Fatalf("Check: %v", err)
		}
		if !got {
			t.Errorf("got: %v - want: %v", got, true)
		}

		got, err = hasher.Check("wrong", hash)
		if err != nil {
			t.Fatalf("Check wrong: %v", err)
		}
		if got {
			t.Errorf("got: %v - want: %v", got, false)
		}
	}
}
