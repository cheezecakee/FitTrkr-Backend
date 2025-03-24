package handler

import (
	"context"
	"net/http"
)

func (cfg *Config) PostRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := cfg.Helper.GetBearerToken(r.Header)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	err = cfg.DB.RevokeRefreshToken(context.Background(), refreshToken)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Refresh Token revoked succesfully!")
}
