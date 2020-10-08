package file

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/ryarnyah/dblock/pkg/model"
	"github.com/ryarnyah/dblock/pkg/provider"
)

type fileProvider struct{}

var (
	fileSource = flag.String("file-source", ".new-schema.json", "New schema in a json file")
)

func init() {
	provider.RegisterProvider("file", fileProvider{})
}

func (s fileProvider) GetCurrentModel() (*model.DatabaseSchema, error) {
	b, err := ioutil.ReadFile(*fileSource)
	if err != nil {
		return nil, err
	}
	var newDatabase model.DatabaseSchema
	err = json.Unmarshal(b, &newDatabase)
	if err != nil {
		return nil, err
	}
	return &newDatabase, nil
}
