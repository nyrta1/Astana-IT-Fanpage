package repository

import (
	"database/sql"
)

type UserTypeRepo interface {
	CreateType(userType string) error
	GetTypeByID(id uint) (string, error)
	GetIDByTypeName(userType string) (uint, error)
	GetAllType() ([]string, error)
	UpdateTypeByID(id uint, userType string) error
	DeleteTypeByID(id uint) error
}

type UserTypeRepository struct {
	db *sql.DB
}

func NewUserTypeRepository(db *sql.DB) *UserTypeRepository {
	return &UserTypeRepository{db}
}

func (r *UserTypeRepository) CreateType(userType string) error {
	_, err := r.db.Exec(
		"INSERT INTO user_type(typename) VALUES ($1)",
		userType,
	)
	return err
}

func (r *UserTypeRepository) GetTypeByID(id uint) (string, error) {
	var typename string
	err := r.db.QueryRow("SELECT typename FROM user_type WHERE id = $1", id).Scan(
		&typename,
	)
	if err != nil {
		return "", err
	}
	return typename, err
}

func (r *UserTypeRepository) GetIDByTypeName(userType string) (uint, error) {
	var id uint
	err := r.db.QueryRow("SELECT id FROM user_type WHERE typename = $1", userType).Scan(
		&id,
	)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserTypeRepository) GetAllType() ([]string, error) {
	rows, err := r.db.Query("SELECT * FROM user_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTypes []string
	for rows.Next() {
		var userType string
		err := rows.Scan(
			userType)
		if err != nil {
			return nil, err
		}
		userTypes = append(userTypes, userType)
	}

	return userTypes, nil
}

func (r *UserTypeRepository) UpdateTypeByID(id uint, userType string) error {
	_, err := r.db.Exec("UPDATE user_type SET typename = $1 WHERE id = $2",
		userType, id)
	return err
}

func (r *UserTypeRepository) DeleteTypeByID(id uint) error {
	_, err := r.db.Exec("DELETE FROM user_type WHERE id = $1", id)
	return err
}
