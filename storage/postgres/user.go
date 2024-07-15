package postgres

import (
	pb "auth-service/generated/user"
	"auth-service/models"
	"database/sql"
	"log/slog"
	"time"
)

type UserRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewUserRepo(db *sql.DB, logger *slog.Logger) *UserRepo {
	return &UserRepo{
		DB:     db,
		Logger: logger,
	}
}

func (repo *UserRepo) CreateUser(user models.Register) (*models.Register, error) {
	var (
		userResp  models.Register
		createdAt time.Time
	)
	err := repo.DB.QueryRow(`
		INSERT INTO users (
			username,
			email,
			password,
			full_name
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id,
			username,
			email,
			full_name,
			created_at
	`, user.Username, user.Email, user.Password, user.FullName).
		Scan(&userResp.ID, &userResp.Username, &userResp.Email, &userResp.FullName, &createdAt)

	if err != nil {
		repo.Logger.Error("Error creating user", slog.String("error", err.Error()))
		return nil, err
	}

	userResp.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	repo.Logger.Info("User created deleted")

	return &userResp, nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*models.UserLogin, error) {
	var userResp models.UserLogin

	err := repo.DB.QueryRow(`
		SELECT
			id,
			username,
			email,
			password
		FROM
			users
		WHERE
			deleted_at = 0 AND email = $1
	`, email).Scan(&userResp.ID, &userResp.Username, &userResp.Email)

	if err != nil {
		repo.Logger.Error("Error get user", slog.String("error", err.Error()))
		return nil, err
	}

	return &userResp, nil
}

func (repo *UserRepo) UpdatePassword(resetPassword models.UpdatePassword) (*models.Success, error) {
	_, err := repo.DB.Exec(`
		UPDATE 
			users 
		SET 
			password = $1 
		WHERE 
			id = $2 and deleted_at = 0
	`, resetPassword.NewPassword,resetPassword.ID )	

	if err != nil {
		repo.Logger.Error("Error in reset password", slog.String("eror", err.Error()))
		return &models.Success{
			Message: "Error in updated password",
		}, err
	}

	return &models.Success{
		Message: "Reset password successfully",
	},nil
}

func (repo *UserRepo) EmailExists(email string) (bool, error) {
	var exists bool
	err := repo.DB.QueryRow(`
		SELECT
			EXISTS (
				SELECT 1
				FROM users
				WHERE email = $1
			)
	`, email).Scan(&exists)

	if err != nil {
		repo.Logger.Error("Emailni bor yo'qligini tekshirishda xatolik", slog.String("error", err.Error()))
		return false, err
	}

	return exists, nil
}

func (repo *UserRepo) GetUserInfo(id string) (*pb.UserInfoResponse, error) {
	var info pb.UserInfoResponse

	err := repo.DB.QueryRow(`
		SELECT
			id,
			username,
			full_name
		FROM
			users
		WHERE
			deleted_at = 0 and id = $1
	`, id).Scan(&info.Id, &info.Username, &info.FullName)

	if err != nil {
		repo.Logger.Error("Error in get user info", slog.String("error", err.Error()))
		return nil, err
	}

	return &info, nil
}

func (repo *UserRepo) GetUserProfile(id string) (*pb.GetProfileResponse, error) {
	var (
		profile   pb.GetProfileResponse
		bio       sql.NullString
		createdAt time.Time
		updatedAt time.Time
	)

	err := repo.DB.QueryRow(`
		SELECT
			id,
			username,
			email,
			full_name,
			bio,
			countries_visited,
			created_at,
			updated_at
		FROM
			users
		WHERE
			id = $1 AND deleted_at = 0
	`, id).Scan(&profile.Id, &profile.Username, &profile.Email, &profile.FullName, &bio, &profile.CountriesVisited, &createdAt, &updatedAt)

	if err != nil {
		repo.Logger.Error("Error Get user profile", slog.String("error", err.Error()))
	}

	if !bio.Valid {
		profile.Bio = "No Bio"
	}

	profile.Bio = bio.String
	profile.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	profile.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	return &profile, nil
}

