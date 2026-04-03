package me

import (
	"errors"
	"net/url"
	"strings"
)

type (
	// Me の vo
	identityID    struct{ value string }
	displayName   struct{ value string }
	displayNameJa struct{ value string }
	role          struct{ value string }
	location      struct{ value string }
	skillCategory struct {
		category  struct{ value string }
		items     []string
		sortOrder struct{ value int }
	}
	Certification struct {
		name   string
		issuer string
		year   int
		month  int
	}
	experience struct {
		company   struct{ value string }
		url       struct{ value string }
		startYear struct{ value int }
		endYear   *struct{ value int }
	}
	Link struct {
		platform string
		url      string
	}
	like struct{ value string }
)

// newIdentityID は identityID オブジェクトを生成
func newIdentityID(
	input string,
) (identityID, error) {
	return identityID{
		value: input,
	}, nil
}

// newDisplayName はdisplayNameオブジェクトを生成
func newDisplayName(
	input string,
) (displayName, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return displayName{}, err
	}
	return displayName{
		value: input,
	}, nil
}

// newDisplayNameJa はnewDisplayNameJaオブジェクトを生成
func newDisplayNameJa(
	input string,
) (displayNameJa, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return displayNameJa{}, err
	}
	return displayNameJa{
		value: input,
	}, nil
}

// newRole はroleオブジェクトを生成
func newRole(
	input string,
) (role, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return role{}, err
	}
	return role{
		value: input,
	}, nil
}

// newLocation はlocationオブジェクトを生成
func newLocation(
	input string,
) (location, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return location{}, err
	}
	return location{
		value: input,
	}, nil
}

// NewCertification はCertificationオブジェクトを生成
func NewCertification(
	inputName, inputIssuer string,
	inputYear, inputMonth int,
) (Certification, error) {
	err := validateNonEmpty(inputName)
	if err != nil {
		return Certification{}, err
	}
	err = validatePositiveInt(inputYear)
	if err != nil {
		return Certification{}, err
	}
	err = validateMonth(inputMonth)
	if err != nil {
		return Certification{}, err
	}

	return Certification{
		name:   inputName,
		issuer: inputIssuer,
		year:   inputYear,
		month:  inputMonth,
	}, nil
}

// newLike はlikeオブジェクトを生成
func newLike(input string) (like, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return like{}, err
	}
	return like{
		value: input,
	}, nil
}

// NewLink はLinkオブジェクトを生成
func NewLink(inputPlatform, inputURL string) (Link, error) {
	err := validateNonEmpty(inputPlatform)
	if err != nil {
		return Link{}, err
	}
	err = validateNonEmpty(inputURL)
	if err != nil {
		return Link{}, err
	}
	_, err = url.ParseRequestURI(inputURL)
	if err != nil {
		return Link{}, err
	}

	return Link{
		platform: inputPlatform,
		url:      inputURL,
	}, nil
}

// Getter

// Value はgetterメソッド
func (vo displayName) Value() string {
	return vo.value
}

// Value はgetterメソッド
func (vo displayNameJa) Value() string {
	return vo.value
}

// Value はgetterメソッド
func (vo role) Value() string {
	return vo.value
}

// Value はgetterメソッド
func (vo location) Value() string {
	return vo.value
}

// Value はgetterメソッド
func (vo like) Value() string {
	return vo.value
}

// Platform はplatformの値を返す
func (l Link) Platform() string {
	return l.platform
}

// URL はurlの値を返す
func (l Link) URL() string {
	return l.url
}

// Name はnameの値を返す
func (c Certification) Name() string {
	return c.name
}

// Issuer はissuerの値を返す
func (c Certification) Issuer() string {
	return c.issuer
}

// Year はyearの値を返す
func (c Certification) Year() int {
	return c.year
}

// Month はmonthの値を返す
func (c Certification) Month() int {
	return c.month
}

// 共通

// バリデーション関数

// validateNonEmpty
func validateNonEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("must not be empty")
	}
	return nil
}

func validatePositiveInt(value int) error {
	if value < 1 {
		return errors.New("must be positive")
	}
	return nil
}

func validateMonth(value int) error {
	if value < 1 || value > 12 {
		return errors.New("must be between 1 and 12")
	}
	return nil
}
