package identity

import (
	"errors"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --- Type ---
type (
	identityID   struct{ value uuid.UUID }
	email        struct{ value string }
	password     struct{ value string }
	passwordHash struct{ value []byte }

	tokenHash struct{ value string }
	status    struct{ value string }
)

// --- Enum ---
var (
	statusActive  = status{value: "active"}
	statusRevoked = status{value: "revoked"}
)

// NewEmail はemailのコンストラクタ
func NewEmail(input string) (email, error) {
	return email{value: input}, nil
}

// NewPasswordHash はRepository to Domain時のみ利用想定
func NewPasswordHash(hash []byte) (passwordHash, error) {
	if len(hash) == 0 {
		return passwordHash{}, errors.New("パスワードが不正です")
	}
	return passwordHash{value: hash}, nil
}

// Password のコンストラクタ
func NewPassword(input string) (password, error) {
	if err := (password{value: input}).Validate(); err != nil {
		return password{}, err
	}
	return password{value: input}, nil
}

// --- Getter ---

func (vo identityID) Value() uuid.UUID {
	return vo.value
}

func (vo email) Value() string {
	return vo.value
}

func (vo passwordHash) Value() []byte {
	return vo.value
}

// Value は Password の値を返す
func (vo password) Value() string {
	return vo.value
}

// --- Validate ---

// Validate はパスワードの値を検証
func (vo password) Validate() error {
	if len(vo.Value()) < 8 {
		return errors.New("パスワードが不正です")
	}
	if len(vo.Value()) > 72 {
		return errors.New("パスワードが不正です")
	}
	// ここで、パスワードが大文字、小文字を含むかをチェック
	if !containsUppercase(vo.Value()) || !containsLowercase(vo.Value()) {
		return errors.New("パスワードが不正です")
	}
	return nil
}

// 大文字を含むかチェック
func containsUppercase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

// 小文字を含むかチェック
func containsLowercase(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}
	return false
}

// --- 振る舞い ---

// HashPassword はパスワードをハッシュ化
func (vo password) HashPassword() (passwordHash, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(vo.Value()), bcrypt.DefaultCost)
	if err != nil {
		return passwordHash{value: nil}, err
	}
	return passwordHash{value: h}, nil
}
