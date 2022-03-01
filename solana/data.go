package solana

import (
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/pkg/errors"
)

type data struct {
	mint                      types.Account
	metadataAddress           common.PublicKey
	receiver                  common.PublicKey
	receiverAssociatedAddress common.PublicKey
}

func (s *solana) getData(receiver common.PublicKey) (data, error) {
	mint := types.NewAccount()

	receiverTokenAccountKey, _, err := common.FindAssociatedTokenAddress(receiver, mint.PublicKey)
	if err != nil {
		return data{}, errors.Wrap(err, "error getting associated receiver account pub key")
	}

	metaPublicKey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		return data{}, errors.Wrap(err, "error getting token metadata pub key")
	}

	return data{
		mint:                      mint,
		metadataAddress:           metaPublicKey,
		receiver:                  receiver,
		receiverAssociatedAddress: receiverTokenAccountKey,
	}, nil
}
