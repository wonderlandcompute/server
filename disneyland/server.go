package disneyland

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Server struct {
	Storage   *DisneylandStorage
	SecretKey []byte
}

func detailedInternalError(err error) error {
	return grpc.Errorf(
		codes.Internal,
		fmt.Sprintf("Error creating job: %v", err),
	)
}

func (s *Server) Init() {
}

func (s *Server) CreateJob(ctx context.Context, in *Job) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	if in.Project == "" {
		in.Project = user.Project
	}
	if user.Project != in.Project && user.Project_access != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.Project")
	}
	createdJob, err := s.Storage.CreateJob(in, user)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return createdJob, nil
}
func (s *Server) CreateMultipleJobs(ctx context.Context, in *ListOfJobs) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)

	jobs_arr, err := s.Storage.CreateMultipleJobs(in.Jobs, user)
	if err != nil {
		return nil, detailedInternalError(err)
	}
	jobsList := &ListOfJobs{Jobs: jobs_arr}
	return jobsList, nil
}

func (s *Server) GetJob(ctx context.Context, in *RequestWithId) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	project := user.Project
	if user.Project_access == "ANY" {
		project = ""
	}
	job, err := s.Storage.GetJob(in.Id, project)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return job, nil
}

func (s *Server) ListJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	ret, err := s.Storage.ListJobs(user.Project)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) ModifyJob(ctx context.Context, in *Job) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	if user.Project != in.Project && user.Project_access != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.Project")
	}

	ret, err := s.Storage.UpdateJob(in)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) PullPendingJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	pts, err := s.Storage.PullJobs(in.HowMany)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return &ListOfJobs{Jobs: pts}, nil
}
