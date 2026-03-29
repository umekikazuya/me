package token_test

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	appidentity "github.com/umekikazuya/me/internal/app/identity"
	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/internal/infra/token"
)

const (
	testSecret   = "test-secret-key-32-bytes-minimum!"
	testATExpiry = 15 * time.Minute
)

// --- helpers ---

func newSvc(t *testing.T) *token.JWTTokenService {
	t.Helper()
	return token.NewJWTTokenService(testSecret, testATExpiry)
}

// bcrypt を避けるため Reconstruct で Identity を生成
func mustNewTestIdentity(t *testing.T) domain.Identity {
	t.Helper()
	idn, err := domain.ReconstructIdentity(
		domain.ReconstructIdentityInput{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: []byte("$2a$10$aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	)
	if err != nil {
		t.Fatalf("mustNewTestIdentity: %v", err)
	}
	return *idn
}

// 期限切れトークンを直接生成（sleep なし）
func makeExpiredToken(t *testing.T, secret, identityID string) string {
	t.Helper()
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.MapClaims{
		"sub": identityID,
		"exp": time.Now().Add(-time.Hour).Unix(),
	})
	signed, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("makeExpiredToken: %v", err)
	}
	return signed
}

// 署名部分を別の文字列に丸ごと置換して改ざんする
func replaceSignature(t *testing.T, validToken, fakeSig string) string {
	t.Helper()
	parts := strings.Split(validToken, ".")
	if len(parts) != 3 {
		t.Fatalf("not a 3-part JWT: %q", validToken)
	}
	return parts[0] + "." + parts[1] + "." + fakeSig
}

// alg:none 攻撃トークンを手動生成（署名なし）
func makeNoneAlgToken(identityID string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString(
		[]byte(fmt.Sprintf(`{"sub":%q,"exp":%d}`, identityID, time.Now().Add(time.Hour).Unix())),
	)
	return header + "." + payload + "."
}

// 別シークレットで AT を生成する
func mustGenerateWithSecret(t *testing.T, identity domain.Identity, secret string) string {
	t.Helper()
	svc := token.NewJWTTokenService(secret, testATExpiry)
	tok, err := svc.GenerateAT(context.Background(), identity)
	if err != nil {
		t.Fatalf("mustGenerateWithSecret: %v", err)
	}
	return tok
}

// --- GenerateAT ---

func TestJWTTokenService_GenerateAT(t *testing.T) {
	t.Parallel()
	svc := newSvc(t)
	ctx := context.Background()

	t.Run("JWT形式: header.payload.signature の3セグメント", func(t *testing.T) {
		t.Parallel()
		tok, err := svc.GenerateAT(ctx, mustNewTestIdentity(t))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if parts := strings.Split(tok, "."); len(parts) != 3 {
			t.Errorf("token segments = %d, want 3 (header.payload.signature)", len(parts))
		}
	})

	t.Run("sub claim = identity.ID()", func(t *testing.T) {
		t.Parallel()
		identity := mustNewTestIdentity(t)
		tok, err := svc.GenerateAT(ctx, identity)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		parsed, _, err := gojwt.NewParser().ParseUnverified(tok, gojwt.MapClaims{})
		if err != nil {
			t.Fatalf("ParseUnverified: %v", err)
		}
		sub, err := parsed.Claims.GetSubject()
		if err != nil {
			t.Fatalf("GetSubject: %v", err)
		}
		if sub != identity.ID() {
			t.Errorf("sub = %q, want %q", sub, identity.ID())
		}
	})

	t.Run("exp = issuedAt + atExpiry の範囲内", func(t *testing.T) {
		t.Parallel()
		identity := mustNewTestIdentity(t)
		before := time.Now()
		tok, err := svc.GenerateAT(ctx, identity)
		after := time.Now()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		parsed, _, err := gojwt.NewParser().ParseUnverified(tok, gojwt.MapClaims{})
		if err != nil {
			t.Fatalf("ParseUnverified: %v", err)
		}
		exp, err := parsed.Claims.GetExpirationTime()
		if err != nil {
			t.Fatalf("no exp claim: %v", err)
		}
		// JWT の exp は Unix 秒単位なので wantMin も秒に丸める
		wantMin := before.Add(testATExpiry).Truncate(time.Second)
		wantMax := after.Add(testATExpiry).Add(time.Second)
		if exp.Before(wantMin) || exp.After(wantMax) {
			t.Errorf("exp %v outside [%v, %v]", exp.Time, wantMin, wantMax)
		}
	})

	t.Run("identity が異なれば別トークン", func(t *testing.T) {
		t.Parallel()
		tok1, _ := svc.GenerateAT(ctx, mustNewTestIdentity(t))
		tok2, _ := svc.GenerateAT(ctx, mustNewTestIdentity(t))
		if tok1 == tok2 {
			t.Error("different identities must produce different tokens")
		}
	})
}

// --- GenerateRT ---

func TestJWTTokenService_GenerateRT(t *testing.T) {
	t.Parallel()
	svc := newSvc(t)
	ctx := context.Background()

	t.Run("非空文字列を返す", func(t *testing.T) {
		t.Parallel()
		rt, err := svc.GenerateRT(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rt == "" {
			t.Error("expected non-empty RT")
		}
	})

	t.Run("10回連続生成して全て異なる値（uniqueness）", func(t *testing.T) {
		t.Parallel()
		seen := make(map[string]struct{}, 10)
		for i := range 10 {
			rt, err := svc.GenerateRT(ctx)
			if err != nil {
				t.Fatalf("call %d: %v", i, err)
			}
			if _, dup := seen[rt]; dup {
				t.Errorf("duplicate RT on call %d", i)
			}
			seen[rt] = struct{}{}
		}
	})

	t.Run("32バイト以上の entropy（base64url で43文字以上）", func(t *testing.T) {
		t.Parallel()
		rt, _ := svc.GenerateRT(ctx)
		if len(rt) < 43 {
			t.Errorf("RT length = %d, want >= 43 (32 bytes base64url)", len(rt))
		}
	})
}

// --- Hash ---

func TestJWTTokenService_Hash(t *testing.T) {
	t.Parallel()
	svc := newSvc(t)
	ctx := context.Background()

	t.Run("空文字はエラー", func(t *testing.T) {
		t.Parallel()
		_, err := svc.Hash(ctx, "")
		if err == nil {
			t.Error("Hash('') should return error")
		}
	})

	// 様々な入力で SHA-256 フォーマット（64文字小文字 hex）を確認
	formatCases := []struct {
		name  string
		input string
	}{
		{"短いトークン", "a"},
		{"典型的な RT 長", strings.Repeat("x", 43)},
		{"記号を含む文字列", "token/with+chars="},
		{"UUID 形式", uuid.New().String()},
	}
	for _, tc := range formatCases {
		t.Run("SHA-256 hex フォーマット: "+tc.name, func(t *testing.T) {
			t.Parallel()
			hash, err := svc.Hash(ctx, tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(hash) != 64 {
				t.Errorf("hash length = %d, want 64", len(hash))
			}
			for _, c := range hash {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
					t.Errorf("non-hex char %q in hash", c)
					break
				}
			}
		})
	}

	t.Run("決定的: 同じ入力は常に同じ出力", func(t *testing.T) {
		t.Parallel()
		h1, _ := svc.Hash(ctx, "token-abc")
		h2, _ := svc.Hash(ctx, "token-abc")
		if h1 != h2 {
			t.Error("Hash must be deterministic")
		}
	})

	t.Run("異なる入力は異なるハッシュ", func(t *testing.T) {
		t.Parallel()
		h1, _ := svc.Hash(ctx, "token-a")
		h2, _ := svc.Hash(ctx, "token-b")
		if h1 == h2 {
			t.Error("different inputs must produce different hashes")
		}
	})

	t.Run("ハッシュは平文と一致しない（非可逆）", func(t *testing.T) {
		t.Parallel()
		in := "plaintexttoken"
		hash, _ := svc.Hash(ctx, in)
		if hash == in {
			t.Error("hash must not equal plaintext")
		}
	})
}

// --- Validate ---

func TestJWTTokenService_Validate_Valid(t *testing.T) {
	t.Parallel()
	svc := newSvc(t)
	ctx := context.Background()

	tok, err := svc.GenerateAT(ctx, mustNewTestIdentity(t))
	if err != nil {
		t.Fatalf("GenerateAT: %v", err)
	}
	if err := svc.Validate(ctx, tok); err != nil {
		t.Errorf("Validate(valid AT) = %v, want nil", err)
	}
}

func TestJWTTokenService_Validate_InvalidCases(t *testing.T) {
	t.Parallel()
	svc := newSvc(t)
	ctx := context.Background()

	identity := mustNewTestIdentity(t)
	validTok, _ := svc.GenerateAT(ctx, identity)

	cases := []struct {
		name  string
		token string
		want  error
	}{
		{
			name:  "空文字",
			token: "",
			want:  appidentity.ErrTokenInvalid,
		},
		{
			name:  "期限切れ",
			token: makeExpiredToken(t, testSecret, identity.ID()),
			want:  appidentity.ErrTokenExpired,
		},
		{
			name:  "署名改ざん",
			token: replaceSignature(t, validTok, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			want:  appidentity.ErrTokenInvalid,
		},
		{
			name:  "別シークレットで署名",
			token: mustGenerateWithSecret(t, identity, "completely-different-secret!!"),
			want:  appidentity.ErrTokenInvalid,
		},
		{
			name:  "JWT でない任意文字列",
			token: "not-a-jwt-at-all",
			want:  appidentity.ErrTokenInvalid,
		},
		{
			name:  "alg:none 攻撃（署名なし）",
			token: makeNoneAlgToken(identity.ID()),
			want:  appidentity.ErrTokenInvalid,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := svc.Validate(ctx, tc.token)
			if !errors.Is(err, tc.want) {
				t.Errorf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}
