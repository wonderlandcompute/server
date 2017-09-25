package optimus

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

const (
	address = "localhost:50051"
	token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjE5ODc0ODYsIm5iZiI6MTUwNjM0OTA4Niwib3B0aW11c19wcm9qZWN0Ijoic2hpcC1zaGllbGQiLCJzdWIiOiJzYXNoYSJ9.4D7N3sMDLKm-mw6LPG7C1FZKOUbyGIbqm7Ic5I5BYqo"
)

func TestGRPCPointCRUD(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	checkTestErr(err, t)
	defer conn.Close()
	c := NewOptimusClient(conn)

	ctx := context.Background()
	md := metadata.Pairs("token", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	created_point, err := c.CreatePoint(ctx, &Point{})
	checkTestErr(err, t)

	t.Log(created_point)
}
