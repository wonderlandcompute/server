package optimus

import (
	"golang.org/x/net/context"
)

type Server struct {
	Storage   *OptimusStorage
	SecretKey []byte
}

func (s *Server) Init() {
}

func (s *Server) CreatePoint(ctx context.Context, in *Point) (*Point, error) {
	user := getAuthUserFromContext(ctx)
	return &Point{Coordinate: user.Username}, nil
}

func (s *Server) GetPoint(ctx context.Context, in *RequestWithId) (*Point, error) {
	return nil, nil
}

func (s *Server) ListPoints(ctx context.Context, in *ListPointsRequest) (*ListOfPoints, error) {
	return nil, nil
}

func (s *Server) ModifyPoint(ctx context.Context, in *Point) (*Point, error) {
	return nil, nil
}
