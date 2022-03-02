package solana

import (
	"context"
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/pkg/errors"
)

type MintConfig struct {
	Receiver            common.PublicKey
	Admin               types.Account
	Creators            []types.Account
	Metadata            string
	PrimarySaleHappened bool

	*data
}

type Solana interface {
	MintToken(metadata Metadata, config MintConfig) (string, error)
}

type solana struct {
	*client.Client
}

func NewSolana(cli *client.Client) Solana {
	return &solana{
		cli,
	}
}

func (s *solana) MintToken(metadata Metadata, config MintConfig) (string, error) {
	err := s.genData(&config)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint data")
	}

	tx, err := s.getMint(metadata, config)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint tx")
	}

	mintHash, err := s.SendRawTransaction(context.Background(), tx)
	if err != nil {
		return "", errors.Wrap(err, "error sending mint tx")
	}

	return mintHash, nil
}
