package cmd

import (
	"fmt"

	"github.com/domgoodwin/bookscan/lookup"
	"github.com/spf13/cobra"
)

var source string

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.PersistentFlags().StringVar(&source, "source", "openlibrary", "ISBN of book")
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a single book",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		book, err := lookup.LookupISBN(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Println(book.FullInfo())
	},
}
