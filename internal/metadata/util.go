package metadata

import (
	"fmt"
)

const IPFSLinkFormat = "%s/api/v0/cat/%s"

func (g *generator) GetLinkIPFS(cid string) string {
	return fmt.Sprintf(IPFSLinkFormat, g.InfuraURL, cid)
}
