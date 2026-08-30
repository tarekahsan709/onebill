// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package crypto // import "onebill.app/internal/crypto"

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"hash/fnv"

	"golang.org/x/crypto/bcrypt"
)

// HashFromBytes returns a non-cryptographic checksum of the input.
func HashFromBytes(value []byte) string {
	h := fnv.New128a()
	h.Write(value)
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateRandomBytes returns random bytes.
func GenerateRandomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

// GenerateRandomStringHex returns a random hexadecimal string.
func GenerateRandomStringHex(size int) string {
	return hex.EncodeToString(GenerateRandomBytes(size))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ConstantTimeCmp(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
