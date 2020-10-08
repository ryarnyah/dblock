package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/glog"
	"github.com/ryarnyah/dblock/pkg/model"
	"github.com/ryarnyah/dblock/pkg/provider"
	"github.com/ryarnyah/dblock/pkg/rules"
	"github.com/ryarnyah/dblock/pkg/version"
)

var (
	providerName   = flag.String("provider", "postgres", "DB provider (supported values: postgres, file)")
	schemaLockFile = flag.String("database-lock-file", ".dblock.lock", "file where database schemas will be persisted")
	errorFile      = flag.String("error-json-file", "", "JSON file to write all errors")
	v              = flag.Bool("version", false, "Print version")
)

const (
	// BANNER for usage.
	BANNER = `________ __________.____                  __    
\______ \\______   \    |    ____   ____ |  | __
 |    |  \|    |  _/    |   /  _ \_/ ___\|  |/ /
 |    ` + "`" + `   \    |   \    |__(  <_> )  \___|    < 
/_______  /______  /_______ \____/ \___  >__|_ \
        \/       \/        \/          \/     \/
 Check db schema compatibility.
 Version: %s
 Build: %s

`
)

func persistModel(m *model.DatabaseSchema) {
	b, err := json.MarshalIndent(&m, "", "    ")
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(*schemaLockFile, b, 0644)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info("Database schemas has been writen successfully.")
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, BANNER, version.VERSION, version.GITCOMMIT)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *v {
		fmt.Printf(BANNER, version.VERSION, version.GITCOMMIT)
		return
	}

	p := provider.RegistredProviders[*providerName]
	currentModel, err := p.GetCurrentModel()
	if err != nil {
		glog.Fatal(err.Error())
	}

	_, err = os.Stat(*schemaLockFile)
	if os.IsNotExist(err) {
		persistModel(currentModel)
		return
	}
	if err != nil {
		glog.Fatal(err)
	}
	b, err := ioutil.ReadFile(*schemaLockFile)
	if err != nil {
		glog.Fatal(err)
	}
	var oldModel model.DatabaseSchema
	err = json.Unmarshal(b, &oldModel)
	if err != nil {
		glog.Fatal(err)
	}

	errors := make([]error, 0)
	for _, rule := range rules.RegistredRules {
		errors = append(errors, rule.CheckCompatibility(&oldModel, currentModel)...)
	}
	for _, err := range errors {
		glog.Info(err)
	}
	if len(errors) == 0 {
		persistModel(currentModel)
	}
	if *errorFile != "" {
		b, err := json.MarshalIndent(errors, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(*errorFile, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(len(errors))
}
