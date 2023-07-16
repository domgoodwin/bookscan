package cmd

import (
	"net/http"

	"github.com/domgoodwin/bookscan/items"
	"github.com/domgoodwin/bookscan/lookup"
	"github.com/domgoodwin/bookscan/store"
	"github.com/gin-gonic/gin"
)

func handleGETBookLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	dbID := c.Query("database_id")
	book, found, err := lookupISBN(dbID, isbn)
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

func lookupISBN(notionDatabaseID, isbn string) (*items.Book, bool, error) {
	var err error
	book, found := store.BookStore.CheckIfItemInCache(notionDatabaseID, isbn)
	if !found {
		book, err = lookup.LookupISBN(isbn)
		if err != nil {
			return nil, false, err
		}
	}
	return book, found, nil
}

func handlePUTBookStore(c *gin.Context) {
	var req putBookRequest
	err := c.Bind(&req)
	if err != nil {
		errorResponse(c, err)
		return
	}

	book, found, err := lookupISBN(req.NotionDatabaseID, req.ISBN)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
	}

	url := ""
	if !found {
		url, err = notionClient(c).AddBookToDatabase(c, book, req.NotionDatabaseID)
		if err != nil {
			errorResponse(c, err)
			return
		}
		store.BookStore.StoreItem(req.NotionDatabaseID, book)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
		"notion_page":    url,
	})
}

type putBookRequest struct {
	ISBN             string `json:"isbn" binding:"required"`
	NotionDatabaseID string `json:"notion_database_id"`
}
