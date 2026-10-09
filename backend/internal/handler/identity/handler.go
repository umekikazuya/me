package identity

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/obs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Handler struct {
	interactor app.Interactor
	tokenSrv   app.TokenService
}

const cookieName = "authState"

var conf = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRETS"),
	Endpoint:     github.Endpoint,
	RedirectURL:  "http://localhost:8050/auth/github/callback",
	Scopes:       []string{"read:user"},
}

func NewHandler(interactor app.Interactor, tokenSrv app.TokenService) *Handler {
	return &Handler{interactor: interactor, tokenSrv: tokenSrv}
}

func getState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (h *Handler) LoginFromGithub(
	w http.ResponseWriter,
	r *http.Request,
) {
	state := getState()
	http.SetCookie(
		w,
		&http.Cookie{Name: cookieName, Value: state, Quoted: false, Expires: time.Now().Add(20 * time.Minute), HttpOnly: true},
	)
	http.Redirect(
		w,
		r,
		conf.AuthCodeURL(state),
		http.StatusFound,
	)
}

func (h *Handler) CallbackGithub(
	w http.ResponseWriter,
	r *http.Request,
) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		errs.WriteProblem(w, r, errs.ErrBadRequest)
		return
	}
	if c.Value != r.URL.Query().Get("state") {
		errs.WriteProblem(w, r, errs.ErrBadRequest)
		return
	}
	slog.InfoContext(r.Context(), "a", "r.URL.Query().Get(code)", r.URL.Query().Get("code"))
	t, err := conf.Exchange(
		r.Context(),
		r.URL.Query().Get("code"),
	)
	if err != nil {
		errs.WriteProblem(w, r, errs.New(errs.ErrBadRequest, err.Error()))
		return
	}
	client := conf.Client(r.Context(), t)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		errs.WriteProblem(w, r, errs.WrapInternal("データ取得エラー", err))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	var values map[string]any
	json.Unmarshal(body, &values)
	slog.InfoContext(
		r.Context(),
		"debug",
		"value",
		values["id"],
	)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	rtCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	err = h.interactor.Logout(
		r.Context(),
		app.InputLogoutDto{
			IdentityID: identityID,
			RT:         rtCookie.Value,
		},
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RevokeSessions(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	var input app.InputRevokeAllSessionsDto
	input.IdentityID = identityID
	err := h.interactor.RevokeAllSessions(
		r.Context(),
		input,
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	rtCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	input := app.InputRefreshTokensDto{
		IdentityID: identityID,
		RT:         rtCookie.Value,
	}
	out, err := h.interactor.RefreshTokens(
		r.Context(),
		input,
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	setTokenCookies(w, out.AT, out.RT)
	w.WriteHeader(http.StatusNoContent)
}
