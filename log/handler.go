package log

import (
	"fmt"
	"net/http"
)

// LevelHandler returns an http.HandlerFunc that sets the current log level
// from a "level" query parameter (debug, info, warn or warning, error),
// e.g. POST /admin/log-level?level=debug
//
// It's a plain http.HandlerFunc so it can be mounted directly with
// net/http, or adapted into another router (e.g. Gin's gin.WrapF) without
// pulling that router's dependency into this module.
func LevelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := r.URL.Query().Get("level")
		if raw == "" {
			http.Error(w, `missing "level" query parameter`, http.StatusBadRequest)
			return
		}

		level, ok := ParseLevel(raw)
		if !ok {
			http.Error(w, fmt.Sprintf("unknown level %q", raw), http.StatusBadRequest)
			return
		}

		SetLevel(level)
		w.WriteHeader(http.StatusNoContent)
	}
}
