package cmd

import (
	"github.com/domgoodwin/bookscan/barcode"
	"github.com/spf13/cobra"
)


func init() {
	rootCmd.AddCommand(scanCmd)
  }
  
  var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Listen for books to scan",
	Run: func(cmd *cobra.Command, args []string) {
	  barcode.WaitForScan()
	},
  }