package router

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

type Option func(*Router) error

func New() *Router {
	return &Router{
		mux: &http.ServeMux{},
	}
}

func (r *Router) Build() http.Handler {
	return r.mux
}
