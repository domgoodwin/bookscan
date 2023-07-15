package cmd

import (
	"net/http"
	"os"

	"github.com/domgoodwin/bookscan/notion"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func handleAuth(c *gin.Context) {
	c.Redirect(http.StatusFound, os.Getenv("NOTION_AUTH_URL"))
}

func handleAuthRedirect(c *gin.Context) {
	code := c.Query("code")

	accessToken, err := notion.GetToken(c, code)
	if err != nil {
		logrus.Error(err)
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}
