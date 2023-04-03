package dataloader

import (
	"encoding/json"
	"fmt"
	"git2.riper.fr/ztec/poulpe/types"
	"github.com/go-zoox/fetch"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	GithubEmojiSourceUrl = "https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json"
)

func FetchEmojiFromGithub() (results []types.EmojiDescription, err error) {
	logrus.Info("Fetch emoji from github")
	response, err := fetch.Get(GithubEmojiSourceUrl)
	if err != nil {
		return
	}
	err = json.Unmarshal(response.Body, &results)
	logrus.Info("Add diversity to emoji list")
	results = enhanceEmojiListWithVariations(results)
	return
}

func enhanceEmojiListWithVariations(list []types.EmojiDescription) []types.EmojiDescription {
	for _, originalEmoji := range list {
		// we only add variations for emoji that supports it
		if originalEmoji.HasSkinTones {
			// we do it for every skin tone
			for skinToneName, tone := range types.Tones {
				// we make a copy of the emojiDescription
				currentEmojiWithSkinTone := originalEmoji

				// This is the important bit that took me hours to figure out
				// we convert the emoji in rune (string -> []rune). An emoji can already be composed of multiple sub UTF8 characters, therefore multiple runes.
				// we append to the list of runes the one for the skin tone.
				// finally, we convert that in string using the type conversion. Using fmt would result in printing all runes independently
				currentEmojiWithSkinTone.Emoji = string(append([]rune(currentEmojiWithSkinTone.Emoji), tone...))

				// we adapt the description and metadata to match the skin tone
				currentEmojiWithSkinTone.Description = fmt.Sprintf("%s %s", currentEmojiWithSkinTone.Description, skinToneName)
				aliases := []string{}
				for _, alias := range currentEmojiWithSkinTone.Aliases {
					// we update all aliases to include the skin tone
					aliases = append(aliases, fmt.Sprintf("%s_%s", alias, strings.ReplaceAll(strings.ReplaceAll(skinToneName, "-", "_"), " ", "_")))
				}
				currentEmojiWithSkinTone.Aliases = aliases
				// I cleared the unicode version because some emoji with skin tone were added way after their original. I could parse the unicode list,
				// but I'm a loafer, so I did not.
				currentEmojiWithSkinTone.UnicodeVersion = ""
				// we add the new emoji to the list
				list = append(list, currentEmojiWithSkinTone)
			}
		}
	}
	return list
}
