package infura

import (
	"bytes"
	"encoding/json"
)

func (i *Infura) AddInfuraJSON(val interface{}) (string, error) {
	jsonMetadata, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	cid, err := i.infura.Add(bytes.NewReader(jsonMetadata))
	if err != nil {
		return "", err
	}

	return cid, nil
}

func (i *Infura) AddInfuraImage(path string) (string, error) {
	cid, err := i.infura.AddDir(path)
	if err != nil {
		return "", err
	}

	return cid, nil
}
