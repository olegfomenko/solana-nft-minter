package metadata

import (
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/olegfomenko/stop_war_ukraine/internal/constants"
)

type Generator interface {
	InitMetadata() (Metadata, string)
}

type generator struct {
	Infura    *ipfs.Shell
	InfuraURL string
}

func NewGenerator(infura *ipfs.Shell, infuraURL string) Generator {
	return &generator{
		Infura:    infura,
		InfuraURL: infuraURL,
	}
}

func (g *generator) InitMetadata() (Metadata, string) {
	cid, err := g.AddInfuraImage(constants.ImagePath)
	if err != nil {
		panic(err.Error())
	}

	metadata := Metadata{
		Name:                 constants.ImageName,
		Symbol:               constants.TokenSymbol,
		Description:          constants.TokenDescription,
		Image:                fmt.Sprintf(IPFSLinkFormat, g.InfuraURL, cid),
		SellerFeeBasisPoints: constants.SellerFee,
		Properties: Properties{
			Files: []Files{
				{
					Uri:  fmt.Sprintf(IPFSLinkFormat, g.InfuraURL, cid),
					Type: constants.TokenImageType,
				},
			},
			Creators: []Creator{
				{
					Address: constants.Admin.PublicKey.String(),
					Share:   100,
				},
			},
		},
		Collection: Collection{
			Name:   constants.CollectionName,
			Family: constants.CollectionFamily,
		},
	}

	cid, err = g.AddInfuraJSON(metadata)
	if err != nil {
		panic(err)
	}

	return metadata, fmt.Sprintf(IPFSLinkFormat, g.InfuraURL, cid)
}
