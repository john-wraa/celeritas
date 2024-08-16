package middleware

import "net/http"

func (m *Middleware) AuthToken(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.Models.Tokens.AuthenticateToken(r)
		if err != nil {
			var payload struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}

			_ = m.App.WriteJSON(w, http.StatusUnauthorized, payload)
		}
	})
}
