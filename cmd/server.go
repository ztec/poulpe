package cmd

import (
	"fmt"
	"git2.riper.fr/ztec/poulpe/dataloader"
	"git2.riper.fr/ztec/poulpe/engine"
	"git2.riper.fr/ztec/poulpe/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var server = &cobra.Command{
	Use:   "server",
	Short: "Start web server",
	Long:  `Start web server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var CurrentEngine engine.Engine
		if engine.IsBleveEngineExist(cachePath) {
			e, err := engine.OpenBleveEngine(cachePath)
			if err != nil {
				logrus.WithError(err).Errorf("Could not load engine from storage at %s", cachePath)
				return err
			}
			CurrentEngine = engine.Engine(&e)
		} else {
			list, err := dataloader.FetchEmojiFromGithub()
			if err != nil {
				logrus.WithError(err).Error("Could not fetch emoji list from source")
				return err
			}
			os.MkdirAll(cachePath, 0777)
			e, err := engine.NewBleveEngineFromEmojiList(cachePath, list)
			if err != nil {
				logrus.WithError(err).Error("Could not start emoji search engine")
				return err
			}
			CurrentEngine = engine.Engine(&e)
		}
		return web.StartServer(
			fmt.Sprintf("%s:%d", serverBindAddress, serverPort),
			CurrentEngine,
		)
	},
}

var (
	serverPort        int
	serverBindAddress string
)

func init() {
	search.Flags().IntVar(&serverPort, "port", 8080, "Port to listen to")
	search.Flags().StringVar(&serverBindAddress, "bind-address", "0.0.0.0", "Interface Ip to bind to. 0.0.0.0 for all interfaces")
}
