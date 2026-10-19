package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	memory  uint32
	time    uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewHashConfig(memory, time, keylen, saltlen uint32, threads uint8) *HashConfig {
	return &HashConfig{
		memory:  memory,
		time:    time,
		threads: threads,
		keyLen:  keylen,
		saltLen: saltlen,
	}
}

func NewRecommendedHashConfig() *HashConfig {
	return &HashConfig{
		memory:  64 * 1024,
		time:    2,
		threads: 2,
		keyLen:  32,
		saltLen: 16,
	}
}

func (h *HashConfig) GenHash(password string) string {
	salt := h.genSalt()
	hash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

	base64Hash := base64.RawStdEncoding.EncodeToString(hash)
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)

	completeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.memory, h.time, h.threads, base64Salt, base64Hash)

	return completeHash
}

func (h *HashConfig) genSalt() []byte {
	salt := make([]byte, h.saltLen)
	rand.Read(salt)

	return salt
}

func Compare(password string, hashed string) error {
	result := strings.Split(hashed, "$")
	if len(result) != 6 {
		return errors.New("Incorrect Hash Format.")
	}

	if result[1] != "argon2id" {
		return errors.New("Incorrect Hash Type.")
	}

	var version int
	if _, err := fmt.Sscanf(result[2], "v=%d", &version); err != nil {
		return err
	}
	if version != argon2.Version {
		return errors.New("Incorrect Hash Version.")
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(result[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(result[4])
	if err != nil {
		return err
	}

	hash, err := base64.RawStdEncoding.DecodeString(result[5])
	if err != nil {
		return err
	}

	newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return errors.New("Hash Mismatch.")
	}

	return nil
}
