package main

import (
	"github.com/olegfomenko/stop_war_ukraine/internal/config"
)

func main() {
	config.ConfigureService()
	infura := config.Infura()
	solana := config.Solana()
	_, url := infura.InitMetadata()
	solana.MintToken(url)
}
