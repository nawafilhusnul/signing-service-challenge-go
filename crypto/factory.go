package crypto

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func NewSignerFromDevice(algorithm string, privateKeyPEM []byte) (Signer, error) {
	switch algorithm {
	case domain.AlgorithmRSA:
		marshaler := NewRSAMarshaler()
		keyPair, err := marshaler.Unmarshal(privateKeyPEM)
		if err != nil {
			return nil, err
		}
		return NewRSASigner(keyPair.Private), nil

	case domain.AlgorithmECC:
		marshaler := NewECCMarshaler()
		keyPair, err := marshaler.Decode(privateKeyPEM)
		if err != nil {
			return nil, err
		}
		return NewECDSASigner(keyPair.Private), nil

	default:
		return nil, domain.ErrInvalidAlgorithm
	}
}

type Generator interface {
	Generate() (KeyPair, error)
}

func NewGenerator(algorithm string) (Generator, error) {
	switch algorithm {
	case domain.AlgorithmRSA:
		return &RSAGenerator{}, nil
	case domain.AlgorithmECC:
		return &ECCGenerator{}, nil
	default:
		return nil, domain.ErrInvalidAlgorithm
	}
}
