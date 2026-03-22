package me

import (
	"errors"
	"net/url"
	"strings"
)

type (
	// Me の vo
	displayName   struct{ value string }
	displayNameJa struct{ value string }
	role          struct{ value string }
	location      struct{ value string }
	skillCategory struct {
		category  struct{ value string }
		items     []string
		sortOrder struct{ value int }
	}
	certification struct {
		name   struct{ value string }
		issuer *struct{ value string }
		year   struct{ value int }
		month  struct{ value int }
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

// Platform はplatformの値を返す
func (l Link) Platform() string {
	return l.platform
}

// URL はurlの値を返す
func (l Link) URL() string {
	return l.url
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
