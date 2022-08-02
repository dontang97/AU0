package dto

import (
	"strings"
	"testing"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClaimsSuite struct {
	suite.Suite
}

func (s *ClaimsSuite) SetupSuite() {
}

func (s *ClaimsSuite) TearDownSuite() {
}

func (s *ClaimsSuite) SetupTest() {
}

func (s *ClaimsSuite) TearDownTest() {
}

func (s *ClaimsSuite) TestToCsv() {
	claims := &Claims{
		IssuedAt:        1659163848000,
		AccountID:       uuid.MustParse("970a489a-38da-484a-b42a-eecd07874dee"),
		VirtualDeviceID: "",
		HandsetDeviceID: "",
		Verified:        true,
		HandSetVerified: false,
		LegalDocToSign:  false,
		CountryCode:     countries.TW,
		DataCenterID:    uuid.MustParse("45407e7c-ab11-4c9c-8a49-569d066ed97c"),
		IDKey:           "",
		IssueTo:         uuid.MustParse("054f6b80-7c2a-4cad-b845-30fd0c5b9b0f"),
		Scopes:          []string{"issuetoken", "email", "birthday", "profile.write", "payment.security.write"},
		ExpiredInterval: 86400,
	}

	csvStr, err := claims.toCsv(",")
	s.Equal(nil, err)

	s.Equal(14, len(strings.Split(csvStr, ",")))
	s.Equal("1659163848000,970a489a-38da-484a-b42a-eecd07874dee,,,true,false,0.0.0.0,false,TW,45407e7c-ab11-4c9c-8a49-569d066ed97c,,054f6b80-7c2a-4cad-b845-30fd0c5b9b0f,AARH,86400", csvStr)
}

func (s *ClaimsSuite) TestFromCsv() {
	claims := &Claims{}

	err := claims.fromCsv("1659163848000,970a489a-38da-484a-b42a-eecd07874dee,,,true,false,0.0.0.0,false,TW,45407e7c-ab11-4c9c-8a49-569d066ed97c,,054f6b80-7c2a-4cad-b845-30fd0c5b9b0f,AARH,86400", ",")
	s.Equal(nil, err)

	s.Equal(&Claims{
		IssuedAt:        1659163848000,
		AccountID:       uuid.MustParse("970a489a-38da-484a-b42a-eecd07874dee"),
		VirtualDeviceID: "",
		HandsetDeviceID: "",
		Verified:        true,
		HandSetVerified: false,
		LegalDocToSign:  false,
		CountryCode:     countries.TW,
		DataCenterID:    uuid.MustParse("45407e7c-ab11-4c9c-8a49-569d066ed97c"),
		IDKey:           "",
		IssueTo:         uuid.MustParse("054f6b80-7c2a-4cad-b845-30fd0c5b9b0f"),
		Scopes:          []string{"issuetoken", "email", "birthday", "profile.write", "payment.security.write"},
		ExpiredInterval: 86400,
	}, claims)
}

func TestRunClaimsSuite(t *testing.T) {
	suite.Run(t, new(ClaimsSuite))
}
