package metadata

import (
	"bytes"
	"encoding/json"
)

func (g *generator) AddInfuraJSON(val interface{}) (string, error) {
	jsonMetadata, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	cid, err := g.Infura.Add(bytes.NewReader(jsonMetadata))
	if err != nil {
		return "", err
	}

	return cid, nil
}

func (g *generator) AddInfuraImage(path string) (string, error) {
	cid, err := g.Infura.AddDir(path)
	if err != nil {
		return "", err
	}

	return cid, nil
}
