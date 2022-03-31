# Solana NFT Minter

[![olegfomenko](https://circleci.com/gh/olegfomenko/solana-nft-minter.svg?style=shield)](https://circleci.com/gh/olegfomenko/solana-nft-minter)
[![Go Report Card](https://goreportcard.com/badge/github.com/olegfomenko/solana-nft-minter)](https://goreportcard.com/report/github.com/olegfomenko/solana-nft-minter)


## Description

Library for minting NFTs from specified image on receiver or admin account. NFT data (image and metadata json) can be
stored on the IPFS Infura service.

## Using Example

Explore [Test Module file](./main_test.go)

Infura part

```go
inf := infura.NewDefaultInfura("https://ipfs.infura.io:5001", "project-id", "project-secret")

// Saving image in Infura
cid, err := inf.AddInfuraImage("./test_image.jpeg")

imageUrl := inf.GetLinkIPFS(cid)

metadata := solana.Metadata{
// Set fields here
}

// Saving metadata in Infura
cid, err = inf.AddInfuraJSON(metadata)
```

Solana part

```go
var admin types.Account{}

sol := solana.Solana{client.NewClient("https://api.devnet.solana.com")}

// Also you can use MintTokenUntilSuccess
hash, err := sol.MintToken(metadata, solana.MintConfig{
    Receiver:            admin.PublicKey,
    Admin:               admin,
    Creators:            []types.Account{admin},
    Metadata:            inf.GetLinkIPFS(cid),
    PrimarySaleHappened: true,
})
```
