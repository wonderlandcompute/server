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
	userProject := user.ProjectAccess
	project := in.Project
	if project == "" {
		if userProject == "ANY" {
			return nil, grpc.Errorf(codes.DataLoss, "Job.Project not specified")
		}
		project = userProject
	}
	if userProject != project && userProject != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.ProjectAccess")
	}
	in.Project = project
	createdJob, err := s.Storage.CreateJob(in, user)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return createdJob, nil
}

func (s *Server) GetJob(ctx context.Context, in *RequestWithId) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	userProject := user.ProjectAccess
	job, err := s.Storage.GetJob(in.Id)
	if userProject != job.Project && userProject != "ANY" {
		err = grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.ProjectAccess")
	}
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return job, nil
}

func (s *Server) ListJobs(ctx context.Context, in *ListJobsRequest) (*ListOfJobs, error) {
	user := getAuthUserFromContext(ctx)
	userKind := user.KindAccess
	userProject := user.ProjectAccess
	kind := in.Kind
	project := in.Project
	if kind == "" {
		if userKind == "ANY" {
			return nil, grpc.Errorf(codes.DataLoss, "ListJobsKindRequest.Kind not specified")
		}
		kind = userKind
	}
	if project == "" {
		if userProject == "ANY" {
			return nil, grpc.Errorf(codes.DataLoss, "ListJobsKindRequest.Project not specified")
		}
		project = userProject
	}

	if userKind != kind && userKind != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Kind ≠ user.KindAccess")
	}
	if userProject != project && userProject != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.ProjectAccess")
	}
	ret, err := s.Storage.ListJobs(project, kind)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) ModifyJob(ctx context.Context, in *Job) (*Job, error) {
	user := getAuthUserFromContext(ctx)
	userProject := user.ProjectAccess
	project := in.Project
	if project == "" && userProject == "ANY" {
		return nil, grpc.Errorf(codes.DataLoss, "Job.Project not specified")
	}
	if userProject != project && userProject != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.ProjectAccess")
	}
	//manag
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
	userProject := user.ProjectAccess
	project := in.Project

	if kind == "" {
		if userKind == "ANY" {
			return nil, grpc.Errorf(codes.DataLoss, "ListJobsKindRequest.Kind not specified")
		}
		kind = userKind
	}
	if project == "" {
		if userProject == "ANY" {
			return nil, grpc.Errorf(codes.DataLoss, "ListJobsKindRequest.Project not specified")
		}
		project = userProject
	}

	if userKind != kind && userKind != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Kind ≠ user.KindAccess")
	}
	if userProject != project && userProject != "ANY" {
		return nil, grpc.Errorf(codes.PermissionDenied, "job.Project ≠ user.ProjectAccess")
	}

	pts, err := s.Storage.PullJobs(in.HowMany, kind)

	if err != nil {
		return nil, detailedInternalError(err)
	}

	return &ListOfJobs{Jobs: pts}, nil
}
