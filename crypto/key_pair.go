package crypto

type KeyPair interface {
	GetPrivateKeyPEM() []byte
	GetPublicKeyPEM() []byte
}
