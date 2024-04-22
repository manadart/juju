package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/canonical/go-dqlite"
	"github.com/canonical/go-dqlite/client"
	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
)

var (
	dbPathFlag = flag.String("db-dir", "/var/lib/juju/dqlite", "Path to the database directory")
)

func main() {
	flag.Parse()

	if *dbPathFlag == "" {
		fmt.Println("Please specify a database path")
		os.Exit(1)
	}

	// Read the desired node definition.
	info, err := nodeInfo(*dbPathFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	servers := []dqlite.NodeInfo{info}

	// Make sure we can use the node cluster store before reconfiguring.
	store, err := nodeClusterStore(*dbPathFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Reconfigure to a cluster of one with the desired node.
	if err := dqlite.ReconfigureMembership(*dbPathFlag, servers); err != nil {
		fmt.Println(errors.Annotate(err, "reconfiguring membership"))
		os.Exit(1)
	}

	// We successfully reconfigured Dqlite, so update the cluster YAML.
	if err := store.Set(context.Background(), servers); err != nil {
		fmt.Println(errors.Annotate(err, "setting new cluster details"))
		os.Exit(1)
	}
}

// nodeInfo reads the desired local node info from
// info.yaml in Juju's Dqlite data directory.
func nodeInfo(dataDir string) (dqlite.NodeInfo, error) {
	var node dqlite.NodeInfo

	data, err := os.ReadFile(path.Join(dataDir, "info.yaml"))
	if err != nil {
		return node, errors.Annotate(err, "reading info.yaml")
	}

	err = yaml.Unmarshal(data, &node)
	return node, errors.Annotate(err, "decoding NodeInfo")
}

// nodeClusterStore starts a store for reading/writing Dqlite cluster info.
// it is backed by cluster.yaml in Juju's Dqlite data directory.
func nodeClusterStore(dataDir string) (*client.YamlNodeStore, error) {
	store, err := client.NewYamlNodeStore(path.Join(dataDir, "cluster.yaml"))
	return store, errors.Annotate(err, "opening Dqlite cluster node store")
}
