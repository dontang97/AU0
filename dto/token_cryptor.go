package dto

type TokenCryptor interface {
	Encrypt(raw []byte) ([]byte, error)
	Decrypt(crypt []byte) ([]byte, error)
}
