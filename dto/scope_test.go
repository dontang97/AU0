package dto_test

import (
	"testing"

	"github.com/dontang97/AU0/dto"
	"github.com/stretchr/testify/suite"
)

type ScopeSuite struct {
	suite.Suite
}

func (s *ScopeSuite) SetupSuite() {
}

func (s *ScopeSuite) TearDownSuite() {
}

func (s *ScopeSuite) SetupTest() {
}

func (s *ScopeSuite) TearDownTest() {
}

func (s *ScopeSuite) TestScopeStrings2Bitmap() {
	//0 1 2 6 10
	scopes := []string{
		"issuetoken",
		"email",
		"birthday",
		"profile.write",
		"payment.security.write",
	}

	bitmap, err := dto.ScopeStrs2bitmap(scopes)

	s.Equal(nil, err)
	s.Equal(3, len(bitmap))
	s.Equal([]byte{0, 4, 71}, bitmap) // 00000000 00000100 010000111

	scopes = []string{"undefined"}
	bitmap, err = dto.ScopeStrs2bitmap(scopes)

	s.EqualError(err, "undefined scope: undefined")
	s.Equal(([]byte(nil)), bitmap)
}

func (s *ScopeSuite) TestBitmap2ScopeStrings() {
	//0 1 2 6 10
	bitmap := []byte{0, 4, 71}

	scopeStrs := dto.Bitmap2ScopeStrs(bitmap)

	s.Equal(5, len(scopeStrs))

	s.Equal(scopeStrs, []string{
		"issuetoken",
		"email",
		"birthday",
		"profile.write",
		"payment.security.write",
	})
}

func TestRunScopeSuite(t *testing.T) {
	suite.Run(t, new(ScopeSuite))
}