func (repo *UserRepo) UpdateUserProfile(req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	var (
		profile   pb.UpdateProfileResponse
		updatedAt time.Time
	)

	err := repo.DB.QueryRow(`
		UPDATE 
			users
		SET 
			full_name = $1,
			bio = $2,
			countries_visited = $3
		WHERE
			id = $1 AND deleted_at = 0
		RETURNING
			id,
			username,
			email,
			full_name,
			bio,
			countries_visited,
			updated_at
		
	`, req.FullName, req.Bio, req.CountriesVisited, req.Id).Scan(&profile.Id, &profile.Username, &profile.Email, &profile.FullName, &profile.Bio, &profile.CountriesVisited, &updatedAt)

	if err != nil {
		repo.Logger.Error("Error Get user profile", slog.String("error", err.Error()))
	}

	profile.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	return &profile, nil
}

func (repo *UserRepo) GetUsers(req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	offset := (req.Page - 1) * req.Limit

	rows, err := repo.DB.Query(`
		SELECT 
			id, 
			username, 
			full_name, 
			countries_visited
        FROM 
			users
        ORDER BY 
			username
        LIMIT $1 
		OFFSET $2
	`, req.Limit, offset)
	if err != nil {
		repo.Logger.Error("error executing query", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Username, &user.FullName, &user.CountriesVisited); err != nil {
			repo.Logger.Error("error scanning row", slog.String("error", err.Error()))
			return nil, err
		}
		users = append(users, &user)
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			users
		WHEEW 
			deleted_at = 0
	`).Scan(&total)

	if err != nil {
		repo.Logger.Error("error counting users", slog.String("error", err.Error()))
		return nil, err
	}

	resp := &pb.ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return resp, nil
}

func (repo *UserRepo) DeleteUser(id string) (*pb.DeleteUserResponse, error) {
	res, err := repo.DB.Exec(`
        UPDATE
            users
        SET
            deleted_at = $1
        WHERE
            deleted_at = 0 AND id = $2
    `, time.Now().Unix(), id)

	if err != nil {
		repo.Logger.Error("Error in user deletion", slog.String("error", err.Error()))
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		repo.Logger.Error("Error getting rows affected", slog.String("error", err.Error()))
		return nil, err
	}

	if rowsAffected == 0 {
		err := sql.ErrNoRows
		repo.Logger.Error("No user found to delete", slog.String("id", id))
		return nil, err
	}

	repo.Logger.Info("User successfully deleted", slog.String("id", id))

	return &pb.DeleteUserResponse{
		Message: "User successfully deleted",
	}, nil
}

func (repo *UserRepo) FollowingUser(req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	var follower pb.FollowUserResponse

	err := repo.DB.QueryRow(`
		INSERT INTO followers (
			follower_id,
			following_id
		)
		VALUES (
			$1,
			$2
		)
		RETURNING
			follower_id,
			following_id,
			followed_at
	`, req.FollowerId, req.FollowingId).Scan(&follower.FollowerId, &follower.FollowingId, &follower.FollowingAt)

	if err != nil {
		repo.Logger.Error("Error in following user", slog.String("error", err.Error()))
		return nil, err
	}
	return &follower, nil
}

func (repo *UserRepo) GetFollowers(req *pb.ListFollowersRequest) (*pb.ListFollowersResponse, error) {
	var followers []*pb.Follower
	offset := (req.Page - 1) * req.Limit
	rows, err := repo.DB.Query(`
		SELECT
			id,
			username,
			full_name
		FROM
			users u
		INNER JOIN
			followers f ON u.id = f.follower_id
		WHERE
			f.following_id = $1 and deleted_at = 0
		OFFSET $2
		LIMIT $3
	`, req.UserId, offset, req.Limit)

	if err != nil {
		repo.Logger.Error("Error in get followers", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var follower pb.Follower
		
		err = rows.Scan(&follower.Id, &follower.Username, &follower.FullName)
		if err != nil {
			repo.Logger.Error("Error in scan follower", slog.String("error", err.Error()))
			return nil, err
		}

		followers = append(followers, &follower)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT
			COUNT(*)
		FROM
			users u
		INNER JOIN
			followers f ON u.id = f.follower_id
		WHERE
			f.following_id = $1 and u.deleted_at = 0
	`).Scan(&total)
	if err != nil {
		repo.Logger.Error("Error in get followers count", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.ListFollowersResponse{
		Followers: followers,
		Total: total,
		Page: req.Page,
		Limit: req.Limit,
	}, nil
}