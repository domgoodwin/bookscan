package cmd

import (
	"net/http"

	"github.com/domgoodwin/bookscan/store"
	"github.com/gin-gonic/gin"
)

func handlePUTUpdateCache(c *gin.Context) {
	var req updateCacheRequest
	err := c.Bind(&req)
	if err != nil {
		errorResponse(c, err)
		return
	}

	var bookLength, recordLength int
	if req.ClearBooksCache {
		bookLength = store.BookStore.ClearCache(req.BooksNotionDatabaseID)
	}
	if req.ClearRecordsCache {
		recordLength = store.RecordStore.ClearCache(req.RecordsNotionDatabaseID)
	}

	// err = store.BookStore.LoadBooksFromNotion(c, req.BooksNotionDatabaseID)
	// if err != nil {
	// 	errorResponse(c, err)
	// 	return
	// }

	// err = store.RecordStore.LoadRecordsFromNotion(c, req.RecordsNotionDatabaseID)
	// if err != nil {
	// 	errorResponse(c, err)
	// 	return
	// }

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

type updateCacheRequest struct {
	BooksNotionDatabaseID   string `json:"books_notion_database_id"`
	RecordsNotionDatabaseID string `json:"records_notion_database_id"`
	ClearBooksCache         bool   `json:"clear_books_cache"`
	ClearRecordsCache       bool   `json:"clear_records_cache"`
}
