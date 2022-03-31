package infura

import (
	ipfs "github.com/ipfs/go-ipfs-api"
	"net/http"
)

type Infura struct {
	infura    *ipfs.Shell
	infuraURL string
}

func NewInfura(shell *ipfs.Shell, infuraURL string) *Infura {
	return &Infura{
		infura:    shell,
		infuraURL: infuraURL,
	}
}

type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func shellClient(projectId, projectSecret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	}
}

func NewDefaultInfura(url string, projectId, projectSecret string) *Infura {
	return NewInfura(ipfs.NewShellWithClient(
		url,
		shellClient(
			projectId,
			projectSecret,
		),
	), url)
}
