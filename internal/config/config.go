package config

import (
	"context"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

const DefaultDgraphURI = "localhost:9080"

func MustGetDgraphConn(ctx context.Context, uri string) (*dgo.Dgraph, func() error) {
	d, err := grpc.DialContext(ctx, uri, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("failed to conn" + err.Error())
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	), d.Close
}
