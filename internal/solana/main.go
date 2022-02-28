package solana

import (
	"context"
	"fmt"
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/program/assotokenprog"
	"github.com/olegfomenko/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/olegfomenko/solana-go-sdk/program/sysprog"
	"github.com/olegfomenko/solana-go-sdk/program/tokenprog"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/olegfomenko/stop_war_ukraine/internal/constants"
	"github.com/pkg/errors"
)

type Solana interface {
	MintToken(meta string)
}

func NewSolana(cli *client.Client) Solana {
	return &solana{cli}
}

type solana struct {
	*client.Client
}

func (s *solana) MintToken(meta string) {
	data, err := s.getMintData()
	if err != nil {
		panic(err)
	}

	tx, err := s.getMint(data, meta)
	if err != nil {
		panic(err)
	}

	mintHash, err := s.SendRawTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Minting tx hash: ", mintHash)
}

type MintTokenData struct {
	Mint                   types.Account
	AdminAssociatedAddress common.PublicKey
	MetadataAddress        common.PublicKey
}

func (s *solana) getMintData() (MintTokenData, error) {
	mint := types.NewAccount()

	adminTokenAccountKey, _, err := common.FindAssociatedTokenAddress(constants.Admin.PublicKey, mint.PublicKey)
	if err != nil {
		return MintTokenData{}, errors.Wrap(err, "error getting associated admin account pub key")
	}

	metaPublicKey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		return MintTokenData{}, errors.Wrap(err, "error getting token metadata pub key")
	}

	return MintTokenData{
		Mint:                   mint,
		AdminAssociatedAddress: adminTokenAccountKey,
		MetadataAddress:        metaPublicKey,
	}, nil
}

func (s *solana) getMint(data MintTokenData, metadata string) ([]byte, error) {
	blockHash, err := s.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error getting recent block hash")
	}

	rent, err := s.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		return nil, errors.Wrap(err, "error getting rent")
	}

	type ToCreate struct {
		Name string
		Do   func() (types.Instruction, error)
	}

	var toCreate = []ToCreate{
		{
			Name: "creating token account",
			Do: func() (types.Instruction, error) {
				return sysprog.CreateAccount(
					constants.Admin.PublicKey,
					data.Mint.PublicKey,
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
					data.Mint.PublicKey,
					constants.Admin.PublicKey,
					common.PublicKey{},
				)
			},
		},

		{
			Name: "creating & initializing token metadata",
			Do: func() (types.Instruction, error) {
				return tokenmeta.CreateMetadataAccount(
					data.MetadataAddress,
					data.Mint.PublicKey,
					constants.Admin.PublicKey,
					constants.Admin.PublicKey,
					constants.Admin.PublicKey,
					true,
					false,
					tokenmeta.Data{
						Name:                 constants.ImageName,
						Symbol:               constants.TokenSymbol,
						Uri:                  metadata,
						SellerFeeBasisPoints: constants.SellerFee,
						Creators: &[]tokenmeta.Creator{
							{
								Address:  constants.Admin.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
					},
				)
			},
		},

		{
			Name: "creating admin account for holding token",
			Do: func() (types.Instruction, error) {
				return assotokenprog.CreateAssociatedTokenAccount(
					constants.Admin.PublicKey,
					constants.Admin.PublicKey,
					data.Mint.PublicKey,
				)
			},
		},

		{
			Name: "minting token to admin's account",
			Do: func() (types.Instruction, error) {
				return tokenprog.MintTo(
					data.Mint.PublicKey,
					data.AdminAssociatedAddress,
					constants.Admin.PublicKey,
					[]common.PublicKey{},
					uint64(1),
				)
			},
		},

		{
			Name: "updating authority and primary sale happened",
			Do: func() (types.Instruction, error) {
				return tokenmeta.UpdateMetadataAccount(
					data.MetadataAddress,
					constants.Admin.PublicKey,
					nil,
					&common.PublicKey{},
					nil,
				)
			},
		},
	}

	var instructions []types.Instruction
	for _, create := range toCreate {
		instruction, err := create.Do()
		if err != nil {
			return nil, errors.Wrap(err, "error while "+create.Name)
		}

		instructions = append(instructions, instruction)
	}

	// -- Creating transaction --
	tx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions:    instructions,
		Signers:         []types.Account{constants.Admin, data.Mint},
		FeePayer:        constants.Admin.PublicKey,
		RecentBlockHash: blockHash.Blockhash,
	})

	if err != nil {
		return nil, errors.Wrap(err, "error creating tx")
	}

	return tx, nil
}
