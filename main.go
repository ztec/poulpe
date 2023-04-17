package main

import (
	"github.com/spf13/cobra"
	"poulpe.ztec.fr/cmd"
)

func main() {
	cobra.CheckErr(cmd.Root.Execute())
}
