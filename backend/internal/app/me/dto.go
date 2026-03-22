package me

import "time"

// InputDto DTO定義
type InputDto struct {
	Certifications []struct {
		Issuer *string `json:"issuer,omitempty"`
		Month  *int    `json:"month,omitempty"`
		Name   string  `json:"name"`
		Year   int     `json:"year"`
	} `json:"certifications"`
	Experiences []struct {
		Company   string  `json:"company"`
		EndYear   *int    `json:"endYear,omitempty"`
		StartYear int     `json:"startYear"`
		URL       *string `json:"url,omitempty"`
	} `json:"experiences"`
	Likes []string `json:"likes"`
	Links []struct {
		Label    *string `json:"label,omitempty"`
		Platform string  `json:"platform"`
		URL      string  `json:"url"`
	} `json:"links"`
	Location    *string `json:"location"`
	DisplayName string  `json:"display"`
	DisplayJa   *string `json:"displayJa,omitempty"`
	Role        *string `json:"role"`
	Skills      []struct {
		Category  string   `json:"category"`
		Items     []string `json:"items"`
		SortOrder int      `json:"sortOrder"`
	} `json:"skills"`
}

// OutputDto DTO定義
type OutputDto struct {
	Certifications []struct {
		Issuer *string `json:"issuer,omitempty"`
		Month  *int    `json:"month,omitempty"`
		Name   string  `json:"name"`
		Year   int     `json:"year"`
	} `json:"certifications,omitempty"`
	Experiences []struct {
		Company   string  `json:"company"`
		EndYear   *int    `json:"endYear,omitempty"`
		StartYear int     `json:"startYear"`
		URL       *string `json:"url,omitempty"`
	} `json:"experiences,omitempty"`
	Likes []string `json:"likes,omitempty"`
	Links []struct {
		Label    *string `json:"label,omitempty"`
		Platform string  `json:"platform"`
		URL      string  `json:"url"`
	} `json:"links,omitempty"`
	Location    string `json:"location,omitempty"`
	DisplayName string `json:"display"`
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
