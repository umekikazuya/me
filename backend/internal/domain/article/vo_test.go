package article

import (
	"testing"
	"time"
)

// --- newID ---

func Test_newID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    id
		wantErr bool
	}{
		{name: "valid id", input: "b0adacee33b2774d7089", want: id{value: "b0adacee33b2774d7089"}, wantErr: false},
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newID(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newID() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newID() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("newID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- newTitle ---

func Test_newTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    title
		wantErr bool
	}{
		{name: "valid title", input: "Factory Methodパターン入門", want: title{value: "Factory Methodパターン入門"}, wantErr: false},
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTitle(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newTitle() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newTitle() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("newTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- newURL ---

func Test_newURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    url
		wantErr bool
	}{
		{name: "valid url", input: "https://qiita.com/umekikazuya/items/b0adacee33b2774d7089", want: url{value: "https://qiita.com/umekikazuya/items/b0adacee33b2774d7089"}, wantErr: false},
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newURL(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newURL() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newURL() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("newURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- newPlatform ---

func Test_newPlatform(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    platform
		wantErr bool
	}{
		{name: "qiita", input: "qiita", want: platform{value: "qiita"}, wantErr: false},
		{name: "zenn", input: "zenn", want: platform{value: "zenn"}, wantErr: false},
		{name: "mochiya", input: "mochiya", want: platform{value: "mochiya"}, wantErr: false},
		{name: "note", input: "note", want: platform{value: "note"}, wantErr: false},
		{name: "unknown platform", input: "twitter", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPlatform(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newPlatform() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newPlatform() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("newPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- newPublishedAt ---

func Test_newPublishedAt(t *testing.T) {
	now := time.Now()
	past := now.AddDate(0, 0, -1)
	future := now.AddDate(0, 0, 1)

	tests := []struct {
		name    string
		input   time.Time
		want    publishedAt
		wantErr bool
	}{
		{name: "past date", input: past, want: publishedAt{value: past}, wantErr: false},
		{name: "current time", input: now, want: publishedAt{value: now}, wantErr: false},
		{name: "future date", input: future, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPublishedAt(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newPublishedAt() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newPublishedAt() succeeded unexpectedly")
			}
			if !got.value.Equal(tt.want.value) {
				t.Errorf("newPublishedAt() = %v, want %v", got.value, tt.want.value)
			}
		})
	}
}

// --- newArticleUpdatedAt ---

func Test_newArticleUpdatedAt(t *testing.T) {
	now := time.Now()
	past := now.AddDate(0, 0, -1)
	future := now.AddDate(0, 0, 1)

	tests := []struct {
		name    string
		input   time.Time
		want    articleUpdatedAt
		wantErr bool
	}{
		{name: "past date", input: past, want: articleUpdatedAt{value: past}, wantErr: false},
		{name: "current time", input: now, want: articleUpdatedAt{value: now}, wantErr: false},
		{name: "future date", input: future, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newArticleUpdatedAt(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("newArticleUpdatedAt() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newArticleUpdatedAt() succeeded unexpectedly")
			}
			if !got.value.Equal(tt.want.value) {
				t.Errorf("newArticleUpdatedAt() = %v, want %v", got.value, tt.want.value)
			}
		})
	}
}

// --- newIsActive ---

func Test_newIsActive(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  isActive
	}{
		{name: "true", input: true, want: isActive{value: true}},
		{name: "false", input: false, want: isActive{value: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newIsActive(tt.input)
			if err != nil {
				t.Errorf("newIsActive() unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("newIsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
