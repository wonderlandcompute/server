package wonderland

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	//"encoding/json"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"crypto/x509"
)

type Server struct {
	Storage   *WonderlandStorage
	SecretKey []byte
}

type Schedulers struct {
	Docker struct {
		SchedulerType []string `json:"scheduler_type"`
		Organisation  []struct {
			Name  string `json:"name"`
			Quota string `json:"quota"`
		} `json:"organisation"`
	} `json:"docker"`
	OptimisationPoint struct {
		SchedulerType []string    `json:"scheduler_type"`
		Organisation  interface{} `json:"organisation"`
	} `json:"optimisation-point"`
}

func detailedInternalError(err error) error {
	return grpc.Errorf(codes.Internal, fmt.Sprintf("Error processing job: %v", err))
}

func (s *Server) CreateJob(ctx context.Context, in *Job) (*Job, error) {
	user := getAuthUserFromContext(ctx)

	// if worker - Cannot create jobs
	if user.IsWorker() {
		return nil, grpc.Errorf(codes.PermissionDenied, "Workers cannot create jobs")
	}
	// if user - Can create jobs in their project
	in.Project = user.ProjectAccess

		createdJob, err := s.Storage.CreateJob(in, user)
		if err != nil {
			return nil, detailedInternalError(err)
		}

	return createdJob, nil
}

func (s *Server) GetJob(ctx context.Context, in *RequestWithId) (*Job, error) {
	user := getAuthUserFromContext(ctx)

	job, err := s.Storage.GetJob(in.Id)

	if err != nil {
		return nil, detailedInternalError(err)
	}
	// if user - Can get jobs from their project
	// if worker - Can get jobs with proper kind
	if !user.CanAccessJob(job) {
		return nil, grpc.Errorf(codes.PermissionDenied, "No access")
	}

	return job, nil
}

func (s *Server) ListJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)

	// if worker - Cannot list jobs
	if user.IsWorker() {
		return nil, grpc.Errorf(codes.PermissionDenied, "Workers cannot list jobs")
	}
	// if user - Can list jobs by kind in their project
	in.Project = user.ProjectAccess

	ret, err := s.Storage.ListJobs(in.HowMany, in.Project, in.Kind)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) ModifyJob(ctx context.Context, in *Job) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	// if user - Can modify jobs in their project
	// if worker - Can modify jobs with proper kind
	if !user.CanAccessJob(in) {
		return nil, grpc.Errorf(codes.PermissionDenied, "No access")
	}

	ret, err := s.Storage.UpdateJob(in)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) PullPendingJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	// if worker - Can pull jobs with proper kind
	if user.IsWorker() {
		in.Kind = user.KindAccess
	}
	// if user - Can pull jobs from their project
	if user.IsUser() {
		in.Project = user.ProjectAccess
	}

	pts, err := s.Storage.PullJobs(in.HowMany, in.Project, in.Kind)

	if err != nil {
		return nil, detailedInternalError(err)
	}

	return pts, nil
}

func (s *Server) DeleteJob(ctx context.Context, in *RequestWithId) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	// if worker - Cannot delete jobs
	if user.IsWorker() {
		return nil, grpc.Errorf(codes.PermissionDenied, "Workers cannot delete jobs")
	}
	// if user - Can delete jobs in their project
	ret, err := s.Storage.DeleteJob(in.Id, user.ProjectAccess)

	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

type WonderlandTestsConfig struct {
	ClientCert  string `yaml:"client_cert"`
	ClientKey   string `yaml:"client_key"`
	CACert      string `yaml:"ca_cert"`
	ConnectTo   string `yaml:"connect_to"`
	DatabaseURI string `yaml:"db_uri"`
}
var TestsConfig *WonderlandTestsConfig

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

func (s *Server) FetchAll(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	//TODO: get config file from os path
	configPath := "/Users/aleksandrsivcov/go/src/gitlab.com/disney/server/config/schedulers.json"
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	scheduler := Schedulers{}
	err = yaml.Unmarshal([]byte(content), &scheduler)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	user := getAuthUserFromContext(ctx)
	// if worker - Can pull jobs with proper kind
	if user.IsWorker() {
		in.Kind = user.KindAccess
	}
	// if user - Can pull jobs from their project
	if user.IsUser() {
		in.Project = user.ProjectAccess
	}

	pts, err := s.Storage.PullJobs(0, "", "docker")

	fmt.Println(pts)
	return nil, nil
}