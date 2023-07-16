package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/domgoodwin/bookscan/database"
	"github.com/domgoodwin/bookscan/notion"
	"github.com/domgoodwin/bookscan/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
)

const (
	headerUserID          = "Bookscan-User-Id"
	headerAPIToken        = "Bookscan-Token"
	contextKeyNotionToken = "NOTION_TOKEN"
	contextKeyNotionPage  = "NOTION_PAGE"
	contextKeyUserID      = "USER_ID"
)

var port string
var tlsPort string

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.PersistentFlags().StringVar(&port, "port", "8443", "Port for API server to listen on")
	apiCmd.PersistentFlags().StringVar(&tlsPort, "tlsPort", "", "Port for API server to listen on with TLS certs")
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start an API server",
	Run: func(cmd *cobra.Command, args []string) {
		err := database.Setup()
		if err != nil {
			panic(err)
		}
		r := gin.Default()
		setupRoutes(r)
		store.SetupStore()

		// TLS handler if port set
		if tlsPort != "" {
			certFolder := os.Getenv("CERT_FOLDER")
			certName := os.Getenv("CERT_NAME")
			keyName := os.Getenv("KEY_NAME")
			logrus.Infof("Running TLS server: %v %v %v", certFolder, certName, keyName)
			r.RunTLS(fmt.Sprintf("0.0.0.0:%s", port), certFolder+certName, certFolder+keyName)
		}
		r.Run(fmt.Sprintf("0.0.0.0:%s", port))

	},
}

func setupRoutes(r *gin.Engine) {
	// Auth based groups (no auth check middleware)
	auth := r.Group("/auth")
	auth.GET("/redirect", handleAuthRedirect)
	auth.GET("/", handleAuth)

	d := r.Group("/")
	d.Use(checkAPIToken, getNotionToken)
	d.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"app": "bookscan"})
	})
	d.GET("/book/lookup", handleGETBookLookup)
	d.PUT("/book/store", handlePUTBookStore)
	d.GET("/record/lookup", handleGETRecordLookup)
	d.PUT("/record/store", handlePUTRecordStore)
	d.PUT("/cache/update", handlePUTUpdateCache)
	d.GET("/cache/info", handleGETCacheInfo)
}

func checkAPIToken(c *gin.Context) {
	userID := getHeader(c, headerUserID)
	token := getHeader(c, headerAPIToken)
	logrus.Debugf("Checking API token %v %v", userID, token)
	valid, err := database.CheckAPIToken(c, userID, token)
	if err != nil {
		logrus.Error(err)
		errorResponse(c, err)
		c.Abort()
		return
	}
	if !valid {
		logrus.Errorf("access code isn't valid")
		errorResponse(c, errors.New("invalid access code"))
		c.Abort()
		return
	}
	logrus.Debug("validated access code")

}

func getNotionToken(c *gin.Context) {
	userID := getHeader(c, headerUserID)
	// We've already authenticated here
	token, err := database.GetNotionTokenByUserID(c, userID)
	if err != nil {
		logrus.Error(err)
		errorResponse(c, err)
		c.Abort()
	}
	c.Set(contextKeyNotionToken, token.AccessToken)
	c.Set(contextKeyNotionPage, token.DuplicatedTemplateID)
	c.Set(contextKeyUserID, userID)
}

func getHeader(c *gin.Context, name string) string {
	values := c.Request.Header[name]
	if len(values) == 0 {
		logrus.Debugf("header: %v not found: %v", name, c.Request.Header)
		return ""
	}
	return values[0]
}

func errorResponse(c *gin.Context, err error) {
	c.JSON(mapErrorToCode(err), gin.H{
		"error": err.Error(),
	})
}

func mapErrorToCode(err error) int {
	if strings.Contains(err.Error(), "404") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func notionClient(c *gin.Context) *notion.NotionClient {
	userIDAny, _ := c.Get(contextKeyUserID)
	userID := userIDAny.(string)
	notionTokenAny, _ := c.Get(contextKeyNotionToken)
	notionToken := notionTokenAny.(string)
	notionPageAny, _ := c.Get(contextKeyNotionPage)
	notionPage := notionPageAny.(string)

	notionClient := notion.GetClient(notionToken)
	return &notion.NotionClient{notionClient, notionToken, userID, notionPage}

}
