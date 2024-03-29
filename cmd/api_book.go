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
	book, found, err := lookupISBN(c, isbn)
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

func lookupISBN(c *gin.Context, isbn string) (*items.Book, bool, error) {
	var err error
	book, found := store.BookStore.CheckIfItemInCache(getContextValue(c, contextKeyNotionBooksDatabaseID), isbn)
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

	book, found, err := lookupISBN(c, req.ISBN)
	if err != nil {
		c.JSON(mapErrorToCode(err), gin.H{
			"error": err.Error(),
		})
	}

	url := ""
	if !found {
		url, err = notionClient(c).AddBookToDatabase(c, book)
		if err != nil {
			errorResponse(c, err)
			return
		}
		store.BookStore.StoreItem(getContextValue(c, contextKeyNotionBooksDatabaseID), book)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book":           book.FullInfoFields(),
		"already_stored": found,
		"notion_page":    url,
	})
}

type putBookRequest struct {
	ISBN string `json:"isbn" binding:"required"`
}
