package me

import "time"

// InputDto DTO定義
type InputDto struct {
	Certifications []struct {
		Issuer string `json:"issuer,omitempty"`
		Month  int    `json:"month,required"`
		Name   string `json:"name"                validate:"required"`
		Year   int    `json:"year"                validate:"required"`
	} `json:"certifications" validate:"dive"`
	Experiences []struct {
		Company   string  `json:"company"             validate:"required"`
		EndYear   *int    `json:"endYear,omitempty"`
		StartYear int     `json:"startYear"           validate:"required"`
		URL       *string `json:"url,omitempty"       validate:"omitempty,url"`
	} `json:"experiences"`
	Likes []string `json:"likes" validate:"omitempty"`
	Links []struct {
		Label    *string `json:"label,omitempty"`
		Platform string  `json:"platform"            validate:"required"`
		URL      string  `json:"url"                 validate:"required,url"`
	} `json:"links" validate:"dive"`
	Location    *string `json:"location"`
	DisplayName string  `json:"displayName"         validate:"required"`
	DisplayJa   *string `json:"displayJa,omitempty"`
	Role        *string `json:"role"`
	Skills      []struct {
		Category  string   `json:"category"            validate:"required"`
		Items     []string `json:"items"               validate:"required,min=1"`
		SortOrder int      `json:"sortOrder"`
	} `json:"skills" validate:"dive"`
}

// OutputDto DTO定義
type OutputDto struct {
	Certifications []struct {
		Issuer string `json:"issuer,omitempty"`
		Month  int    `json:"month"`
		Name   string `json:"name"`
		Year   int    `json:"year"`
	} `json:"certifications,omitempty"`
	Experiences []struct {
		Company   string  `json:"company"`
		EndYear   *int    `json:"endYear,omitempty"`
		StartYear int     `json:"startYear"`
		URL       *string `json:"url,omitempty"`
	} `json:"experiences,omitempty"`
	Likes []string `json:"likes,omitempty"`
	Links []struct {
		Platform string `json:"platform"`
		URL      string `json:"url"`
	} `json:"links,omitempty"`
	Location    string `json:"location,omitempty"`
	DisplayName string `json:"displayName"`
	DisplayJa   string `json:"displayJa,omitempty"`
	Role        string `json:"role,omitempty"`
	Skills      []struct {
		Category  string   `json:"category"`
		Items     []string `json:"items"`
		SortOrder int      `json:"sortOrder"`
	} `json:"skills,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
