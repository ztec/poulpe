package engine

import (
	"encoding/json"
	"fmt"
	"git2.riper.fr/ztec/poulpe/types"
	"github.com/AkinAD/emoji"
	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

const emojisFileName = "emojis.json"
const bleveIndexFileName = "bleve.index"

type BleveEngine struct {
	index  bleve.Index
	emojis []types.EmojiDescription
}

func IsBleveEngineExist(path string) bool {
	bleveIndexFileName := fmt.Sprintf("%s/%s", path, bleveIndexFileName)
	_, bleveIndexFileNameStat := os.Stat(bleveIndexFileName)
	emojiListFileName := fmt.Sprintf("%s/%s", path, emojisFileName)
	_, emojiListFileNameStat := os.Stat(emojiListFileName)

	if bleveIndexFileNameStat == nil && emojiListFileNameStat == nil {
		return true
	}
	return false
}

func OpenBleveEngine(path string) (engine BleveEngine, err error) {
	bleveIndexFileName := fmt.Sprintf("%s/%s", path, bleveIndexFileName)
	emojiListFileName := fmt.Sprintf("%s/%s", path, emojisFileName)
	index, err := bleve.Open(bleveIndexFileName)
	if err != nil {
		return
	}
	emojisData, err := os.ReadFile(emojiListFileName)
	if err != nil {
		return
	}
	var emojis []types.EmojiDescription
	err = json.Unmarshal(emojisData, &emojis)
	if err != nil {
		return
	}
	engine.index = index
	engine.emojis = emojis
	return
}

func NewFileBleveEngineFromEmojiList(path string, emojiList []types.EmojiDescription) (engine BleveEngine, err error) {
	logrus.Info("Creating Bleve index")
	bleveIndexFileName := fmt.Sprintf("%s/%s", path, bleveIndexFileName)
	emojiListFileName := fmt.Sprintf("%s/%s", path, emojisFileName)

	// we create a new indexMapping. I used the default one that will index all fields of my EmojiDescription
	mapping := bleve.NewIndexMapping()
	// we create the index instance
	bleveIndex, err := bleve.New(bleveIndexFileName, mapping)
	if err != nil {
		return
	}

	for eNumber, eDescription := range emojiList {
		// this will index each item one by one. No need to be quick here for me, I can wait few ms for the program to start.
		err = bleveIndex.Index(fmt.Sprintf("%d", eNumber), eDescription)
		if err != nil {
			logrus.WithField("emojiDescription", eDescription).WithError(err).Error("Could not index an emoji")
		}
	}
	data, err := json.Marshal(emojiList)
	if err != nil {
		return
	}
	err = os.WriteFile(emojiListFileName, data, 0777)
	if err != nil {
		return
	}
	engine.index = bleveIndex
	engine.emojis = emojiList
	return
}

func NewMemoryBleveEngineFromEmojiList(emojiList []types.EmojiDescription) (engine BleveEngine, err error) {
	// we create a new indexMapping. I used the default one that will index all fields of my EmojiDescription
	mapping := bleve.NewIndexMapping()
	// we create the index instance
	bleveIndex, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return
	}

	for eNumber, eDescription := range emojiList {
		// this will index each item one by one. No need to be quick here for me, I can wait few ms for the program to start.
		err = bleveIndex.Index(fmt.Sprintf("%d", eNumber), eDescription)
		if err != nil {
			logrus.WithField("emojiDescription", eDescription).WithError(err).Error("Could not index an emoji")
		}
	}
	engine.index = bleveIndex
	engine.emojis = emojiList
	return
}

func (b *BleveEngine) Search(q string) (results []types.EmojiDescription, err error) {
	if b.index == nil {
		// Pas de palais, pas de palais !
		return
	}

	qReverse, ok := emoji.FindReverse(removeSkinToneFromString(q))
	if ok {
		q = fmt.Sprintf("%s", strings.ReplaceAll(strings.ReplaceAll(qReverse, ":", ""), "_", " "))
		logrus.WithField("reverse", q).Info("Reverse search")
	}

	// we create a joinedQuery as bleve expect.
	queires := []query.Query{}
	// Querystring will interpret all fancy stuff
	queryString := bleve.NewQueryStringQuery(q)
	// Fuzzy will try to find terms near the one typed
	fuzzyQuery := bleve.NewFuzzyQuery(q)
	queires = append(queires, queryString, fuzzyQuery)

	if len(q) > 2 {
		// Will find anything that start by the joinedQuery
		prefixQuery := bleve.NewPrefixQuery(q)
		// will match anything containing the term
		wildcardQuery := bleve.NewWildcardQuery(fmt.Sprintf("*%s*", q))
		queires = append(queires, prefixQuery, wildcardQuery)
	}
	joinedQuery := bleve.NewDisjunctionQuery(queires...)

	// we define the search options and limit to 200 results. This should be enough.
	searchrequest := bleve.NewSearchRequestOptions(joinedQuery, 200, 0, false)
	// we do the search itself. This is the longest. Approximately few hundreds of us
	searchresults, err := b.index.Search(searchrequest)
	if err != nil {
		logrus.WithError(err).Error("Could not search for an emoji")
		return
	}

	// we return the results. I use the index to find my original object stored in `emojis` because it's simpler. Optimisation possible.
	for _, result := range searchresults.Hits {
		numIndex, _ := strconv.ParseInt(result.ID, 10, 64)
		results = append(results, b.emojis[numIndex])
	}
	return
}

func removeSkinToneFromString(q string) string {
	qRunes := []rune(q)
	finalRunes := []rune{}
	found := false
	for _, r := range qRunes {
		found = false
		for _, tone := range types.Tones {
			if r == tone[0] {
				found = true
			}
		}
		if !found {
			finalRunes = append(finalRunes, r)
		}
	}
	return string(finalRunes)
}
