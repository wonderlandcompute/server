package optimus

import (
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type OptimusTestsConfig struct {
	ClientCert  string `yaml:"client_cert"`
	ClientKey   string `yaml:"client_key"`
	CACert      string `yaml:"ca_cert"`
	ConnectTo   string `yaml:"connect_to"`
	DatabaseURI string `yaml:"db_uri"`
}

var TestsConfig *OptimusTestsConfig

func initTestsConfig() {
	TestsConfig = &OptimusTestsConfig{}
	config_path := os.Getenv("OPTIMUS_TESTS_CONFIG")
	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), TestsConfig)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

}

func getTransportCredentials() (*credentials.TransportCredentials, error) {
	peerCert, err := tls.LoadX509KeyPair(TestsConfig.ClientCert, TestsConfig.ClientKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(TestsConfig.CACert)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tc := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{peerCert},
		RootCAs:      caCertPool,
	})

	return &tc, nil
}

func checkPointsEqual(a *Point, b *Point) bool {
	return (a.Project == b.Project) &&
		(a.Id == b.Id) &&
		(a.Status == b.Status) &&
		(a.Coordinate == b.Coordinate) &&
		(a.MetricValue == b.MetricValue) &&
		(a.Metadata == b.Metadata)
}

func TestGRPCPointCRUD(t *testing.T) {
	initTestsConfig()
	tc, err := getTransportCredentials()
	if err != nil {
		t.Fail()
	}

	conn, err := grpc.Dial(TestsConfig.ConnectTo, grpc.WithTransportCredentials(*tc))
	checkTestErr(err, t)
	defer conn.Close()
	c := NewOptimusClient(conn)

	ctx := context.Background()

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

	pulled_points, err := c.PullPendingPoints(ctx, &ListPointsRequest{HowMany: 2})
	checkTestErr(err, t)

	if len(pulled_points.Points) != 1 {
		t.Fail()
	}
}
