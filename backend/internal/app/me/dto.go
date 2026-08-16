package me

// InputDto DTO定義
type InputDto struct {
	ID          string `json:"-"`
	DisplayName string `json:"displayName"         validate:"required"`
}

type (
	InputUpdateProfile struct {
		Location    string `json:"location,omitempty"`
		DisplayName string `json:"displayName"         validate:"required"`
		DisplayJa   string `json:"displayJa,omitempty"`
		Role        string `json:"role,omitempty"`
	}
	InputUpdateLinks []InputLink
	InputLink        struct {
		Label    string `json:"label,omitempty"`
		Platform string `json:"platform"            validate:"required"`
		URL      string `json:"url"                 validate:"required,url"`
	}
	InputUpdateLikes []string
)

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
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
