package main

import (
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// Connect to dgraph
// Read and par the cvs file
// For each line -> insert into dgraph
//
// https://docs.dgraph.io/master/clients/#go

func main() {
	dCli := newClient()

}

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func test(c *dgo.Dgraph) {
	txn := c.NewTxn()
	defer txn.Discard()

}
