package password

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"strings"

	"github.com/umekikazuya/me/internal/app/port"
	"golang.org/x/crypto/argon2"
)

type Argon2PasswordManager struct{}

// Hash implements [port.PasswordManager].
//
// time: 1
// memory: 64MB
// threads: 4
// keyLen: 32
func (a *Argon2PasswordManager) Hash(
	ctx context.Context,
	input string,
) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	hash := argon2.IDKey(
		[]byte(input),
		salt,
		1,
		64*1024,
		4,
		32,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Appendf(
		nil,
		"$argon2id$v=19$m=65536,t=1,p=4$%s$%s",
		b64Salt,
		b64Hash,
	), nil
}

// Verify implements [port.PasswordManager].
func (a *Argon2PasswordManager) Verify(ctx context.Context, hashedPassword string, plainPassword string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return ErrInvalidHashFormat
	}

	// 2. アルゴリズムの確認
	if parts[1] != "argon2id" {
		return ErrInvalidHashFormat
	}

	// 3. バージョンの検証
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return ErrInvalidHashFormat
	}
	if version != argon2.Version {
		return ErrIncompatibleVersion
	}

	// 4. パラメータ (memory, time, threads) の抽出
	var memory uint32
	var time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return ErrInvalidHashFormat
	}

	// 5. ソルトとハッシュ値のBase64デコード
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidHashFormat
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidHashFormat
	}

	// 6. 抽出したパラメータを使って入力パスワードを再度ハッシュ化
	// 注意: 最後の引数 (keyLen) は、保存されているハッシュ値の長さに合わせる
	keyLen := len(hash)
	if keyLen > math.MaxUint32 {
		return ErrInvalidHashFormat
	}
	computedHash := argon2.IDKey(
		[]byte(plainPassword),
		salt,
		time,
		memory,
		threads,
		uint32(keyLen),
	)

	// 7. 定数時間比較（タイミング攻撃対策）
	if subtle.ConstantTimeCompare(hash, computedHash) != 1 {
		return ErrMismatchedHashAndPassword
	}

	return nil
}

var _ port.PasswordManager = (*Argon2PasswordManager)(nil)
