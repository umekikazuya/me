package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	appme "github.com/umekikazuya/me/internal/app/me"
	handlerme "github.com/umekikazuya/me/internal/handler/me"
	"github.com/umekikazuya/me/pkg/middleware"
)

func main() {
	ctx := context.Background()

	// ロガー初期化
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo, err := setupRepo(ctx)
	if err != nil {
		slog.Error("インフラの初期化に失敗しました", "error", err)
		os.Exit(1)
	}
	interactor := appme.NewInteractor(repo)
	me := handlerme.NewHandler(interactor)

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
	r.HandleFunc("GET /me", me.Get)
	r.HandleFunc("POST /me", me.Create)
	r.HandleFunc("PUT /me", me.Update)

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
