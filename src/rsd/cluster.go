package rsd

import (
	"strings"
	"sync"

	"github.com/apex/log"
	"github.com/go-redis/redis"
)

// Exporter implements the prometheus.Exporter interface, and exports Redis metrics.
type Exporter struct {
	redisAddr string
}

// NewRedisExporter returns a new exporter of Redis metrics.
func NewRedisExporter(redisURI string) (*Exporter, error) {
	e := Exporter{
		redisAddr: redisURI,
	}

	return &e, nil
}

// Redis Cluster implementation

type ClusterExporter struct {
	sync.Mutex

	RedisEndpoint string
	ClusterNodes  []*ClusterNode
}

type ClusterNode struct {
	Exporter    *Exporter
	NodeID      string
	Addr        string
	Flags       string
	MasterID    string
	PingSent    string
	PongRecv    string
	ConfigEpoch string
	LinkState   string
	Slots       string
}

func NewClusterExporter(addr string) (*ClusterExporter, error) {

	ce := ClusterExporter{
		RedisEndpoint: addr,
	}

	return &ce, nil
}

// DiscoveryCluster discovery all nodes on the cluster.
func (ce *ClusterExporter) DiscoveryCluster() error {

	// erase all metrics before start the new discovery
	ce.cleanDiscovery()

	rcc := redis.NewClient(&redis.Options{
		Addr: ce.RedisEndpoint,
	})

	rccNodes, err := rcc.ClusterNodes().Result()
	if err != nil {
		log.Errorf("Unable to conn to redis endpoint %s . %s", ce.RedisEndpoint, err)
		return err
	}

	nodesCounter := 0.0
	for _, v := range strings.Split(rccNodes, "\n") {

		sArr := strings.Split(v, " ")
		if len(sArr) < 4 {
			continue
		}
		node := ClusterNode{
			NodeID:   sArr[0],
			Addr:     sArr[1],
			Flags:    sArr[2],
			MasterID: sArr[3],
		}

		node.PingSent = sArr[4]
		node.ConfigEpoch = sArr[5]
		node.PongRecv = sArr[6]
		node.LinkState = sArr[7]

		// Slave only
		if len(sArr) > 8 {
			node.Slots = sArr[8]
		}

		ne, err := NewRedisExporter(node.Addr)
		if err != nil {
			log.Errorf("ERROR creating Redis Exporter.")
		}
		node.Exporter = ne
		ce.ClusterNodes = append(ce.ClusterNodes, &node)

		nodesCounter++
	}

	return nil
}

func (ce *ClusterExporter) cleanDiscovery() {
	for i := len(ce.ClusterNodes) - 1; i >= 0; i-- {
		ce.ClusterNodes = removeSliceClusterNode(ce.ClusterNodes, i)
	}
}

func removeSliceClusterNode(s []*ClusterNode, i int) []*ClusterNode {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
