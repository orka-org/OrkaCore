package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/orka-org/orkacore/api/auth/v1"
	"github.com/orka-org/orkacore/internal/biz"
)

type AuthServiceService struct {
	pb.UnimplementedAuthServiceServer
	uc     biz.AuthUsecase
	logger log.Helper
}

func NewAuthServiceService(uc biz.AuthUsecase, logger log.Logger) *AuthServiceService {
	return &AuthServiceService{
		uc:     uc,
		logger: *log.NewHelper(logger),
	}
}

func (s *AuthServiceService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := s.uc.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{
		AccessToken:  res.Access,
		RefreshToken: res.Refresh,
	}, nil
}

func (s *AuthServiceService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := s.uc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		AccessToken:  res.Access,
		RefreshToken: res.Refresh,
	}, nil
}

func (s *AuthServiceService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	res, err := s.uc.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.InternalServer("no user found", "user not found")
	}

	resUser := &pb.User{}
	resUser.Id = res.ID
	resUser.Email = res.Email
	resUser.Username = res.Username
	if res.FirstName != nil {
		resUser.FirstName = res.FirstName
	}
	if res.LastName != nil {
		resUser.LastName = res.LastName
	}
	if res.Phone != nil {
		resUser.Phone = res.Phone
	}

	created := res.CreatedAt
	str := created.String()
	resUser.CreatedAt = &str
	updated := res.UpdatedAt
	str = updated.String()
	resUser.UpdatedAt = &str

	if res.DeletedAt != nil {
		str := res.DeletedAt.String()
		resUser.DeletedAt = &str
	}
	if res.SuspendedAt != nil {
		str := res.SuspendedAt.String()
		resUser.SuspendedAt = &str
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User:  resUser,
	}, nil
}

func (s *AuthServiceService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	res, err := s.uc.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.InternalServer("no user found", "user not found")
	}

	resUser := &pb.User{}
	resUser.Id = res.ID
	resUser.Email = res.Email
	resUser.Username = res.Username
	if res.FirstName != nil {
		resUser.FirstName = res.FirstName
	}
	if res.LastName != nil {
		resUser.LastName = res.LastName
	}
	if res.Phone != nil {
		resUser.Phone = res.Phone
	}

	created := res.CreatedAt
	str := created.String()
	resUser.CreatedAt = &str
	updated := res.UpdatedAt
	str = updated.String()
	resUser.UpdatedAt = &str

	if res.DeletedAt != nil {
		str := res.DeletedAt.String()
		resUser.DeletedAt = &str
	}
	if res.SuspendedAt != nil {
		str := res.SuspendedAt.String()
		resUser.SuspendedAt = &str
	}

	return &pb.GetUserResponse{
		User: resUser,
	}, nil
}
