package cryptor_test

import (
	"testing"

	"github.com/dontang97/AU0/dto"
	"github.com/dontang97/AU0/util/cryptor"
	"github.com/stretchr/testify/suite"
)

type HTCTokenCryptorSuite struct {
	suite.Suite
	cryptor dto.TokenCryptor
}

func (s *HTCTokenCryptorSuite) SetupSuite() {
	s.cryptor, _ = cryptor.NewHTCTokenCryptor(
		[]byte("LIjhBWwR1A9BjiBCdsMs0KiZ3x50Ce9auGFBqqj69Q4="),
		[]byte("tKIlYqe00NtAuXhDy1UfYQ=="),
	)
}

func (s *HTCTokenCryptorSuite) TearDownSuite() {
}

func (s *HTCTokenCryptorSuite) SetupTest() {
}

func (s *HTCTokenCryptorSuite) TearDownTest() {
}

func (s *HTCTokenCryptorSuite) TestEncrypt() {
	raw := "1659163848000,9e7475f8-97ce-485f-81de-e57d59333278,,,true,false,0.0.0.0,false,,5ecf0bd8-9e9c-4def-85f0-b59fc6458b00,,7b398d3b-6d8d-4030-9057-63c476fd2d2e,AARH,86400"
	enc, err := s.cryptor.Encrypt([]byte(raw))

	s.Equal(nil, err)
	s.Equal("C3G+LcrldXHM1umtdflv/5WABOVfsohCbQpmkD8X0fVNFsg/AsvBWPUG1unirWF0lXvGn40ORnzettjsCUav+fS/nXYxHve1si8Z1ZuE2fKgZhI6gLJIWVOUFvuvbTucSyqO4MWWfsUPk78R4wZzdJ3gugWa/n7mlfwbatyoEoAoltY0ndQp0vzv7YjGRgUl5Gv0mpxIiAL7AL80YKWt5JjTRMZ0CtQnXFykuWZSlz53VTOSCE1rXwb4sarGlbj6ZI6Z6MbaPU8EjHC6q5zXmqqqG9xotDIbOzn9fJaOu5NolXL9f/Zjlu3mbF76+ktei4oJUkfksO6/IqoBAhd9PwTfyj1f5Ws4aWa/4oYRS3Dri8kBKPWzpeQf4G9s92HITzmUabhZsFeSnnkuH2b28MuX21AOU+fTxgoowj8Wk1YPuWgtmt+0oNgCuxp23Ux8", string(enc))
}

func (s *HTCTokenCryptorSuite) TestDecrypt() {
	dec, err := s.cryptor.Decrypt([]byte("C3G+LcrldXHM1umtdflv/5WABOVfsohCbQpmkD8X0fVNFsg/AsvBWPUG1unirWF0lXvGn40ORnzettjsCUav+fS/nXYxHve1si8Z1ZuE2fKgZhI6gLJIWVOUFvuvbTucSyqO4MWWfsUPk78R4wZzdJ3gugWa/n7mlfwbatyoEoAoltY0ndQp0vzv7YjGRgUl5Gv0mpxIiAL7AL80YKWt5JjTRMZ0CtQnXFykuWZSlz53VTOSCE1rXwb4sarGlbj6ZI6Z6MbaPU8EjHC6q5zXmqqqG9xotDIbOzn9fJaOu5NolXL9f/Zjlu3mbF76+ktei4oJUkfksO6/IqoBAhd9PwTfyj1f5Ws4aWa/4oYRS3Dri8kBKPWzpeQf4G9s92HITzmUabhZsFeSnnkuH2b28MuX21AOU+fTxgoowj8Wk1YPuWgtmt+0oNgCuxp23Ux8"))

	s.Equal(nil, err)
	s.Equal("1659163848000,9e7475f8-97ce-485f-81de-e57d59333278,,,true,false,0.0.0.0,false,TW,5ecf0bd8-9e9c-4def-85f0-b59fc6458b00,,7b398d3b-6d8d-4030-9057-63c476fd2d2e,AARH,86400", string(dec))
}

func TestRunHTCTokenCryptorSuite(t *testing.T) {
	suite.Run(t, new(HTCTokenCryptorSuite))
}
