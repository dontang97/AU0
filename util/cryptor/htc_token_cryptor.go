package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"

	"github.com/dontang97/AU0/dto"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

// refer to https://gist.github.com/hothero/7d085573f5cb7cdb5801d7adcf66dcf3
type HTCTokenCryptor struct {
	block  cipher.Block
	ciphe  cipher.BlockMode
	ciphd  cipher.BlockMode
	utf16e *encoding.Encoder
	utf16d *encoding.Decoder
}

func NewHTCTokenCryptor(key []byte, initVector []byte) (dto.TokenCryptor, error) {
	rawkey, err := b64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return nil, err
	}

	rawIV, err := b64.StdEncoding.DecodeString(string(initVector))
	if err != nil {
		return nil, err
	}

	b, err := aes.NewCipher(rawkey)
	if err != nil {
		return nil, err
	}

	return &HTCTokenCryptor{
		block:  b,
		ciphe:  cipher.NewCBCEncrypter(b, rawIV),
		ciphd:  cipher.NewCBCDecrypter(b, rawIV),
		utf16e: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(),
		utf16d: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(),
	}, nil
}

func (cryptor *HTCTokenCryptor) Encrypt(raw []byte) ([]byte, error) {
	utf16b, err := cryptor.utf16e.Bytes(raw)
	if err != nil {
		return nil, err
	}

	blockSize := cryptor.block.BlockSize()

	content := make([]byte, len(utf16b))
	copy(content, utf16b)
	padding := blockSize - len(content)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	padtext = append(content, padtext...)

	crypted := make([]byte, len(padtext))
	cryptor.ciphe.CryptBlocks(crypted, padtext)

	return []byte(b64.StdEncoding.EncodeToString(crypted)), nil
}

func (cryptor *HTCTokenCryptor) Decrypt(crypt []byte) ([]byte, error) {
	crypt, err := b64.StdEncoding.DecodeString(string(crypt))
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(crypt))
	cryptor.ciphd.CryptBlocks(decrypted, crypt)

	padding := decrypted[len(decrypted)-1]
	return cryptor.utf16d.Bytes(decrypted[:len(decrypted)-int(padding)])
}
