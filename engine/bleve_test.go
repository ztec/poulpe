package engine

import (
	"git2.riper.fr/ztec/poulpe/dataloader"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	engine *BleveEngine
)

func getEngine() BleveEngine {
	if engine == nil {
		cachePath := "/tmp/poulpe_tests"
		if IsBleveEngineExist(cachePath) {
			e, err := OpenBleveEngine(cachePath)
			if err != nil {
				logrus.WithError(err).Errorf("Could not load engine from storage at %s", cachePath)
				panic(err)
			}
			engine = &e
		} else {
			list, err := dataloader.FetchEmojiFromGithub()
			if err != nil {
				logrus.WithError(err).Error("Could not fetch emoji list from source")
				panic(err)
			}
			os.MkdirAll(cachePath, 0777)
			e, err := NewFileBleveEngineFromEmojiList(cachePath, list)
			if err != nil {
				logrus.WithError(err).Error("Could not start emoji search engine")
				panic(err)
			}
			engine = &e
		}
	}
	return *engine
}

func TestSearchFindAllHandsEmoji(t *testing.T) {
	be := getEngine()
	results, err := be.Search("hand")
	assert.NoError(t, err)
	assert.Equal(t, 135, len(results))
}

func TestSearchFindAllHappyEmoji(t *testing.T) {
	be := getEngine()
	results, err := be.Search("tags: happy")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(results))
}

func TestSearchFindCountryEmoji(t *testing.T) {
	be := getEngine()
	results, err := be.Search("france")
	assert.NoError(t, err)
	found := false
	for _, emojiDesc := range results {
		for _, alias := range emojiDesc.Aliases {
			if alias == "fr" {
				found = true
			}
		}
	}
	assert.True(t, found)
}

func TestSearchFindHugsEmojiWithHugEmoji(t *testing.T) {
	be := getEngine()
	results, err := be.Search("hug")
	assert.NoError(t, err)
	found := false
	for _, emojiDesc := range results {
		for _, alias := range emojiDesc.Aliases {
			if alias == "hugs" {
				found = true
			}
		}
	}
	assert.True(t, found)
}

func TestSearchFindDirectFromEmoji(t *testing.T) {
	be := getEngine()
	results, err := be.Search("🤭")
	assert.NoError(t, err)
	found := false
	for _, emojiDesc := range results {
		for _, alias := range emojiDesc.Aliases {
			if alias == "hand_over_mouth" {
				found = true
			}
		}
	}
	assert.True(t, found)
}
