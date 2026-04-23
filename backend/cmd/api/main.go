package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/umekikazuya/me/internal/app/eventhandler"
	"github.com/umekikazuya/me/internal/app/identity"
	appme "github.com/umekikazuya/me/internal/app/me"
	handlerarticle "github.com/umekikazuya/me/internal/handler/article"
	handleridentity "github.com/umekikazuya/me/internal/handler/identity"
	handlerme "github.com/umekikazuya/me/internal/handler/me"
	infraevent "github.com/umekikazuya/me/internal/infra/event"
	"github.com/umekikazuya/me/internal/infra/token"
	"github.com/umekikazuya/me/pkg/middleware"
	"github.com/umekikazuya/me/pkg/obs"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 観測性基盤初期化: logs / traces / metrics を stdout に出す。
	// アクセスログはインフラ層 (API Gateway/ALB 等) の責務とし、アプリでは出さない。
	prov, shutdown, err := obs.Bootstrap(ctx, obs.Config{
		ServiceName:   "api",
		Level:         obs.ParseLevel(os.Getenv("LOG_LEVEL")),
		SensitiveKeys: []string{"password", "password_hash", "authorization", "cookie", "set-cookie", "token", "refresh_token"},
		AddSource:     true,
		EnableTraces:  true,
		EnableMetrics: true,
	})
	if err != nil {
		slog.Error("観測性基盤の初期化に失敗しました", "error", err)
		os.Exit(1)
	}
	defer func() { _ = shutdown(ctx) }()
	slog.SetDefault(prov.Logger)

	// 具像実装の初期化
	meRepo, identityRepo, sessionRepo, articleInteractor, err := setupRepo(ctx)
	if err != nil {
		slog.Error("インフラの初期化に失敗しました", "error", err)
		os.Exit(1)
	}
	articleHandler := handlerarticle.NewHandler(articleInteractor)
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

	dispatcher := infraevent.NewSyncEventDispatcher()
	dispatcher.Register(eventhandler.NewIdentityRegisteredHandler(meInteractor))
	identityInteractor := identity.NewInteractor(identityRepo, sessionRepo, tokenSrv, dispatcher)
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

	// Articles (public)
	r.HandleFunc("GET /articles", articleHandler.Search)
	r.HandleFunc("GET /articles/meta/tags", articleHandler.GetTagsAll)
	r.HandleFunc("GET /articles/meta/suggest", articleHandler.GetSuggests)

	// Articles (admin)
	r.Handle("POST /articles", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(articleHandler.Register),
		),
	))
	r.Handle("PUT /articles/{externalId}", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(articleHandler.Update),
		),
	))
	r.Handle("DELETE /articles/{externalId}", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(articleHandler.Remove),
		),
	))

	// Me
	r.HandleFunc("GET /me", meHandler.Get)
	r.Handle("PUT /me", handleridentity.CSRFMiddleware(
		identityHandler.AuthMiddleware(
			http.HandlerFunc(meHandler.Update),
		),
	))

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
	// middleware chain (外側 → 内側):
	//   RequestID (obs.WithRequestID で context に積む)
	//     → otelhttp (root span を作成、trace_id を context に載せる)
	//       → Recover (panic → 500 ProblemDetail + ERROR ログ、trace_id が自動で付く)
	//         → router
	srv := &http.Server{
		Addr: ":8080",
		Handler: middleware.RequestID(
			otelhttp.NewHandler(
				middleware.Recover(r),
				"api",
			),
		),
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

