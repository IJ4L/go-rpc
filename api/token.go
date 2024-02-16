package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	db "simplebank.com/db/sqlgen"
)

type renewAccsessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccsessTokenResponse struct {
	AccsesToken          string    `json:"access_token"`
	AccsesTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (server *Server) renewAccsessToken(ctx *gin.Context) {
	var req renewAccsessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err = errors.New("invalid refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err = errors.New("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err = errors.New("username does not match")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(refreshPayload.ExpiredAt) {
		err = errors.New("refresh token is expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accsesToken, accsessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, "", server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccsessTokenResponse{
		AccsesToken:          accsesToken,
		AccsesTokenExpiredAt: accsessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
