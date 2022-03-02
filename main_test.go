package main

import (
	"fmt"
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/olegfomenko/solana-nft-minter/infura"
	"github.com/olegfomenko/solana-nft-minter/solana"
	"testing"
)

func TestMintToAdmin(t *testing.T) {
	inf := infura.NewDefaultInfura("https://ipfs.infura.io:5001", "", "")
	sol := solana.NewSolana(client.NewClient("https://api.devnet.solana.com"))

	admin, _ := types.AccountFromBase58("")

	cid, err := inf.AddInfuraImage("./test_image.jpeg")
	if err != nil {
		panic(err)
	}

	metadata := solana.Metadata{
		Name:        "Test",
		Symbol:      "Test",
		Description: "Test",
		Properties: solana.Properties{
			Creators: []solana.Creator{
				solana.Creator{
					Address: admin.PublicKey.String(),
					Share:   100,
				},
			},
		},
		Image:                inf.GetLinkIPFS(cid),
		SellerFeeBasisPoints: 10,
	}

	cid, err = inf.AddInfuraJSON(metadata)
	if err != nil {
		panic(err)
	}

	hash, err := sol.MintToken(metadata, solana.MintConfig{
		Receiver:            admin.PublicKey,
		Admin:               admin,
		Creators:            []types.Account{admin},
		Metadata:            inf.GetLinkIPFS(cid),
		PrimarySaleHappened: true,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}

func TestMintToReceiver(t *testing.T) {
	inf := infura.NewDefaultInfura("https://ipfs.infura.io:5001", "", "")
	sol := solana.NewSolana(client.NewClient("https://api.devnet.solana.com"))

	admin, _ := types.AccountFromBase58("")
	receiver := common.PublicKeyFromString("")
	cid, err := inf.AddInfuraImage("./test_image.jpeg")
	if err != nil {
		panic(err)
	}

	metadata := solana.Metadata{
		Name:        "Test",
		Symbol:      "Test",
		Description: "Test",
		Properties: solana.Properties{
			Creators: []solana.Creator{
				solana.Creator{
					Address: admin.PublicKey.String(),
					Share:   100,
				},
			},
		},
		Image:                inf.GetLinkIPFS(cid),
		SellerFeeBasisPoints: 10,
	}

	cid, err = inf.AddInfuraJSON(metadata)
	if err != nil {
		panic(err)
	}

	hash, err := sol.MintToken(metadata, solana.MintConfig{
		Receiver:            receiver,
		Admin:               admin,
		Creators:            []types.Account{admin},
		Metadata:            inf.GetLinkIPFS(cid),
		PrimarySaleHappened: true,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}
