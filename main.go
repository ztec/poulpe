package main

import (
	"git2.riper.fr/ztec/poulpe/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.Root.Execute())
}
