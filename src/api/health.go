package api

import (
	"net/http"

	"github.com/birabittoh/go-lift/src/database"
)

type AutheliaUserInfo struct {
	DisplayName string   `json:"display_name"`
	Emails      []string `json:"emails"`
	Method      string   `json:"method"`
	HasTOTP     bool     `json:"has_totp"`
	HasWebAuthn bool     `json:"has_webauthn"`
	HasDuo      bool     `json:"has_duo"`
}

type AutheliaUserInfoResponse struct {
	Status string           `json:"status"`
	Data   AutheliaUserInfo `json:"data"`
}

var mockAutheliaResponse = AutheliaUserInfoResponse{
	Status: "OK",
	Data: AutheliaUserInfo{
		DisplayName: "Admin",
		Emails:      []string{"a***n@admin.com"},
		Method:      "totp",
	},
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"message": "pong"})
}
func connectionHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB.DB()
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to get database connection:", err.Error())
			return
		}

		err = sqlDB.Ping()
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Database ping failed:", err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, map[string]string{"message": "Database connection is healthy"})
	}
}

func mockAutheliaHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, mockAutheliaResponse)
}
