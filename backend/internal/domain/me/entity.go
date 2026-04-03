package me

import (
	"errors"
	"time"
)

type Me struct {
	identityID     identityID
	displayName    displayName
	displayNameJa  *displayNameJa
	role           *role
	location       *location
	skills         []skillCategory
	certifications []Certification
	experiences    []experience
	links          []Link
	likes          []like
	createdAt      time.Time
	updatedAt      time.Time
}

type OptFunc func(*Me) error

// --- Factory 関数 ---

// NewMe はMeエンティティを作成する
func NewMe(id string, name string, opts ...OptFunc) (*Me, error) {
	identityID, err := newIdentityID(id)
	if err != nil {
		return nil, err
	}
	dn, err := newDisplayName(name)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	e := &Me{
		identityID:  identityID,
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

// ReconstructInput はReconstructの入力型
type ReconstructInput struct {
	ID             string
	Name           string
	DisplayJa      *string
	Role           *string
	Location       *string
	Likes          []string
	Links          []Link
	Certifications []Certification
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Reconstruct はDBから取得した信頼済みデータでエンティティを復元する
func Reconstruct(input ReconstructInput) *Me {
	e := &Me{
		identityID:  identityID{value: input.ID},
		displayName: displayName{value: input.Name},
		createdAt:   input.CreatedAt,
		updatedAt:   input.UpdatedAt,
	}
	if input.DisplayJa != nil {
		v := displayNameJa{value: *input.DisplayJa}
		e.displayNameJa = &v
	}
	if input.Role != nil {
		v := role{value: *input.Role}
		e.role = &v
	}
	if input.Location != nil {
		v := location{value: *input.Location}
		e.location = &v
	}
	for _, s := range input.Likes {
		e.likes = append(e.likes, like{value: s})
	}
	e.links = input.Links
	e.certifications = input.Certifications
	return e
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
	next.certifications = []Certification{}
	next.experiences = []experience{}
	next.links = []Link{}
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

// OptLinks はlinksを設定するオプション
func OptLinks(input []Link) OptFunc {
	return func(m *Me) error {
		m.links = input
		return nil
	}
}

// OptLikes はlinksを設定するオプション
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

// OptCertifications はcertificationsを設定するオプション
func OptCertifications(
	input []Certification,
) OptFunc {
	return func(m *Me) error {
		m.certifications = input
		return nil
	}
}

// --- Getter ---

// ID はIDの値を返す
func (e *Me) ID() string {
	return e.identityID.value
}

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

// Links はlinksの値を返す
func (e *Me) Links() []Link {
	return e.links
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

// Certifications はcertificationsの値を返す
func (e *Me) Certifications() []Certification {
	return e.certifications
}

// CreatedAt はcreatedAtフィールドのgetter
func (e *Me) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt はupdatedAtフィールドのgetter
func (e *Me) UpdatedAt() time.Time {
	return e.updatedAt
}
