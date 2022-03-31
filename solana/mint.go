package solana

import (
	"context"
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/program/assotokenprog"
	"github.com/olegfomenko/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/olegfomenko/solana-go-sdk/program/sysprog"
	"github.com/olegfomenko/solana-go-sdk/program/tokenprog"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/pkg/errors"
)

type toCreate struct {
	Name string
	Do   func() (types.Instruction, error)
}

func (s *Solana) getMint(meta Metadata, config MintConfig) ([]byte, error) {
	blockHash, rent, err := s.getSystemData()
	if err != nil {
		return nil, err
	}

	creators, creatorInstructions := s.parseCreators(meta.Properties.Creators, config)

	var rawInstructions = []toCreate{
		{
			Name: "creating token account",
			Do: func() (types.Instruction, error) {
				return sysprog.CreateAccount(
					config.Admin.PublicKey,
					config.mint.PublicKey,
					common.TokenProgramID,
					rent,
					tokenprog.MintAccountSize,
				)
			},
		},

		{
			Name: "initializing token account",
			Do: func() (types.Instruction, error) {
				return tokenprog.InitializeMint(
					0,
					config.mint.PublicKey,
					config.Admin.PublicKey,
					common.PublicKey{},
				)
			},
		},

		{
			Name: "creating & initializing token metadata",
			Do: func() (types.Instruction, error) {
				return tokenmeta.CreateMetadataAccount(
					config.metadataAddress,
					config.mint.PublicKey,
					config.Admin.PublicKey,
					config.Admin.PublicKey,
					config.Admin.PublicKey,
					true,
					false,
					tokenmeta.Data{
						Name:                 meta.Name,
						Symbol:               meta.Symbol,
						Uri:                  config.Metadata,
						SellerFeeBasisPoints: uint16(meta.SellerFeeBasisPoints),
						Creators:             &creators,
					},
				)
			},
		},

		{
			Name: "creating admin account for holding token",
			Do: func() (types.Instruction, error) {
				return assotokenprog.CreateAssociatedTokenAccount(
					config.Admin.PublicKey,
					config.Receiver,
					config.mint.PublicKey,
				)
			},
		},

		{
			Name: "minting token to receiver account",
			Do: func() (types.Instruction, error) {
				return tokenprog.MintTo(
					config.mint.PublicKey,
					config.receiverAssociatedAddress,
					config.Admin.PublicKey,
					[]common.PublicKey{},
					uint64(1),
				)
			},
		},

		{
			Name: "updating authority",
			Do: func() (types.Instruction, error) {
				return tokenmeta.UpdateMetadataAccount(
					config.metadataAddress,
					config.Admin.PublicKey,
					nil,
					&common.PublicKey{},
					&config.PrimarySaleHappened,
				)
			},
		},
	}

	rawInstructions = append(rawInstructions, creatorInstructions...)

	var instructions []types.Instruction
	for _, create := range rawInstructions {
		instruction, err := create.Do()
		if err != nil {
			return nil, errors.Wrap(err, "error while "+create.Name)
		}

		instructions = append(instructions, instruction)
	}

	// -- Creating transaction --
	tx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions:    instructions,
		Signers:         append(config.Creators, config.Admin, config.mint),
		FeePayer:        config.Admin.PublicKey,
		RecentBlockHash: blockHash,
	})

	if err != nil {
		return nil, errors.Wrap(err, "error creating tx")
	}

	return tx, nil
}

func (s *Solana) getSystemData() (string, uint64, error) {
	blockHash, err := s.GetRecentBlockhash(context.Background())
	if err != nil {
		return "", 0, errors.Wrap(err, "error getting recent block hash")
	}

	rent, err := s.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		return "", 0, errors.Wrap(err, "error getting rent")
	}

	return blockHash.Blockhash, rent, nil
}

func (s *Solana) parseCreators(creators []Creator, config MintConfig) ([]tokenmeta.Creator, []toCreate) {
	tokenCreators := make([]tokenmeta.Creator, len(creators))
	instructions := make([]toCreate, 0, len(creators))

	for i, creator := range creators {
		verified := creator.Address == config.Admin.PublicKey.String()
		address := common.PublicKeyFromString(creator.Address)

		tokenCreators[i] = tokenmeta.Creator{
			Address:  address,
			Verified: verified,
			Share:    creator.Share,
		}

		if verified {
			instructions = append(instructions, toCreate{
				Name: "Signing creator " + creator.Address,
				Do: func() (types.Instruction, error) {
					return tokenmeta.SignMetadata(config.metadataAddress, address)
				},
			})
		}
	}

	return tokenCreators, instructions
}
