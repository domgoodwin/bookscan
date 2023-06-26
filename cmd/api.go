package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/domgoodwin/bookscan/lookup"
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
		r.GET("/lookup", func(c *gin.Context) {
			isbn := c.Query("isbn")
			book, err := lookup.LookupISBN(isbn)
			if err != nil {
				c.JSON(mapErrorToCode(err), gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, book.FullInfoFields())
		})
		r.Run(fmt.Sprintf("0.0.0.0:%s", port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	},
}

func mapErrorToCode(err error) int {
	if strings.Contains(err.Error(), "404") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
