package di

import (
	"encoding/json"
	"net/http"

	"github.com/umekikazuya/me/internal/handler/article"
	"github.com/umekikazuya/me/internal/handler/identity"
	"github.com/umekikazuya/me/internal/handler/me"
)

type Handlers struct {
	Me       me.Handler
	Article  article.Handler
	Identity identity.Handler
}

func NewRouter(handlers Handlers) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /up", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct { //nolint:errcheck
			Status string `json:"status"`
		}{Status: "ok"})
	})
	// Articles (public)
	r.HandleFunc(
		"GET /articles",
		handlers.Article.Search,
	)
	r.HandleFunc(
		"GET /articles/meta/tags",
		handlers.Article.GetTagsAll,
	)
	r.HandleFunc(
		"GET /articles/meta/suggest",
		handlers.Article.GetSuggests,
	)

	// Articles (admin)
	r.Handle(
		"POST /articles",
		identity.CSRFMiddleware(
			handlers.Identity.AuthMiddleware(
				http.HandlerFunc(handlers.Article.Register),
			),
		),
	)
	r.Handle(
		"PUT /articles/{externalId}",
		identity.CSRFMiddleware(
			handlers.Identity.AuthMiddleware(http.HandlerFunc(handlers.Article.Update)),
		),
	)
	r.Handle("DELETE /articles/{externalId}", identity.CSRFMiddleware(
		handlers.Identity.AuthMiddleware(
			http.HandlerFunc(handlers.Article.Remove),
		),
	))

	// Me
	r.HandleFunc("GET /me", handlers.Me.Get)
	r.Handle("PUT /me", identity.CSRFMiddleware(
		handlers.Identity.AuthMiddleware(
			http.HandlerFunc(handlers.Me.Update),
		),
	))

	// --- Identity ---
	// login
	r.Handle("POST /auth/login", identity.CSRFMiddleware(
		http.HandlerFunc(handlers.Identity.Login),
	),
	)
	// logout
	r.Handle("POST /auth/logout", identity.CSRFMiddleware(
		handlers.Identity.AuthMiddleware(
			http.HandlerFunc(handlers.Identity.Logout),
		),
	))
	// refresh TODO: https://github.com/umekikazuya/me/pull/33#discussion_r3017640414
	r.Handle("POST /auth/refresh", identity.CSRFMiddleware(
		handlers.Identity.AuthMiddleware(
			http.HandlerFunc(handlers.Identity.RefreshToken),
		),
	))
	// register
	r.Handle("POST /auth/register", identity.CSRFMiddleware(
		http.HandlerFunc(handlers.Identity.Register),
	))
	// resetPassword
	r.Handle("PUT /auth/password", identity.CSRFMiddleware(
		handlers.Identity.AuthMiddleware(
			http.HandlerFunc(handlers.Identity.ResetPassword),
		),
	))
	// changeEmail
	r.Handle("PUT /auth/email", identity.CSRFMiddleware(
		handlers.Identity.AuthMiddleware(
			http.HandlerFunc(handlers.Identity.ChangeEmailAddress),
		),
	))

	return r
}
