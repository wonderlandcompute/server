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
	if user.Project != in.Project && user.ProjectAccess != "ANY" {
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
	if user.Project != job.Project && user.ProjectAccess != "ANY" {
		err = grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.Project")
	}
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return job, nil
}

func (s *Server) ListJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	userKind := user.KindAccess
	kind := in.Kind
	if kind == "" {
		if userKind == "ANY" {
			return nil, grpc.Errorf(codes.PermissionDenied, "ListJobsKindRequest.Kind not specified")
		}
		kind = userKind
	}

	if userKind != kind && userKind != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Kind ≠ user.KindAccess")
	}
	ret, err := s.Storage.ListJobs(user.Project, kind)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) ModifyJob(ctx context.Context, in *Job) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	if user.Project != in.Project && user.ProjectAccess != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.Project")
	}

	ret, err := s.Storage.UpdateJob(in)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) PullPendingJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	kind := in.Kind
	userKind := user.KindAccess
	if kind == "" {
		if userKind == "ANY" {
			return nil, grpc.Errorf(codes.PermissionDenied, "ListJobsKindRequest.Kind not specified")
		}
		kind = userKind
	}
	if userKind != kind && userKind != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Kind ≠ user.KindAccess")
	}

	pts, err := s.Storage.PullJobs(in.HowMany, kind)

	if err != nil {
		return nil, detailedInternalError(err)
	}

	return &ListOfJobs{Jobs: pts}, nil
}
