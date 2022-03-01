package main

import (
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-nft-minter/infura"
	"github.com/olegfomenko/solana-nft-minter/solana"
)

type Generator struct {
	infura.Infura
	solana.Solana
}

func (g *Generator) Generate(
	metadata solana.Metadata,
	pathToFile string,
	receiver common.PublicKey,
	primarySaleHappened bool,
) (url string, hash string, err error) {
	url, err = g.InitMetadata(metadata, pathToFile)
	if err != nil {
		return
	}

	hash, err = g.MintToken(metadata, url, receiver, primarySaleHappened)
	return
}
