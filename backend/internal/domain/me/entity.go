package me

import (
	"errors"
	"time"
)

type Me struct {
	displayName    displayName
	displayNameJa  *displayNameJa
	role           *role
	location       *location
	skills         []skillCategory
	certifications []certification
	experiences    []experience
	links          []link
	likes          []like
	createdAt      time.Time
	updatedAt      time.Time
}

type OptFunc func(*Me) error

// --- Factory 関数 ---

// NewMe はMeエンティティを作成する
func NewMe(name string, opts ...OptFunc) (*Me, error) {
	dn, err := newDisplayName(name)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	e := &Me{
		displayName: dn,
		createdAt:   now,
		updatedAt:   now,
	}

	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("nil option is not allowed")
		}
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// Reconstruct はMeエンティティを再構築する
func Reconstruct() {
}

// --- 振る舞い ---

// Update は更新関数
func (e *Me) Update(name string, opts ...OptFunc) error {
	dn, err := newDisplayName(name)
	if err != nil {
		return err
	}

	next := *e
	next.displayName = dn
	next.displayNameJa = nil
	next.role = nil
	next.location = nil
	next.skills = []skillCategory{}
	next.certifications = []certification{}
	next.experiences = []experience{}
	next.links = []link{}
	next.likes = []like{}
	for _, opt := range opts {
		if opt == nil {
			return errors.New("nil option is not allowed")
		}
		if err := opt(&next); err != nil {
			return err
		}
	}
	next.updatedAt = time.Now()
	*e = next

	return nil
}

// --- Functional Option 関数 ---

// OptDisplayNameJa はdisplayNameJaを設定するオプション
func OptDisplayNameJa(input string) OptFunc {
	return func(m *Me) error {
		value, err := newDisplayNameJa(input)
		if err != nil {
			return err
		}
		m.displayNameJa = &value
		return nil
	}
}

// OptRole はroleを設定するオプション
func OptRole(input string) OptFunc {
	return func(m *Me) error {
		value, err := newRole(input)
		if err != nil {
			return err
		}
		m.role = &value
		return nil
	}
}

// OptLocation はlocationを設定するオプション
func OptLocation(input string) OptFunc {
	return func(m *Me) error {
		value, err := newLocation(input)
		if err != nil {
			return err
		}
		m.location = &value
		return nil
	}
}

// OptLikes はLikesを設定するオプション
func OptLikes(input []string) OptFunc {
	return func(m *Me) error {
		likes := []like{}
		for _, s := range input {
			value, err := newLike(s)
			if err != nil {
				return err
			}
			likes = append(likes, value)
		}
		m.likes = likes
		return nil
	}
}

// --- Getter ---

// DisplayName はdisplayNameの値を返す
func (e *Me) DisplayName() string {
	return e.displayName.Value()
}

// DisplayNameJa はdisplayNameJaの値を返す。未設定の場合は空文字を返す。
func (e *Me) DisplayNameJa() string {
	if e.displayNameJa == nil {
		return ""
	}
	return e.displayNameJa.Value()
}

// Role はroleの値を返す。未設定の場合は空文字を返す。
func (e *Me) Role() string {
	if e.role == nil {
		return ""
	}
	return e.role.Value()
}

// Location はlocationの値を返す。未設定の場合は空文字を返す。
func (e *Me) Location() string {
	if e.location == nil {
		return ""
	}
	return e.location.Value()
}

// Likes はlikesの値を返す
func (e *Me) Likes() []string {
	if len(e.likes) == 0 {
		return []string{}
	}
	val := make([]string, 0, len(e.likes))
	for _, o := range e.likes {
		val = append(val, o.Value())
	}
	return val
}

// CreatedAt はcreatedAtフィールドのgetter
func (e *Me) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt はupdatedAtフィールドのgetter
func (e *Me) UpdatedAt() time.Time {
	return e.updatedAt
}
