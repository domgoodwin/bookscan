package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var userLicense string

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bookscan.yaml)")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.SetDefault("author", "Dom Goodwin git@dgood.win")
	viper.SetDefault("license", "apache")
}

var rootCmd = &cobra.Command{
	Use:   "bookscan",
	Short: "Bookscan is a tool to interact with your book library.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
