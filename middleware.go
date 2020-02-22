package main

import "net/http"

type middleware func(http.HandlerFunc) http.HandlerFunc

func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		last := final
		for i := len(mw) - 1; i >= 0; i-- {
			last = mw[i](last)
		}
		return func(w http.ResponseWriter, r *http.Request) {
			last(w, r)
		}
	}
}
