package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/domgoodwin/bookscan/items"
	"github.com/domgoodwin/bookscan/lookup"
	"github.com/domgoodwin/bookscan/notion"
	"github.com/domgoodwin/bookscan/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
)

var port string

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
		store.SetupStore()
		r.Run(fmt.Sprintf("0.0.0.0:%s", port))
	},
}

func setupRoutes(r *gin.Engine) {
	r.GET("/book/lookup", handleGETBookLookup)
	r.PUT("/book/store", handlePUTBookStore)
	r.GET("/record/lookup", handleGETRecordLookup)
	r.PUT("/record/store", handlePUTRecordStore)
	r.PUT("/cache/update", handlePUTUpdateCache)
	r.GET("/cache/info", handleGETCacheInfo)
}

func handleGETBookLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	book, found, err := lookupISBN(isbn)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
	})
}

func lookupISBN(isbn string) (*items.Book, bool, error) {
	var err error
	book, found := store.BookStore.CheckIfItemInCache(isbn)
	if !found {
		book, err = lookup.LookupISBN(isbn)
		if err != nil {
			return nil, false, err
		}
	}
	return book, found, nil
}

func handleGETRecordLookup(c *gin.Context) {
	isbn := c.Query("barcode")
	record, found, err := lookupRecordBarcode(isbn)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           record.FullInfoFields(),
		"already_stored": found,
	})
}

func lookupRecordBarcode(barcode string) (*items.Record, bool, error) {
	var err error
	record, found := store.RecordStore.CheckIfItemInCache(barcode)
	if !found {
		logrus.Infof("record not found in cache: %v", barcode)
		record, err = lookup.LookupRecordBarcode(barcode)
		if err != nil {
			logrus.Error(err)
			return nil, false, err
		}
		logrus.Infof("record looked up: %v", record)
	}
	return record, found, nil
}

func handlePUTBookStore(c *gin.Context) {
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
		store.BookStore.StoreItem(book)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
		"notion_page":    url,
	})
}

func handlePUTRecordStore(c *gin.Context) {
	var req putRecordRequest
	err := c.Bind(&req)
	if err != nil {
		errorResponse(c, err)
		return
	}

	record, found, err := lookupRecordBarcode(req.Barcode)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	url := ""
	if !found {
		url, err = notion.AddRecordToDatabase(c, record, req.NotionDatabaseID)
		if err != nil {
			errorResponse(c, err)
			return
		}
		store.RecordStore.StoreItem(record)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"record":         record.FullInfoFields(),
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

	var bookLength, recordLength int
	if req.ClearBooksCache {
		bookLength = store.BookStore.ClearCache()
	}
	if req.ClearRecordsCache {
		recordLength = store.RecordStore.ClearCache()
	}

	err = store.BookStore.LoadBooksFromNotion(c, req.BooksNotionDatabaseID)
	if err != nil {
		errorResponse(c, err)
		return
	}

	err = store.RecordStore.LoadRecordsFromNotion(c, req.RecordsNotionDatabaseID)
	if err != nil {
		errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"books": map[string]interface{}{
			"deleted_cache_items_count": bookLength,
			"cache_size":                store.BookStore.Length(),
		},
		"records": map[string]interface{}{
			"deleted_cache_items_count": recordLength,
			"cache_size":                store.RecordStore.Length(),
		},
	})
}

func handleGETCacheInfo(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"books": map[string]interface{}{
			"database_id": store.BookStore.DatabaseID(),
			"cache_size":  store.BookStore.Length(),
		},
		"records": map[string]interface{}{
			"database_id": store.RecordStore.DatabaseID(),
			"cache_size":  store.RecordStore.Length(),
		},
	})
}

func errorResponse(c *gin.Context, err error) {
	c.JSON(mapErrorToCode(err), gin.H{
		"error": err.Error(),
	})
}

type putBookRequest struct {
	ISBN             string `json:"isbn" binding:"required"`
	NotionDatabaseID string `json:"notion_database_id"`
}

type putRecordRequest struct {
	Barcode          string `json:"barcode" binding:"required"`
	NotionDatabaseID string `json:"notion_database_id"`
}

type updateCacheRequest struct {
	BooksNotionDatabaseID   string `json:"books_notion_database_id"`
	RecordsNotionDatabaseID string `json:"records_notion_database_id"`
	ClearBooksCache         bool   `json:"clear_books_cache"`
	ClearRecordsCache       bool   `json:"clear_records_cache"`
}

func mapErrorToCode(err error) int {
	if strings.Contains(err.Error(), "404") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
