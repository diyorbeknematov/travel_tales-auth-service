package service

import (
	pb "auth-service/generated/user"
	"auth-service/storage/postgres"
	"auth-service/storage/redis"
	"context"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedAuthServiceServer
	UserRepo    *postgres.UserRepo
	Logger      *slog.Logger
	RedisClient *redis.RedisClient
}

func (s *UserService) UserInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	resp, err := s.UserRepo.GetUserInfo(in.Id)
	if err != nil {
		s.Logger.Error("Userni ma'lumotlarini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	resp, err := s.UserRepo.GetUserProfile(in.Id)
	if err != nil {
		s.Logger.Error("Xatolik get user pofilda serviceda", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	resp, err := s.UserRepo.UpdateUserProfile(in)
	if err != nil {
		s.Logger.Error("Error userni yangilashda", slog.String("error", err.Error()))
		return nil, err
	}
	return resp, nil
}

func (s *UserService) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	resp, err := s.UserRepo.GetUsers(in)
	if err != nil {
		s.Logger.Error("Error userlar ro'yxatini olishda", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	resp, err := s.UserRepo.DeleteUser(in.Id)
	if err != nil {
		s.Logger.Error("Userni o'chirishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}

func (s *UserService) FollowUser(ctx context.Context, in *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	resp, err := s.UserRepo.FollowingUser(in)
	if err != nil {
		s.Logger.Error("Userga follower bo'lishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	return resp, nil
}

func (s *UserService) ListFollowers(ctx context.Context, in *pb.ListFollowersRequest) (*pb.ListFollowersResponse, error) {
	resp, err := s.UserRepo.GetFollowers(in)
	if err != nil {
		s.Logger.Error("userni followerlarini olishda xatolik", slog.String("error", err.Error()))
		return nil, err
	}
	return resp, nil
}

func (s *UserService) GetUserActivity(ctx context.Context, in *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error) {
	resp, err := s.UserRepo.GetUserActivity(in.Id)
	if err != nil {
		s.Logger.Error("Error in get user activity", slog.String("error", err.Error()))
		return nil, err
	}

	return resp, nil
}
