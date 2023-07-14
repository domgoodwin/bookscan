package cmd

import (
	"net/http"

	"github.com/domgoodwin/bookscan/notion"
	"github.com/gin-gonic/gin"
)

func handleAuthRedirect(c *gin.Context) {
	code := c.Query("code")

	accessToken, err := notion.GetToken(c, code)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}
