package dto

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/biter777/countries"
	"github.com/google/uuid"
)

const (
	claimsFiledNum = 13
)

type Claims struct {
	IssuedAt        int64 //unix milli seconds timestamp
	AccountID       uuid.UUID
	VirtualDeviceID string
	HandsetDeviceID string
	Verified        bool
	HandSetVerified bool
	LegalDocToSign  bool
	CountryCode     countries.CountryCode
	DataCenterID    uuid.UUID
	IDKey           string
	IssueTo         uuid.UUID
	Scopes          []string
	ExpiredInterval int
}

func (claims *Claims) toCsv(sep string) (string, error) {
	var strs []string

	strs = append(strs, strconv.FormatInt(claims.IssuedAt, 10))
	strs = append(strs, claims.AccountID.String())
	strs = append(strs, claims.VirtualDeviceID)
	strs = append(strs, claims.HandsetDeviceID)
	strs = append(strs, strconv.FormatBool(claims.Verified))
	strs = append(strs, strconv.FormatBool(claims.HandSetVerified))

	// constant
	strs = append(strs, "0.0.0.0")

	strs = append(strs, strconv.FormatBool(claims.LegalDocToSign))
	strs = append(strs, claims.CountryCode.Alpha2())
	strs = append(strs, claims.DataCenterID.String())
	strs = append(strs, claims.IDKey)
	strs = append(strs, claims.IssueTo.String())

	scopesBitmap, err := ScopeStrs2bitmap(claims.Scopes)
	if err != nil {
		return "", err
	}

	strs = append(strs, b64.StdEncoding.EncodeToString(scopesBitmap))

	strs = append(strs, strconv.Itoa(claims.ExpiredInterval))

	return strings.Join(strs, sep), nil
}

func (claims *Claims) fromCsv(csvStr string, sep string) error {
	strs := strings.Split(csvStr, sep)

	if strLen := len(strs); strLen != (claimsFiledNum + 1) {
		return fmt.Errorf("invalid claim fields num (expected: %d, got: %d)", (claimsFiledNum + 1), strLen)
	}

	issuedAt, err := strconv.ParseInt(strs[0], 10, 64)
	if err != nil {
		return err
	}
	claims.IssuedAt = issuedAt

	accountID, err := uuid.Parse(strs[1])
	if err != nil {
		return err
	}
	claims.AccountID = accountID

	claims.VirtualDeviceID = strs[2]
	claims.HandsetDeviceID = strs[3]

	verified, err := strconv.ParseBool(strs[4])
	if err != nil {
		return err
	}
	claims.Verified = verified

	handSetVerified, err := strconv.ParseBool(strs[5])
	if err != nil {
		return err
	}
	claims.HandSetVerified = handSetVerified

	/*
		if strs[6] != "0.0.0.0" {
			return errors.New("7th column must be 0.0.0.0")
		}
	*/

	legalDocToSign, err := strconv.ParseBool(strs[7])
	if err != nil {
		return err
	}
	claims.LegalDocToSign = legalDocToSign

	claims.CountryCode = countries.ByName(strs[8])

	dataCenterID, err := uuid.Parse(strs[9])
	if err != nil {
		return err
	}
	claims.DataCenterID = dataCenterID

	claims.IDKey = strs[10]

	issueTo, err := uuid.Parse(strs[11])
	if err != nil {
		return err
	}
	claims.IssueTo = issueTo

	scopesBitmap, err := b64.StdEncoding.DecodeString(strs[12])
	if err != nil {
		return err
	}

	claims.Scopes = Bitmap2ScopeStrs(scopesBitmap)

	expiredInterval, err := strconv.Atoi(strs[13])
	if err != nil {
		return err
	}
	claims.ExpiredInterval = expiredInterval

	return nil
}

func (claims *Claims) Encrypt(cryptor TokenCryptor) ([]byte, error) {
	csvStr, err := claims.toCsv(",")
	if err != nil {
		return nil, err
	}
	return cryptor.Encrypt([]byte(csvStr))
}

func (claims *Claims) Decrypt(cryptor TokenCryptor, raw []byte) error {
	d, err := cryptor.Decrypt(raw)
	if err != nil {
		return err
	}

	return claims.fromCsv(string(d), ",")
}
