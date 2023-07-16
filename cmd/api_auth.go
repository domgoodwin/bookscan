package cmd

import (
	"fmt"
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
	redirectURI := c.Query("redirect_uri")

	userID, apiToken, err := notion.GetToken(c, code, redirectURI)
	if err != nil {
		logrus.Error(err)
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(
		"bookscan://auth?api_token=%v&user_id=%v",
		apiToken,
		userID,
	))
}
