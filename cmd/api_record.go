package cmd

import (
	"net/http"

	"github.com/domgoodwin/bookscan/items"
	"github.com/domgoodwin/bookscan/lookup"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
		"record":         record.FullInfoFields(),
		"already_stored": found,
	})
}

func lookupRecordBarcode(barcode string) (*items.Record, bool, error) {
	var err error
	// record, found := store.RecordStore.CheckIfItemInCache(notionDatabaseID, barcode)
	var record *items.Record
	found := false
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
		url, err = notionClient(c).AddRecordToDatabase(c, record)
		if err != nil {
			errorResponse(c, err)
			return
		}
		// store.RecordStore.StoreItem(record)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"record":         record.FullInfoFields(),
		"already_stored": found,
		"notion_page":    url,
	})
}

type putRecordRequest struct {
	Barcode string `json:"barcode" binding:"required"`
}
