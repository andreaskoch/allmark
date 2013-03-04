package indexer

import (
	"errors"
	"strings"
)

type Block struct {
	Name  string
	Value string
}

func EmptyBlock() Block {
	return Block{}
}

func NewBlock(name string, value string) (Block, error) {

	if strings.TrimSpace(name) == "" {
		return EmptyBlock(), errors.New("Cannot create a block without a name")
	}

	return Block{
		Name:  name,
		Value: value,
	}, nil
}
