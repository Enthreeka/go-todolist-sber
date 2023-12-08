package usecase

import (
	"crypto/subtle"
	"encoding/hex"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"golang.org/x/crypto/argon2"
)

type Argon interface {
	GenerateHashFromPassword(password string) (string, error)
	VerifyPassword(hashPassword string, password string) error
}

type argon struct {
	salt    []byte
	version int
	threads uint8
	time    uint32
	memory  uint32
	keyLen  uint32
}

func NewArgonPassword(salt string) Argon {
	return &argon{
		salt:    []byte(salt),
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
}

func (a *argon) GenerateHashFromPassword(password string) (string, error) {
	if !entity.IsPasswordValid(password) {
		return "", apperror.ErrDataNotValid
	}

	hashPasswordByte := argon2.IDKey([]byte(password), a.salt, a.time, a.memory, a.threads, a.keyLen)
	hashPasswordString := hex.EncodeToString(hashPasswordByte)

	return hashPasswordString, nil
}

func (a *argon) VerifyPassword(hashPassword string, password string) error {
	newHashByte := argon2.IDKey([]byte(password), a.salt, a.time, a.memory, a.threads, a.keyLen)

	newHashString := hex.EncodeToString(newHashByte)

	if subtle.ConstantTimeCompare([]byte(hashPassword), []byte(newHashString)) != 1 {
		return apperror.ErrHashPasswordsNotEqual
	}

	return nil
}
