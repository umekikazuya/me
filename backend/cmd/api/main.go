package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	// ロガー初期化
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ルーターを初期化
	r := http.NewServeMux()

	// ヘルスチェック
	r.HandleFunc(
		"GET /up",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)
			if err := json.NewEncoder(w).Encode(
				struct {
					Status  string `json:"status"`
					Message string `json:"message"`
				}{Status: "ok", Message: "Server is running"},
			); err != nil {
				slog.Error("JSONエンコードエラー", "error", err)
				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}
		},
	)

	slog.Info("サーバーを起動します")

	// サーバー起動
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("起動エラー", "error", err)
		os.Exit(1) // TODO: SIGINT/SIGTERM シグナルを受信した際のグレースフルシャットダウン実装
	}
}