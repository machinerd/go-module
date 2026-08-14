package log

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LevelHandler returns an http.HandlerFunc for reading and changing the
// current log level:
//
//	GET  /admin/log-level             -> {"level":"info"}
//	POST /admin/log-level?level=debug -> 204
//
// It's a plain http.HandlerFunc so it can be mounted directly with
// net/http, or adapted into another router (e.g. Gin's gin.WrapF) without
// pulling that router's dependency into this module.
func LevelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			level, ok := CurrentLevel()
			if !ok {
				http.Error(w, "current logger does not report its level", http.StatusNotImplemented)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"level": level.String()})

		case http.MethodPost:
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

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
