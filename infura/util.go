package infura

import (
	"fmt"
)

const IPFSLinkFormat = "%s/api/v0/cat/%s"

func (i *infura) GetLinkIPFS(cid string) string {
	return fmt.Sprintf(IPFSLinkFormat, i.infuraURL, cid)
}
