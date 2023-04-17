package engine

import "poulpe.ztec.fr/types"

type Engine interface {
	Search(q string) ([]types.EmojiDescription, error)
}
