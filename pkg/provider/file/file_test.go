package file

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/ryarnyah/dblock/pkg/model"
)

func TestFileNotFound(t *testing.T) {
	fileProvider := fileProvider{}
	*fileSource = "non-existent-file"

	_, err := fileProvider.GetCurrentModel()
	if err == nil {
		t.Fatal("file must not exist")
	}
}

func TestDefunctFile(t *testing.T) {
	fileProvider := fileProvider{}
	f, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	*fileSource = f.Name()

	_, _ = f.WriteString("very bad file")
	_, err = fileProvider.GetCurrentModel()
	if err == nil {
		t.Fatal("must be unable to parse file")
	}
}

func TestLoadModelFromFile(t *testing.T) {
	fileProvider := fileProvider{}
	f, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	*fileSource = f.Name()

	expectedSchema := &model.DatabaseSchema{
		Schemas: []model.Schema{
			{
				SchemaName: "test",
				Tables: []model.TableSchema{
					{
						TableName: "test",
						Columns: []model.ColumnSchema{
							{
								ColunmName: "ID",
								ColumnType: "INTEGER",
							},
						},
					},
				},
			},
		},
	}
	b, err := json.Marshal(expectedSchema)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write(b)
	if err != nil {
		t.Fatal(err)
	}
	testSchema, err := fileProvider.GetCurrentModel()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testSchema, expectedSchema) {
		t.Fatalf("got %+v, expected %+v", testSchema, expectedSchema)
	}
}
