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

func (s *solana) getMint(data data, meta Metadata, metadata string, primarySaleHappened bool) ([]byte, error) {
	blockHash, rent, err := s.getSystemData()
	if err != nil {
		return nil, err
	}

	creators, creatorInstructions, creatorAccouns := s.parseCreators(meta.Properties.Creators, data)

	var rawInstructions = []toCreate{
		{
			Name: "creating token account",
			Do: func() (types.Instruction, error) {
				return sysprog.CreateAccount(
					s.admin.PublicKey,
					data.mint.PublicKey,
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
					data.mint.PublicKey,
					s.admin.PublicKey,
					common.PublicKey{},
				)
			},
		},

		{
			Name: "creating & initializing token metadata",
			Do: func() (types.Instruction, error) {
				return tokenmeta.CreateMetadataAccount(
					data.metadataAddress,
					data.mint.PublicKey,
					s.admin.PublicKey,
					s.admin.PublicKey,
					s.admin.PublicKey,
					true,
					false,
					tokenmeta.Data{
						Name:                 meta.Name,
						Symbol:               meta.Symbol,
						Uri:                  metadata,
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
					s.admin.PublicKey,
					data.receiver,
					data.mint.PublicKey,
				)
			},
		},

		{
			Name: "minting token to receiver account",
			Do: func() (types.Instruction, error) {
				return tokenprog.MintTo(
					data.mint.PublicKey,
					data.receiverAssociatedAddress,
					s.admin.PublicKey,
					[]common.PublicKey{},
					uint64(1),
				)
			},
		},

		{
			Name: "updating authority",
			Do: func() (types.Instruction, error) {
				return tokenmeta.UpdateMetadataAccount(
					data.metadataAddress,
					s.admin.PublicKey,
					nil,
					&common.PublicKey{},
					&primarySaleHappened,
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
		Signers:         append(creatorAccouns, s.admin, data.mint),
		FeePayer:        s.admin.PublicKey,
		RecentBlockHash: blockHash,
	})

	if err != nil {
		return nil, errors.Wrap(err, "error creating tx")
	}

	return tx, nil
}

func (s *solana) getSystemData() (string, uint64, error) {
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

func (s *solana) parseCreators(creators []Creator, data data) ([]tokenmeta.Creator, []toCreate, []types.Account) {
	tokenCreators := make([]tokenmeta.Creator, len(creators))
	instructions := make([]toCreate, 0, len(creators))
	accounts := make([]types.Account, 0, len(creators))

	for i, creator := range creators {
		verified := creator.Address == s.admin.PublicKey.String()
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
					return tokenmeta.SignMetadata(data.metadataAddress, address)
				},
			})

			accounts = append(accounts, s.creators[creator.Address])
		}
	}

	return tokenCreators, instructions, accounts
}
