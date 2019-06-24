package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"

	flag "github.com/spf13/pflag"

	"github.com/mtulio/prometheus-sd-redis/src/rsd"
)

func main() {
	// fInURL := flag.String("in-url", os.Getenv("REDIS_SD_IN_URL"), "In config URL. Single discovery")
	fInFile := flag.String("in-file", os.Getenv("REDIS_SD_IN_FILE"), "In config file")
	fOutFile := flag.String("out-file", os.Getenv("REDIS_SD_OUT_FILE"), "Out config file")
	flag.Parse()

	var services rsd.Services
	var sdconfigs rsd.SDConfigs

	// if *fInURL != "" {
	// 	fmt.Println("Using %s", *fInURL)
	// }

	jFile, err := os.Open(*fInFile)
	if err != nil {
		log.Error(fmt.Sprint(err))
	}

	byteValue, err := ioutil.ReadAll(jFile)
	if err != nil {
		log.Error(fmt.Sprint(err))
	}

	json.Unmarshal(byteValue, &services)

	for _, v := range services.Services {
		exp, err := rsd.NewClusterExporter(
			v.URL,
		)
		if err != nil {
			log.Error(fmt.Sprint(err))
		}
		err = exp.DiscoveryCluster()
		if err != nil {
			log.Errorf("ERROR discovering the cluster.")
		}
		c := rsd.SDConfig{}
		for _, n := range exp.ClusterNodes {
			c.Targets = append(c.Targets, n.Addr)
		}
		c.Labels = v.Labels
		c.Labels["job"] = v.Job
		sdconfigs = append(sdconfigs, c)
	}

	file, err := json.MarshalIndent(sdconfigs, "", " ")
	if err != nil {
		log.Error(fmt.Sprint(err))
	}

	err = ioutil.WriteFile(*fOutFile, file, 0644)
	if err != nil {
		log.Error(fmt.Sprint(err))
	}
	log.Infof("Service discovered was saved on file: %s", *fOutFile)
}
