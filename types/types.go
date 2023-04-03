package types

var (
	Tones = map[string][]rune{
		"light skin tone":        []rune("\U0001F3FB"),
		"medium-light skin tone": []rune("\U0001F3FC"),
		"medium skin tone":       []rune("\U0001F3FD"),
		"medium-dark skin tone":  []rune("\U0001F3FE"),
		"dark skin tone":         []rune("\U0001F3FF"),
	}
)

type EmojiDescription struct {
	Emoji          string   `json:"emoji"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	HasSkinTones   bool     `json:"skin_tones,omitempty"`
	UnicodeVersion string   `json:"unicode_version"`
}
