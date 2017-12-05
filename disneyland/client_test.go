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
	configPath := os.Getenv("DISNEYLAND_TESTS_CONFIG")
	content, err := ioutil.ReadFile(configPath)
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
	//first
	createdJob, err := c.CreateJob(ctx, &Job{Status: Job_PENDING})
	checkTestErr(err, t)

	readJob, err := c.GetJob(ctx, &RequestWithId{Id: createdJob.Id})
	checkTestErr(err, t)

	if !checkJobsEqual(createdJob, readJob) {
		t.Fail()
	}

	createdJob.Status = Job_PENDING
	createdJob.Metadata = "updated_test"
	createdJob.Output = "updated_test"
	createdJob.Kind = "docker"

	updatedJob, err := c.ModifyJob(ctx, createdJob)
	checkTestErr(err, t)

	if !checkJobsEqual(createdJob, updatedJob) {
		t.Fail()
	}
	//second
	createdJob, err = c.CreateJob(ctx, &Job{Kind: "docker"})
	checkTestErr(err, t)

	allJobs, err := c.ListJobs(ctx, &ListJobsRequest{HowMany: 0})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 2 {
		t.Fail()
	}

	allJobs, err = c.ListJobs(ctx, &ListJobsRequest{HowMany: 1})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 1 {
		t.Fail()
	}

	allJobs, err = c.ListJobs(ctx, &ListJobsRequest{Kind: "docker", HowMany:2})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 2 {
		t.Fail()
	}

	pulledJobs, err := c.PullPendingJobs(ctx, &ListJobsRequest{HowMany: 1})
	checkTestErr(err, t)

	if len(pulledJobs.Jobs) != 1 {
		t.Fail()
	}
	//third
	createdJob, err = c.CreateJob(ctx, &Job{Kind: "remove"})
	checkTestErr(err, t)

	_, err = c.DeleteJob(ctx, &RequestWithId{Id: createdJob.Id})
	checkTestErr(err, t)
}
