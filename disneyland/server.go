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
