package main

import (
	"github.com/spf13/cobra"
)

type Config struct {
  TargetURL         string
  Profile           string
  Headless          bool
  Headers           []string
  SkipDefaultRegexp bool
  CustomRegexps     []string
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:   "gocrawldp",
	Short: "SPA crawler using CDP and chromium",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	rootCmd.Flags().StringVar(
		&cfg.TargetURL,
		"url",
		"",
		"Target URL",
	)
	
	rootCmd.Flags().StringVar(
	  &cfg.Profile,
	  "profile",
	  "",
	  "Path to chrome profile",
	)

	rootCmd.Flags().BoolVar(
		&cfg.Headless,
		"headless",
		false,
		"Run Chrome in headless mode",
	)

	rootCmd.Flags().StringArrayVar(
		&cfg.Headers,
		"header",
		nil,
		"Extra HTTP header (Name: Value)",
	)
	
	rootCmd.Flags().StringArrayVar(
	  &cfg.CustomRegexps,
	  "regexp",
	  nil,
	  "Extra regular expressions (Name: `pattern`)",
	)
	
	rootCmd.Flags().BoolVar(
	  &cfg.SkipDefaultRegexp,
	  "nodefrexp",
	  false,
	  "Omit loading of default regular expressions",
	)

	if err := rootCmd.MarkFlagRequired("url"); err != nil {
		panic(err)
	}
}
