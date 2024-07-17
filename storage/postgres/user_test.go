package postgres

import (
	pb "auth-service/generated/user"
	"auth-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.CreateUser(models.RegisterRequest{
		Username: "sanjarbek",
		Email:    "sanjarbek2007@gmail.com",
		Password: "sqwerty007",
		FullName: "Sanjarbek",
	})

	if err != nil {
		t.Fatal(err)
	}

	response := &models.RegisterResponse{
		Username: "sanjarbek",
		Email:    "sanjarbek2007@gmail.com",
		FullName: "Sanjarbek",
	}

	assert.Equal(t, resp.Username, response.Username)
	assert.Equal(t, resp.Email, response.Email)
	assert.Equal(t, resp.FullName, response.FullName)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.NotEmpty(t, resp.ID)
}

func TestGetUserByEmail(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	getResp, err := userRepo.GetUserByEmail("diyorbeknematov@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	response := &models.LoginResponse{
		ID:       "975799c4-bd72-43c8-b0c5-93bd9461e033",
		Username: "diyorbek0321",
		Email:    "diyorbeknematov@gmail.com",
		Password: "$2a$10$dvdDaTR3ZBb5x5MMZf1xPOk3Nb5zbx35D8444y0OViFRsgQawEcu.",
	}

	assert.Equal(t, getResp, response)
}

func TestUpdatePassword(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.UpdatePassword(models.UpdatePassword{
		ID:          "9b0cf2c8-308c-4896-a737-511bff1bb991",
		NewPassword: "anvarnarzdeveloper",
	})

	if err != nil {
		t.Fatal(err)
	}

	response := &models.Success{
		Message: "Reset password successfully",
	}

	assert.Equal(t, resp, response)
}

func TestEmailExists(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.EmailExists("diyorbeknematov@gmail.com")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, resp, true)
}

func TestGetUserInfo(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.GetUserInfo("975799c4-bd72-43c8-b0c5-93bd9461e033")

	if err != nil {
		t.Fatal(err)
	}

	response := &pb.UserInfoResponse{
		Id:       "975799c4-bd72-43c8-b0c5-93bd9461e033",
		Username: "diyorbek0321",
		FullName: "string",
	}

	assert.Equal(t, resp, response)
}

func TestGetUserProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.GetUserProfile("975799c4-bd72-43c8-b0c5-93bd9461e033")

	if err != nil {
		t.Fatal(err)
	}

	response := &pb.GetProfileResponse{
		Id:               "975799c4-bd72-43c8-b0c5-93bd9461e033",
		Username:         "diyorbek0321",
		Email:            "diyorbeknematov@gmail.com",
		FullName:         "string",
		Bio:              "Go Backend Developer",
		CountriesVisited: 2,
		CreatedAt:        "2024-07-16 01:45:18",
		UpdatedAt:        "2024-07-16 01:45:18",
	}

	assert.Equal(t, resp, response)
}

func TestUpdateUserProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.UpdateUserProfile(&pb.UpdateProfileRequest{
		Id:               "975799c4-bd72-43c8-b0c5-93bd9461e033",
		FullName:         "Diyorbek Ne'matov",
		Bio:              "Go Backend Developer",
		CountriesVisited: 4,
	})

	if err != nil {
		t.Fatal(err)
	}

	response := &pb.UpdateProfileResponse{
		Id:               "975799c4-bd72-43c8-b0c5-93bd9461e033",
		Username:         "diyorbek0321",
		Email:            "diyorbeknematov@gmail.com",
		FullName:         "Diyorbek Ne'matov",
		Bio:              "Go Backend Developer",
		CountriesVisited: 4,
		UpdatedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}

	assert.Equal(t, resp, response)
}

func TestGetUsers(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.GetUsers(&pb.ListUsersRequest{
		Page:  1,
		Limit: 2,
	})

	if err != nil {
		t.Fatal(err)
	}

	expectedUsers := []*pb.User{
		{
			Id:               "975799c4-bd72-43c8-b0c5-93bd9461e033",
			Username:         "diyorbek0321",
			FullName:         "Diyorbek Ne'matov",
			CountriesVisited: 4,
		},
		{
			Id:               "9b0cf2c8-308c-4896-a737-511bff1bb991",
			Username:         "anvarnarz",
			FullName:         "Anvar Narzullayev",
			CountriesVisited: 0,
		},
	}

	expectedResp := &pb.ListUsersResponse{
		Users: expectedUsers,
		Total: 5,
		Page:  1,
		Limit: 2,
	}

	assert.Equal(t, resp, expectedResp)
}

func TestDeleteUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.DeleteUser("e1b9af75-931d-4d3b-acd7-a00e2571fa92")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, resp, &pb.DeleteUserResponse{
		Message: "User successfully deleted",
	})
}

func TestFollowingUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	resp, err := userRepo.FollowingUser(&pb.FollowUserRequest{
		FollowingId: "9b0cf2c8-308c-4896-a737-511bff1bb991",
		FollowerId:  "975799c4-bd72-43c8-b0c5-93bd9461e033",
	})

	if err != nil {
		t.Fatal(err)
	}

	response := &pb.FollowUserResponse{
		FollowingId: "9b0cf2c8-308c-4896-a737-511bff1bb991",
		FollowerId:  "975799c4-bd72-43c8-b0c5-93bd9461e033",
	}

	assert.Equal(t, resp.FollowerId, response.FollowerId)
	assert.Equal(t, resp.FollowingId, response.FollowingId)
}

func TestGetFollowers(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	// Call the method
	resp, err := repo.GetFollowers(&pb.ListFollowersRequest{
		UserId: "9b0cf2c8-308c-4896-a737-511bff1bb991",
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		t.Fatalf("Error getting followers: %v", err)
	}

	// Assertions
	expected := &pb.ListFollowersResponse{
		Followers: []*pb.Follower{
			{
				Id:       "975799c4-bd72-43c8-b0c5-93bd9461e033",
				Username: "diyorbek0321",
				FullName: "Diyorbek Ne'matov",
			},
		},
		Total: 5,
		Page:  1,
		Limit: 1,
	}

	assert.Equal(t, resp, expected)
}

func TestGetUserActivity(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)
	// Call the method
	resp, err := repo.GetUserActivity("975799c4-bd72-43c8-b0c5-93bd9461e033")
	if err != nil {
		t.Fatalf("Error getting user activity: %v", err)
	}

	// Assertions
	expected := &pb.GetUserActivityResponse{
		UserId:           "975799c4-bd72-43c8-b0c5-93bd9461e033",
		CountriesVisited: 4,
		LastActive:       "2024-07-17T10:45:29.803978+05:00",
	}

	assert.Equal(t, resp, expected)
}
