package me

import (
	"errors"
	"slices"
	"time"
)

type (
	Me struct {
		profile        profile
		skills         Skills
		certifications []Certification
		experiences    []experience
		links          []Link
		likes          []like
		updatedAt      time.Time
	}
	ReconstructInput struct {
		DisplayName    string
		DisplayNameJa  string
		Role           string
		Location       string
		Likes          []string
		Links          []Link
		Skills         Skills
		Certifications []Certification
		UpdatedAt      time.Time
	}
	OptFunc        func(*Me) error
	OptProfileFunc func(*profile) error
)

func Default() *Me {
	return &Me{
		updatedAt: time.Now(),
	}
}

// Reconstruct はDBから取得した信頼済みデータでエンティティを復元する
func Reconstruct(input ReconstructInput) *Me {
	e := &Me{
		updatedAt: input.UpdatedAt,
	}
	e.profile.displayName = input.DisplayName
	e.profile.displayNameJa = input.DisplayNameJa
	e.profile.role = input.Role
	e.profile.location = input.Location
	for _, s := range input.Likes {
		e.likes = append(e.likes, like{value: s})
	}
	e.links = input.Links
	e.certifications = input.Certifications
	e.skills = input.Skills
	return e
}

func (e *Me) UpdateProfile(baseTime time.Time, in ...OptProfileFunc) error {
	if in == nil {
		return errors.New("入力内容が不正です")
	}
	if len(in) == 0 {
		return errors.New("入力内容が不正です")
	}
	current := e.profile
	for _, opt := range in {
		err := opt(&current)
		if err != nil {
			return err
		}
	}
	e.updatedAt = baseTime
	e.profile = current
	return nil
}

func (e *Me) UpdateLikes(in []string, baseTime time.Time) error {
	optsFn := OptLikes(in)
	err := optsFn(e)
	if err != nil {
		return err
	}
	e.updatedAt = baseTime
	return nil
}

func (e *Me) UpdateLinks(in []Link, baseTime time.Time) error {
	optsFn := OptLinks(in)
	err := optsFn(e)
	if err != nil {
		return err
	}
	e.updatedAt = baseTime
	return nil
}

func (e *Me) AddSkill(
	itemName, categoryName string,
	baseTime time.Time,
) error {
	if e.skills == nil {
		e.skills = make(Skills)
	}
	if slices.Contains(e.skills[categoryName].Items, itemName) {
		return errors.New("既に登録済み")
	}
	current := e.skills[categoryName]
	current.Items = append(current.Items, itemName)
	e.skills[categoryName] = current
	e.updatedAt = baseTime
	return nil
}

func (e *Me) RemoveSkill(itemName, categoryName string, baseTime time.Time) error {
	if e.skills == nil {
		e.skills = make(Skills)
	}
	if !slices.Contains(e.skills[categoryName].Items, itemName) {
		return errors.New("登録されていません")
	}
	current := e.skills[categoryName]
	current.Items = slices.DeleteFunc(
		current.Items,
		func(s string) bool { return s == itemName },
	)
	if len(current.Items) == 0 {
		delete(e.skills, categoryName)
	} else {
		e.skills[categoryName] = current
	}
	e.updatedAt = baseTime
	return nil
}

// Update は更新関数
func (e *Me) Update(name string, opts ...OptFunc) error {
	next := *e
	// next.skills = []skillCategory{}
	next.certifications = []Certification{}
	next.experiences = []experience{}
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

func OptDisplayName(in string) OptProfileFunc {
	return func(p *profile) error {
		p.displayName = in
		return nil
	}
}

// OptDisplayNameJa はdisplayNameJaを設定するオプション
func OptDisplayNameJa(in string) OptProfileFunc {
	return func(p *profile) error {
		p.displayNameJa = in
		return nil
	}
}

// OptRole はroleを設定するオプション
func OptRole(in string) OptProfileFunc {
	return func(p *profile) error {
		p.role = in
		return nil
	}
}

// OptLocation はlocationを設定するオプション
func OptLocation(in string) OptProfileFunc {
	return func(p *profile) error {
		p.location = in
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

// DisplayName はdisplayNameの値を返す
func (e *Me) DisplayName() string {
	return e.profile.displayName
}

// DisplayNameJa はdisplayNameJaの値を返す。未設定の場合は空文字を返す。
func (e *Me) DisplayNameJa() string {
	return e.profile.displayNameJa
}

// Role はroleの値を返す。未設定の場合は空文字を返す。
func (e *Me) Role() string {
	return e.profile.role
}

// Location はlocationの値を返す。未設定の場合は空文字を返す。
func (e *Me) Location() string {
	return e.profile.location
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

func (e *Me) Skills() Skills {
	return e.skills
}

// Certifications はcertificationsの値を返す
func (e *Me) Certifications() []Certification {
	return e.certifications
}

// UpdatedAt はupdatedAtフィールドのgetter
func (e *Me) UpdatedAt() time.Time {
	return e.updatedAt
}
