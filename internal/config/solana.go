package config

import (
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/olegfomenko/stop_war_ukraine/internal/solana"
	"github.com/spf13/viper"
)

const (
	solanaAddress = "solana.addr"
)

func Solana() solana.Solana {
	return solana.NewSolana(client.NewClient(viper.GetString(solanaAddress)))
}
