package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/umekikazuya/me/internal/app/eventhandler"
	"github.com/umekikazuya/me/internal/app/identity"
	appme "github.com/umekikazuya/me/internal/app/me"
	handleridentity "github.com/umekikazuya/me/internal/handler/identity"
	handlerme "github.com/umekikazuya/me/internal/handler/me"
	infraevent "github.com/umekikazuya/me/internal/infra/event"
	"github.com/umekikazuya/me/internal/infra/token"
	"github.com/umekikazuya/me/pkg/middleware"
)

func main() {
	ctx := context.Background()

	// ロガー初期化
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 具像実装の初期化
	meRepo, identityRepo, sessionRepo, err := setupRepo(ctx)
	if err != nil {
		slog.Error("インフラの初期化に失敗しました", "error", err)
		os.Exit(1)
	}
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		slog.Error("JWT_SECRET が未設定です")
		os.Exit(1)
	}
	tokenSrv := token.NewJWTTokenService(
		jwtSecret,
		15*time.Minute,
	)
	meInteractor := appme.NewInteractor(meRepo)
	meHandler := handlerme.NewHandler(meInteractor)

	bus := infraevent.NewLocalEventBus()
	bus.Register(eventhandler.NewIdentityRegisteredHandler(meInteractor))
	identityInteractor := identity.NewInteractor(identityRepo, sessionRepo, tokenSrv, bus)
	identityHandler := handleridentity.NewHandler(identityInteractor, tokenSrv)

	// ルーター初期化
	r := http.NewServeMux()

	// ヘルスチェック
	r.HandleFunc("GET /up", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct { //nolint:errcheck
			Status string `json:"status"`
		}{Status: "ok"})
	})

	// Me
	r.HandleFunc("GET /me", meHandler.Get)
	r.Handle("PUT /me", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(meHandler.Update),
		),
	))
	r.HandleFunc("POST /me", meHandler.Create) // TODO: 認証プロファイル作成時のイベントでMe集約がセットアップされるのが理想

	// --- Identity ---
	// login
	r.Handle("POST /auth/login", handleridentity.CSRFMiddleware(
		http.HandlerFunc(identityHandler.Login),
	),
	)
	// logout
	r.Handle("POST /auth/logout", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(identityHandler.Logout),
		),
	))
	// refresh TODO: https://github.com/umekikazuya/me/pull/33#discussion_r3017640414
	r.Handle("POST /auth/refresh", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(identityHandler.RefreshToken),
		),
	))
	// register
	r.Handle("POST /auth/register", handleridentity.CSRFMiddleware(
		http.HandlerFunc(identityHandler.Register),
	))
	// resetPassword
	r.Handle("PUT /auth/password", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(identityHandler.ResetPassword),
		),
	))
	// changeEmail
	r.Handle("PUT /auth/email", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(identityHandler.ChangeEmailAddress),
		),
	))

	slog.Info("サーバーを起動します")

	// サーバー起動
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           middleware.Logging(r),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("起動エラー", "error", err)
		os.Exit(1) // TODO: SIGINT/SIGTERM グレースフルシャットダウン
	}
}
