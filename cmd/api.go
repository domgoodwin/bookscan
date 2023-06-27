package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/domgoodwin/bookscan/lookup"
	"github.com/domgoodwin/bookscan/notion"
	"github.com/domgoodwin/bookscan/store"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
)

var port string
var s *store.Store

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.PersistentFlags().StringVar(&port, "port", "8443", "Port for API server to listen on")
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start an API server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		setupRoutes(r)
		notion.SetupClient()
		s = store.SetupStore()
		r.Run(fmt.Sprintf("0.0.0.0:%s", port))
	},
}

func setupRoutes(r *gin.Engine) {
	r.GET("/lookup", handleGETLookup)
	r.PUT("/store", handlePUTStore)
}

func handleGETLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	book, err := lookup.LookupISBN(isbn)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, book.FullInfoFields())
}

func handlePUTStore(c *gin.Context) {
	var info bookInfo
	var isbn string
	if c.Bind(&info) == nil {
		isbn = info.ISBN
	}
	book, err := lookup.LookupISBN(isbn)
	if err != nil {
		errorResponse(c, err)
		return
	}

	found := s.StoreBook(book)
	// Only store in CSV if not found
	if !found {
		err = book.StoreInCSV()
		if err != nil {
			errorResponse(c, err)
			return
		}
	}
	url, err := notion.AddBookToDatabase(c, book)
	if err != nil {
		errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
		"notion_page":    url,
	})
}

func errorResponse(c *gin.Context, err error) {
	c.JSON(mapErrorToCode(err), gin.H{
		"error": err.Error(),
	})
}

type bookInfo struct {
	ISBN string `json:"isbn" binding:"required"`
}

func mapErrorToCode(err error) int {
	if strings.Contains(err.Error(), "404") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
