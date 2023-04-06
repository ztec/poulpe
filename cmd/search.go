package cmd

import (
	"fmt"
	"git2.riper.fr/ztec/poulpe/dataloader"
	engine "git2.riper.fr/ztec/poulpe/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var search = &cobra.Command{
	Use:   "search",
	Short: "Search for an emoji from text",
	Long:  `Search for an emoji from text`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
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
			e, err := engine.NewFileBleveEngineFromEmojiList(cachePath, list)
			if err != nil {
				logrus.WithError(err).Error("Could not start emoji search engine")
				return err
			}
			CurrentEngine = engine.Engine(&e)
		}

		results, err := CurrentEngine.Search(strings.Join(args, " "))

		if err != nil {
			logrus.WithError(err).Error("Could not search for emoji")
			return err
		}

		for _, emoji := range results {
			fmt.Printf(
				"%s\t%s\t\"%s\"\t%s\t%s\n",
				emoji.Emoji,
				strings.Join(emoji.Aliases, " "),
				emoji.Description,
				emoji.Category,
				strings.Join(emoji.Tags, ","),
			)
		}
		return nil
	},
}

var (
	cachePath string
)

func init() {
	search.Flags().StringVar(&cachePath, "cachePath", "/tmp/poulpe", "path to store emoji and index. Will be created if it does not exist")
}
