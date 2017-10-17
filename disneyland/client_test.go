package disneyland

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

type DisneylandTestsConfig struct {
	ClientCert  string `yaml:"client_cert"`
	ClientKey   string `yaml:"client_key"`
	CACert      string `yaml:"ca_cert"`
	ConnectTo   string `yaml:"connect_to"`
	DatabaseURI string `yaml:"db_uri"`
}

var TestsConfig *DisneylandTestsConfig

func initTestsConfig() {
	TestsConfig = &DisneylandTestsConfig{}
	config_path := os.Getenv("DISNEYLAND_TESTS_CONFIG")
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

func checkJobsEqual(a *Job, b *Job) bool {
	return (a.Project == b.Project) &&
		(a.Id == b.Id) &&
		(a.Status == b.Status) &&
		(a.Metadata == b.Metadata) &&
		(a.Kind == b.Kind) &&
		(a.Output == b.Output) &&
		(a.Input == b.Input)
}

func TestGRPCJobCRUD(t *testing.T) {
	initTestsConfig()
	tc, err := getTransportCredentials()
	if err != nil {
		t.Fail()
	}

	conn, err := grpc.Dial(TestsConfig.ConnectTo, grpc.WithTransportCredentials(*tc))
	checkTestErr(err, t)
	defer conn.Close()
	c := NewDisneylandClient(conn)

	ctx := context.Background()

	created_job, err := c.CreateJob(ctx, &Job{Status: Job_PENDING})
	checkTestErr(err, t)

	read_job, err := c.GetJob(ctx, &RequestWithId{Id: created_job.Id})
	checkTestErr(err, t)

	if !checkJobsEqual(created_job, read_job) {
		t.Fail()
	}

	created_job.Status = Job_PENDING
	created_job.Metadata = "meta_test"
	created_job.Input = "input_test"
	created_job.Output = "output_test"
	created_job.Kind = "docker"

	updated_job, err := c.ModifyJob(ctx, created_job)
	checkTestErr(err, t)

	if !checkJobsEqual(created_job, updated_job) {
		t.Fail()
	}

	created_job, err = c.CreateJob(ctx, &Job{Project: "abc"})
	checkTestErr(err, t)

	all_jobs, err := c.ListJobs(ctx, &ListJobsRequest{HowMany: 2})
	checkTestErr(err, t)

	if len(all_jobs.Jobs) < 1 {
		t.Fail()
	}

	pulled_jobs, err := c.PullPendingJobs(ctx, &ListJobsKindRequest{HowMany: 2, Kind: "docker"})
	checkTestErr(err, t)

	if len(pulled_jobs.Jobs) < 1 {
		t.Fail()
	}

}
