package solana

import (
	"context"
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/pkg/errors"
)

type Solana interface {
	MintToken(metadata Metadata, meta string, receiver common.PublicKey, primarySaleHappened bool) (string, error)
}

func NewSolana(cli *client.Client, admin types.Account, creators map[string]types.Account) Solana {
	return &solana{
		cli,
		admin,
		creators,
	}
}

type solana struct {
	*client.Client
	admin    types.Account
	creators map[string]types.Account
}

func (s *solana) MintToken(metadata Metadata, meta string, receiver common.PublicKey, primarySaleHappened bool) (string, error) {
	data, err := s.getData(receiver)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint data")
	}

	tx, err := s.getMint(data, metadata, meta, primarySaleHappened)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint tx")
	}

	mintHash, err := s.SendRawTransaction(context.Background(), tx)
	if err != nil {
		return "", errors.Wrap(err, "error sending mint tx")
	}

	return mintHash, nil
}
