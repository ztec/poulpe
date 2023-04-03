package cmd

import (
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "poulpe search query",
	Short: "Emoji search engine",
	Long:  `Emoji search engine to find your best emoji from text on the cli or via a marvelous web interface`,
}

func init() {
	Root.AddCommand(search)
	Root.AddCommand(server)
}
