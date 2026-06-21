package crypto

import (
	"fmt"

	"golang.org/x/crypto/argon2"
)

type KDFParams struct {
	MemoryKib   uint32 `json:"memory_kib"`
	Iterations  uint32 `json:"iterations"`
	Parallelism uint8  `json:"parallelism"`
	Salt        []byte `json:"salt"`
}

type Profile string

const (
	ProfileFast     Profile = "fast"
	ProfileDefault  Profile = "default"
	ProfileParanoid Profile = "paranoid"
)

var ProfileParams = map[Profile]KDFParams{
	ProfileFast: {
		MemoryKib:   32768,
		Iterations:  2,
		Parallelism: 2,
	},
	ProfileDefault: {
		MemoryKib:   65536,
		Iterations:  3,
		Parallelism: 4,
	},
	ProfileParanoid: {
		MemoryKib:   262144,
		Iterations:  4,
		Parallelism: 4,
	},
}

func GetProfileParams(profile Profile) (KDFParams, error) {
	params, ok := ProfileParams[profile]
	if !ok {
		return KDFParams{}, fmt.Errorf("unknown profile: %s", profile)
	}
	return params, nil
}

func DeriveKey(password string, params KDFParams) ([]byte, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	if len(params.Salt) == 0 {
		return nil, fmt.Errorf("salt cannot be empty")
	}

	key := argon2.IDKey(
		[]byte(password),
		params.Salt,
		params.Iterations,
		params.MemoryKib,
		params.Parallelism,
		32,
	)

	return key, nil
}

func (p KDFParams) WithSalt(salt []byte) KDFParams {
	return KDFParams{
		MemoryKib:   p.MemoryKib,
		Iterations:  p.Iterations,
		Parallelism: p.Parallelism,
		Salt:        salt,
	}
}
