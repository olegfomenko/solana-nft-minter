package infura

import (
	"fmt"
)

const IPFSLinkFormat = "https://gateway.ipfs.io/ipfs/%s"

func (i *Infura) GetLinkIPFS(cid string) string {
	return fmt.Sprintf(IPFSLinkFormat, cid)
}
