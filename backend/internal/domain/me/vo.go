package me

import (
	"errors"
	"net/url"
	"strings"
)

type (
	// Me の vo
	profile struct {
		displayName   string
		displayNameJa string
		role          string
		location      string
	}
	skillCategory struct {
		// category  struct{ value string }
		// items     []string
		// sortOrder struct{ value int }
	}
	Certification struct {
		name   string
		issuer string
		year   int
		month  int
	}
	experience struct {
		// company   struct{ value string }
		// url       struct{ value string }
		// startYear struct{ value int }
		// endYear   *struct{ value int }
	}
	Link struct {
		platform string
		url      string
	}
	like struct{ value string }
)

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
