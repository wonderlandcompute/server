package wonderland

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

type WonderlandTestsConfig struct {
	ClientCert  string `yaml:"client_cert"`
	ClientKey   string `yaml:"client_key"`
	CACert      string `yaml:"ca_cert"`
	ConnectTo   string `yaml:"connect_to"`
	DatabaseURI string `yaml:"db_uri"`
}

var TestsConfig *WonderlandTestsConfig

func initTestsConfig() {
	TestsConfig = &WonderlandTestsConfig{}
	configPath := os.Getenv("WONDERLAND_TESTS_CONFIG")
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
	c := NewWonderlandClient(conn)

	ctx := context.Background()

	// first job
	createdJob, err := c.CreateJob(ctx, &Job{Kind: "docker"})
	checkTestErr(err, t)

	readJob, err := c.GetJob(ctx, &RequestWithId{Id: createdJob.Id})
	checkTestErr(err, t)

	if readJob.Status != Job_PENDING {
		t.Fail()
	}

	if !checkJobsEqual(createdJob, readJob) {
		t.Fail()
	}

	createdJob.Status = Job_PENDING
	createdJob.Metadata = "updated_test"
	createdJob.Output = "updated_test"

	updatedJob, err := c.ModifyJob(ctx, createdJob)
	checkTestErr(err, t)

	if !checkJobsEqual(createdJob, updatedJob) {
		t.Fail()
	}

	// second job
	createdJob, err = c.CreateJob(ctx, &Job{Kind: "docker"})
	checkTestErr(err, t)
	if createdJob.Status != Job_PENDING {
		t.Fail()
	}

	// listing
	allJobs, err := c.ListJobs(ctx, &ListJobsRequest{HowMany: 0})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 2 {
		t.Fail()
	}

	allJobs, err = c.ListJobs(ctx, &ListJobsRequest{HowMany: 0, Kind: "ock"})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 0 {
		t.Fail()
	}

	allJobs, err = c.ListJobs(ctx, &ListJobsRequest{HowMany: 1})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 1 {
		t.Fail()
	}

	allJobs, err = c.ListJobs(ctx, &ListJobsRequest{Kind: "docker", HowMany: 2})
	checkTestErr(err, t)

	if len(allJobs.Jobs) != 2 {
		t.Fail()
	}

	pulledJobs, err := c.PullPendingJobs(ctx, &ListJobsRequest{HowMany: 1})
	checkTestErr(err, t)

	if len(pulledJobs.Jobs) != 1 {
		t.Fail()
	}

	// third
	createdJob, err = c.CreateJob(ctx, &Job{Kind: "remove"})
	checkTestErr(err, t)

	_, err = c.DeleteJob(ctx, &RequestWithId{Id: createdJob.Id})
	checkTestErr(err, t)

	//kill job
	createdJob, err = c.CreateJob(ctx, &Job{Kind: "kill"})
	checkTestErr(err, t)

	readJob, err = c.KillJob(ctx, &RequestWithId{Id: createdJob.Id})
	checkTestErr(err, t)

	if readJob.Status != Job_KILLED {
		t.Fail()
	}
}
