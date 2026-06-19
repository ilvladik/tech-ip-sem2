package student

import (
	"context"
	"strings"

	"example.com/pz2-grpc/gen/studentpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	studentpb.UnimplementedStudentServiceServer
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Ping(_ context.Context, req *studentpb.PingRequest) (*studentpb.PingResponse, error) {
	msg := req.GetMessage()
	if msg == "" {
		msg = "ping"
	}
	return &studentpb.PingResponse{
		Message: "Server received: " + msg,
	}, nil
}

func (s *Service) GetStudentByID(_ context.Context, req *studentpb.GetStudentRequest) (*studentpb.GetStudentResponse, error) {
	id := req.GetId()
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid student id")
	}

	st, err := s.repo.GetByID(id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "student not found")
	}

	return &studentpb.GetStudentResponse{Student: st}, nil
}

func (s *Service) ListStudents(_ context.Context, _ *studentpb.Empty) (*studentpb.ListStudentsResponse, error) {
	return &studentpb.ListStudentsResponse{
		Students: s.repo.ListAll(),
	}, nil
}

func (s *Service) CreateStudent(_ context.Context, req *studentpb.CreateStudentRequest) (*studentpb.GetStudentResponse, error) {
	fullName := strings.TrimSpace(req.GetFullName())
	group := strings.TrimSpace(req.GetGroup())
	email := strings.TrimSpace(req.GetEmail())
	specialization := strings.TrimSpace(req.GetSpecialization())

	if fullName == "" || group == "" || email == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name, group and email are required")
	}

	st, err := s.repo.Create(fullName, group, email, specialization)
	if err != nil {
		if err == ErrEmailExists {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &studentpb.GetStudentResponse{Student: st}, nil
}
