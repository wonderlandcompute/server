package optimus

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type User struct {
	Username string
	Project  string
}

func getAuthUserFromContext(ctx context.Context) User {
	v := ctx.Value("authorized-user")
	if v != nil {
		return v.(User)
	}
	return User{}
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	allowed_endpoints := map[string]bool{
	// "/OptimusOptimus/CreatePoint":         true,
	}
	if allow, ok := allowed_endpoints[fullMethodName]; allow && ok {
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Error validating metadata")
	}
	token, ok := md["token"]
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Auth token not provided")
	}

	parsedToken, err := jwt.Parse(token[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "Unexpected algorithm!")
		}
		return s.SecretKey, nil
	})

	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid auth token")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		user := User{
			Username: claims["sub"].(string),
			Project:  claims["optimus_project"].(string),
		}
		return context.WithValue(ctx, "authorized-user", user), nil
	} else {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid auth token")
	}
}
