package main

import (
	"github.com/spf13/cobra"
)

var (
	cbrTargetURL     string
	cbrProfile       string
	cbrHeadless      bool
	cbrHeaders       []string
	cbrDefaultRegexp bool
	cbrRegexps       []string
)

var rootCmd = &cobra.Command{
	Use:   "gocrawldp",
	Short: "SPA crawler using CDP and chromium",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	rootCmd.Flags().StringVar(
		&cbrTargetURL,
		"url",
		"",
		"Target URL",
	)
	
	rootCmd.Flags().StringVar(
	  &cbrProfile,
	  "profile",
	  "",
	  "Path to chrome profile",
	)

	rootCmd.Flags().BoolVar(
		&cbrHeadless,
		"headless",
		false,
		"Run Chrome in headless mode",
	)

	rootCmd.Flags().StringArrayVar(
		&cbrHeaders,
		"header",
		nil,
		"Extra HTTP header (Name: Value)",
	)
	
	rootCmd.Flags().StringArrayVar(
	  &cbrRegexps,
	  "regexp",
	  nil,
	  "Extra regular expressions (Name: `pattern`)",
	)
	
	rootCmd.Flags().BoolVar(
	  &cbrDefaultRegexp,
	  "nodefrexp",
	  false,
	  "Omit loading of default regular expressions",
	)

	if err := rootCmd.MarkFlagRequired("url"); err != nil {
		panic(err)
	}
}
