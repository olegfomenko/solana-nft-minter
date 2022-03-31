package solana

import (
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/olegfomenko/solana-go-sdk/types"
)

type data struct {
	mint                      types.Account
	metadataAddress           common.PublicKey
	receiverAssociatedAddress common.PublicKey
}

func (s *Solana) genData(config *MintConfig) (err error) {
	config.data = &data{}
	config.mint = types.NewAccount()

	config.data.receiverAssociatedAddress, _, err = common.FindAssociatedTokenAddress(config.Receiver, config.mint.PublicKey)
	if err != nil {
		return
	}

	config.data.metadataAddress, err = tokenmeta.GetTokenMetaPubkey(config.mint.PublicKey)
	if err != nil {
		return
	}

	return
}
