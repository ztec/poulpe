package engine

import "git2.riper.fr/ztec/poulpe/types"

type Engine interface {
	Search(q string) ([]types.EmojiDescription, error)
}
