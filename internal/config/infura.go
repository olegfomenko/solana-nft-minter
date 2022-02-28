package config

import (
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/olegfomenko/stop_war_ukraine/internal/metadata"
	"github.com/spf13/viper"
	"net/http"
)

const (
	infuraUrl    = "infura.url"
	infuraId     = "infura.project_id"
	infuraSecret = "infura.project_secret"
)

func Infura() metadata.Generator {
	return metadata.NewGenerator(ipfs.NewShellWithClient(
		viper.GetString(infuraUrl),
		NewClient(
			viper.GetString(infuraId),
			viper.GetString(infuraSecret),
		),
	), viper.GetString(infuraUrl),
	)
}

func NewClient(projectId, projectSecret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	}
}

type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}
