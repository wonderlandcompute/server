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

func checkPointsEqual(a *Point, b *Point) bool {
	return (a.Project == b.Project) &&
		(a.Id == b.Id) &&
		(a.Status == b.Status) &&
		(a.Coordinate == b.Coordinate) &&
		(a.MetricValue == b.MetricValue) &&
		(a.Metadata == b.Metadata)
}

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

	read_point, err := c.GetPoint(ctx, &RequestWithId{Id: created_point.Id})
	checkTestErr(err, t)

	if !checkPointsEqual(created_point, read_point) {
		t.Fail()
	}

	created_point.Status = Point_FAILED
	created_point.MetricValue = "metric_test"
	created_point.Metadata = "meta_test"

	updated_point, err := c.ModifyPoint(ctx, created_point)
	checkTestErr(err, t)

	if !checkPointsEqual(created_point, updated_point) {
		t.Fail()
	}

	created_point, err = c.CreatePoint(ctx, &Point{Coordinate: "second"})
	checkTestErr(err, t)

	all_points, err := c.ListPoints(ctx, &ListPointsRequest{})
	checkTestErr(err, t)

	if len(all_points.Points) != 2 {
		t.Fail()
	}
}
