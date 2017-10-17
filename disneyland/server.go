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
		fmt.Sprintf("Error processing job: %v", err),
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

func (s *Server) GetJob(ctx context.Context, in *RequestWithId) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	job, err := s.Storage.GetJob(in.Id)
	if user.Project != job.Project && user.Project_access != "ANY" {
		err = grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.Project")
	}
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return job, nil
}

func (s *Server) ListJobs(ctx context.Context, in *ListJobsKindRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	user_kind := user.Kind_access
	kind := in.Kind
	if kind == "" {
		if user_kind == "ANY" {
			return nil, grpc.Errorf(codes.PermissionDenied, "ListJobsKindRequest.Kind not specified")
		}
		kind = user_kind
	}

	if user_kind != kind && user_kind != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Kind ≠ user.Kind_access")
	}
	ret, err := s.Storage.ListJobs(user.Project, kind)
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

func (s *Server) PullPendingJobs(ctx context.Context, in *ListJobsKindRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	kind := in.Kind
	user_kind := user.Kind_access
	if kind == "" {
		if user_kind == "ANY" {
			return nil, grpc.Errorf(codes.PermissionDenied, "ListJobsKindRequest.Kind not specified")
		}
		kind = user_kind
	}
	if user_kind != kind && user_kind != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Kind ≠ user.Kind_access")
	}

	pts, err := s.Storage.PullJobs(in.HowMany, kind)

	if err != nil {
		return nil, detailedInternalError(err)
	}

	return &ListOfJobs{Jobs: pts}, nil
}
