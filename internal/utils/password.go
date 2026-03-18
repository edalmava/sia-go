package utils

import (
	"github.com/alexedwards/argon2id"
)

var defaultParams = &argon2id.Params{
	Memory:      65536,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, defaultParams)
}

func VerifyPassword(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return match && err == nil
}
