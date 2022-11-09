package delivery

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	userId, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("userId", userId)
}

func getUserId(c *gin.Context) (int, bool) {
	userId, ok := c.Get("userId")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "couldn't find user id")
		return 0, false
	}

	userIdInt, ok := userId.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "couldn't find user id")
		return 0, false
	}
	return userIdInt, true
}
