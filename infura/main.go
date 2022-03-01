package infura

import (
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/olegfomenko/solana-nft-minter/solana"
	"github.com/pkg/errors"
)

type Infura interface {
	InitMetadata(metadata solana.Metadata, path string) (string, error)
}

type infura struct {
	infura    *ipfs.Shell
	infuraURL string
}

func NewGenerator(shell *ipfs.Shell, infuraURL string) Infura {
	return &infura{
		infura:    shell,
		infuraURL: infuraURL,
	}
}

func (i *infura) InitMetadata(metadata solana.Metadata, path string) (string, error) {
	cid, err := i.AddInfuraImage(path)
	if err != nil {
		return "", errors.Wrap(err, "error saving image to infura")
	}

	cid, err = i.AddInfuraJSON(metadata)
	if err != nil {
		return "", errors.Wrap(err, "error saving json to infura")
	}

	return fmt.Sprintf(IPFSLinkFormat, i.infuraURL, cid), nil
}
