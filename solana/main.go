package solana

import (
	"context"
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/olegfomenko/solana-go-sdk/common"
	"github.com/olegfomenko/solana-go-sdk/types"
	"github.com/pkg/errors"
	"time"
)

const (
	DefaultRetries = 5
	DefaultDelay   = time.Second
)

type MintConfig struct {
	Receiver            common.PublicKey
	Admin               types.Account
	Creators            []types.Account
	Metadata            string
	PrimarySaleHappened bool

	*data
}

type Config struct {
	Retries int
	Delay   time.Duration
}

type Solana struct {
	*client.Client
}

func (s *Solana) MintToken(metadata Metadata, config MintConfig) (string, error) {
	err := s.genData(&config)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint data")
	}

	tx, err := s.getMint(metadata, config)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint tx")
	}

	mintHash, err := s.SendRawTransaction(context.Background(), tx)
	if err != nil {
		return "", errors.Wrap(err, "error sending mint tx")
	}

	return mintHash, nil
}

func (s *Solana) MintTokenUntilSuccess(metadata Metadata, mintConfig MintConfig, config *Config) (string, error) {
	if config == nil {
		config = &Config{
			Retries: DefaultRetries,
			Delay:   DefaultDelay,
		}
	}

	err := s.genData(&mintConfig)
	if err != nil {
		return "", errors.Wrap(err, "error generating mint data")
	}

	return s.mintTokenUntilSuccess(metadata, mintConfig, config.Retries, config.Delay)
}

func (s *Solana) mintTokenUntilSuccess(metadata Metadata, config MintConfig, retries int, delay time.Duration) (string, error) {
	for _i := 0; _i < retries; _i++ {
		tx, err := s.getMint(metadata, config)
		if err != nil {
			return "", errors.Wrap(err, "error generating mint tx")
		}

		mintHash, err := s.SendRawTransaction(context.Background(), tx)
		if err != nil {
			return "", errors.Wrap(err, "error sending mint tx")
		}

		time.Sleep(delay)

		ok, err := s.checkTxConfirmed(mintHash)
		if err != nil {
			return "", err
		}

		if ok {
			return mintHash, nil
		}
	}

	return "", errors.New("failed to send and check transaction: max retries proceed")
}

func (s *Solana) checkTxConfirmed(hash string) (bool, error) {
	statuses, err := s.GetSignatureStatuses(context.Background(), []string{hash})
	if err != nil {
		return false, errors.Wrap(err, "error checking tx confirm")
	}

	if len(statuses) < 1 {
		return false, errors.Wrap(err, "error checking tx confirm (array is empty)")
	}

	return statuses[0].Err == nil && (statuses[0].ConfirmationStatus != nil || statuses[0].Slot != 0), nil
}
