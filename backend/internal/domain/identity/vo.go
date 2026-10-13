package identity

import (
	"errors"

	"github.com/google/uuid"
)

// --- Type ---

type (
	id struct{ value uuid.UUID }

	tokenHash struct{ value string }
	status    struct{ value string }
)

// --- Enum ---
var (
	statusActive  = status{value: "active"}
	statusRevoked = status{value: "revoked"}
)

// newIdentityID はidentityIDのコンストラクタ
func newIdentityID(input uuid.UUID) id {
	return id{value: input}
}

// NewTokenHash はtokenHashのコンストラクタ
func NewTokenHash(input string) (tokenHash, error) {
	if input == "" {
		return tokenHash{}, errors.New("tokenHashが不正です")
	}
	return tokenHash{value: input}, nil
}

// --- Getter ---

func (vo id) Value() string {
	return vo.value.String()
}

// Value は tokenHash の値を返す
func (vo tokenHash) Value() string {
	return vo.value
}

// Value は status の値を返す
func (vo status) Value() string {
	return vo.value
}
