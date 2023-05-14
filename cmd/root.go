/*
Copyright © 2023 Luiz Melo
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webcrawler",
	Short: "Given a starting URL, the crawler visits each URL it finds on the same domain",
	Long: `CLI for exploring a website structure (sitemap)

It prints each URL visited, and a list of links found on that page. The crawler is 
limited to one subdomain - so when you start with *https://somedomain.com/*, it will 
crawl all pages on the somedomain.com website, but not follow external links, 
for example to otherdomain.com or community.somedomain.com.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
