package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/domgoodwin/bookscan/book"
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
	r.PUT("/cache/update", handlePUTUpdateCache)
	r.GET("/cache/info", handleGETCacheInfo)
}

func handleGETLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	book, found, err := lookupISBN(isbn)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
	})
}

func lookupISBN(isbn string) (*book.Book, bool, error) {
	var err error
	book, found := s.CheckIfBookInCache(isbn)
	if !found {
		book, err = lookup.LookupISBN(isbn)
		if err != nil {
			return nil, false, err
		}
	}
	return book, found, nil
}

func handlePUTStore(c *gin.Context) {
	var req putBookRequest
	err := c.Bind(&req)
	if err != nil {
		errorResponse(c, err)
		return
	}

	book, found, err := lookupISBN(req.ISBN)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
	}

	url := ""
	// Only store in CSV if not found
	if !found {
		err = book.StoreInCSV()
		if err != nil {
			errorResponse(c, err)
			return
		}
		url, err = notion.AddBookToDatabase(c, book, req.NotionDatabaseID)
		if err != nil {
			errorResponse(c, err)
			return
		}
		s.StoreBook(book)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
		"notion_page":    url,
	})
}

func handlePUTUpdateCache(c *gin.Context) {
	var req updateCacheRequest
	err := c.Bind(&req)
	if err != nil {
		errorResponse(c, err)
		return
	}

	var length int
	if req.ClearCache {
		length = s.ClearCache()
	}

	err = s.LoadBooksFromNotion(c, req.NotionDatabaseID)
	if err != nil {
		errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"deleted_cache_items_count": length,
		"cache_size":                s.Length(),
	})
}

func handleGETCacheInfo(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"database_id": s.DatabaseID(),
		"cache_size":  s.Length(),
	})
}

func errorResponse(c *gin.Context, err error) {
	c.JSON(mapErrorToCode(err), gin.H{
		"error": err.Error(),
	})
}

type putBookRequest struct {
	ISBN             string `json:"isbn" binding:"required"`
	NotionDatabaseID string `json:"notion_database_id" binding:"required"`
}

type updateCacheRequest struct {
	NotionDatabaseID string `json:"notion_database_id"`
	ClearCache       bool   `json:"clear_cache"`
}

func mapErrorToCode(err error) int {
	if strings.Contains(err.Error(), "404") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
