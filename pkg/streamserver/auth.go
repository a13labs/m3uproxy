package streamserver

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/a13labs/m3uproxy/pkg/auth"
)

func bearerAuth(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
		}

		if authParts[0] != "Bearer" {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		token := authParts[1]
		if !auth.VerifyToken(token) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		next(w, r)
	})
}

func basicAuth(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		if authParts[0] != "Basic" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		token, err := auth.CreateToken(credentials[0], credentials[1])
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		next(w, r)
	})
}

func adminAccess(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return bearerAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		authParts := strings.SplitN(authHeader, " ", 2)
		token := authParts[1]

		role, err := auth.GetRoleFromToken(token)
		if err != nil || role != "admin" {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}))
}
