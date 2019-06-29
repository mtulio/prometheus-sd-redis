package rsd

import (
	"fmt"

	"github.com/apex/log"
	"github.com/mtulio/prometheus-sd-redis/src/rsd"
)

func ExampleNewClusterExporter() {
	// Redis node URL (one Redis node ingressed on the cluster)
	exp, err := rsd.NewClusterExporter(
		"my-redis-cluster-url.example.com",
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
		fmt.Printf("Redis cluster's node IP: %s\n", n.Addr)
	}
}

func ExampleNewClusterExporter_second() {
	fmt.Println("Second sample")
}
