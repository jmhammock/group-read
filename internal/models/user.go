package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	UserTypeId string `json:"user_type_id"`
}

type Users []User

type UserModel struct {
	db DB
}

func NewUserModel(db DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (u UserModel) Get(limit, offset uint32) (*Users, error) {
	q := `SELECT *
		FROM users
		LIMIT $1
		OFFSET $2;`

	rows, err := u.db.Query(q, limit, offset)
	if err != nil {
		return nil, err
	}

	var users Users
	for rows.Next() {
		var user User
		err := rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (u UserModel) GetById(id string) (*User, error) {
	q := `SELECT *
		FROM users
		WHERE id = $1;`

	var user *User
	err := u.db.QueryRow(q, id).Scan(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserModel) Insert(user *User) (string, error) {
	q := `INSERT INTO users (first_name, last_name, email, password, user_type_id)
		VALUES($1, $2, $3, $3, $4, $5)
		RETURNING id;`

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	var id string
	args := []interface{}{
		user.FirstName,
		user.LastName,
		user.Email,
		string(password),
		user.UserTypeId,
	}
	err = u.db.QueryRow(q, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (u UserModel) Update(user *User) error {
	q := `UPDATE users
		SET first_name = $2,
			last_name = $3,
			email = $4,
			user_type_id = $5
		WHERE id = $1;`

	args := []interface{}{
		user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.UserTypeId,
	}
	_, err := u.db.Execute(q, args...)
	return err
}

func (u UserModel) UpdatePassword(id, pw string) error {
	q := `UPDATE users
		SET password = $2
		WHERE id = $1;`

	password, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = u.db.Execute(q, string(password))
	return err
}

func (u UserModel) Delete(id string) error {
	q := `DELETE FROM users
		WHERE id = $1;`

	_, err := u.db.Execute(q, id)
	return err
}

type UserType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
