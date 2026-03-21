package me

import "testing"

func Test_DisplayName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid display name",
			args:    args{input: "abc def"},
			want:    "abc def",
			wantErr: false,
		},
		{
			name:    "empty input",
			args:    args{input: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "space only input",
			args:    args{input: "  "},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newDisplayName(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDisplayName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("newDisplayName().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func Test_DisplayNameJa(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid japanese display name",
			args:    args{input: "田中 太郎"},
			want:    "田中 太郎",
			wantErr: false,
		},
		{
			name:    "empty input",
			args:    args{input: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newDisplayNameJa(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDisplayNameJa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("newDisplayNameJa().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func Test_Role(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid role",
			args:    args{input: "Software Engineer"},
			want:    "Software Engineer",
			wantErr: false,
		},
		{
			name:    "empty input",
			args:    args{input: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRole(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("newRole().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func Test_Location(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid location",
			args:    args{input: "Tokyo, Japan"},
			want:    "Tokyo, Japan",
			wantErr: false,
		},
		{
			name:    "empty input",
			args:    args{input: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newLocation(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("newLocation().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func Test_Like(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid like",
			args:    args{input: "Go"},
			want:    "Go",
			wantErr: false,
		},
		{
			name:    "empty input",
			args:    args{input: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newLike(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("newLike().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func Test_ValidatePositiveInt(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"positive", 1, false},
		{"large positive", 2026, false},
		{"zero", 0, true},
		{"negative", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePositiveInt(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validatePositiveInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ValidateMonth(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"min valid", 1, false},
		{"max valid", 12, false},
		{"zero", 0, true},
		{"too large", 13, true},
		{"negative", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMonth(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateMonth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ValidateNonEmpty(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid string", "hello", false},
		{"valid string with space", "hello world", false},
		{"empty string", "", true},
		{"space only", "   ", true},
		{"tab only", "\t", true},
		{"newline only", "\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateNonEmpty(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateNonEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
