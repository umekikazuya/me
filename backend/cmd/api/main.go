package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

// shutdownTimeout は SIGINT/SIGTERM 受信後にインフライトリクエストを捌き切る猶予。
// ALB / API Gateway のドレイン時間との整合を意識して設定する。
const shutdownTimeout = 30 * time.Second

func main() {
	// run に集約するのは os.Exit が defer をスキップするため。
	// 直接 main で os.Exit すると obs の shutdown が走らず traces/metrics が flush されない。
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

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
		return 1
	}
	// shutdown は長寿命な ctx に縛られないよう、呼び出し時に bounded な context を渡す。
	// ListenAndServe から戻った時点では元 ctx が cancel 済みの可能性があり、
	// その場合 tracer/meter の flush が即座に諦められてしまう。
	// signal 受信時に作成した共有 shutdownCtx を使うことで、HTTP と observability の
	// shutdown が同一タイムアウト予算を共有し、合計で shutdownTimeout 以内に収まる。
	var shutdownCtx context.Context
	var shutdownCancel context.CancelFunc
	defer func() {
		if shutdownCtx != nil {
			_ = shutdown(shutdownCtx)
		}
		if shutdownCancel != nil {
			shutdownCancel()
		}
	}()
	slog.SetDefault(prov.Logger)

	// 具像実装の初期化
	meRepo, identityRepo, sessionRepo, articleInteractor, err := setupRepo(ctx)
	if err != nil {
		slog.Error("インフラの初期化に失敗しました", "error", err)
		return 1
	}
	articleHandler := handlerarticle.NewHandler(articleInteractor)
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		slog.Error("JWT_SECRET が未設定です")
		return 1
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

	// signal 受信準備を ListenAndServe の goroutine 起動より前に行う。
	// 先に goroutine を起動すると、早期の SIGINT/SIGTERM がデフォルトハンドラで
	// 処理され、graceful shutdown が実行されないリスクがある。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// ListenAndServe は別 goroutine で動かし、main は signal 受信を待つ。
	// 起動直後に失敗 (port 競合など) した場合は srvErrCh から即座に戻る。
	srvErrCh := make(chan error, 1)
	go func() {
		srvErrCh <- srv.ListenAndServe()
	}()

	slog.Info("サーバーを起動します")

	select {
	case err := <-srvErrCh:
		// SIGINT/SIGTERM を受ける前に ListenAndServe が返った = 起動失敗。
		// Shutdown 経由の終了ではないので ErrServerClosed は来ない想定だが、念のため除外。
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("起動エラー", "error", err)
			return 1
		}
		return 0
	case sig := <-sigCh:
		slog.Info("シャットダウン開始", "signal", sig.String())
	}

	// HTTP shutdown と observability shutdown の両方で共有する単一の context を作成。
	// 合計で shutdownTimeout 以内に両方の shutdown を完了させる。
	shutdownCtx, shutdownCancel = context.WithTimeout(context.Background(), shutdownTimeout)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		// Shutdown が猶予内に完了できなかった = インフライトが残っている。
		// それでも ListenAndServe は Close 済みで戻るので、goroutine は解放される。
		slog.Error("サーバーシャットダウン失敗", "error", err)
	}
	// Shutdown 後は ListenAndServe が ErrServerClosed で戻る。goroutine のクローズ待ち。
	if err := <-srvErrCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("起動エラー", "error", err)
		return 1
	}
	slog.Info("サーバーを停止しました")
	return 0
}
