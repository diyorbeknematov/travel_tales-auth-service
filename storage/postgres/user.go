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

func (repo *UserRepo) CreateUser(user *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var (
		userResp  pb.RegisterResponse
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
		Scan(&userResp.Id, &userResp.Username, &userResp.Email, &userResp.FullName, &createdAt)

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

func (repo *UserRepo) ResetPassword(resetPassword models.UpdatePassword) (error) {
	_, err := repo.DB.Exec(`
		UPDATE 
			users 
		SET 
			password = $1 
		WHERE 
			id = $2
	`, resetPassword.NewPassword,resetPassword.ID )	

	if err != nil {
		repo.Logger.Error("Error in reset password", slog.String("eror", err.Error()))
		return nil
	}

	return nil
}
