package authentication

import (
	"database/sql"
	"ecommerce/models"
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sqlx.DB, email string, password string) (int, error) {
	var login models.Store
	SQL := `SELECT id, password FROM users WHERE user_email = $1 and archive_at is null`
	err := db.Get(&login, SQL, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return login.Id, errors.New("user not found")
		}
		return login.Id, err
	}
	err = bcrypt.CompareHashAndPassword(login.Password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return login.Id, errors.New("Incorrect password")
	} else if err != nil {
		return login.Id, err
	}
	return login.Id, nil
}
func ValidateAdmin(db *sqlx.DB, userId int) (bool, error) {
	var userType string
	SQL := `select user_type from user_role where user_id=$1 and user_type='admin'`
	err := db.QueryRowx(SQL, userId).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Create(tx *sqlx.Tx, info models.RegisterUser) (int, error) {
	var userId int
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return userId, err
	}
	SQL := `insert into users(user_name,user_email,password) values ($1,$2,$3) returning id`
	err = tx.QueryRowx(SQL, info.Name, info.Email, hash).Scan(&userId)
	if err != nil {
		return userId, err
	}
	return userId, nil
}

func AddRole(tx *sqlx.Tx, userId int, role string) error {
	SQL := `insert into user_role(user_id,user_type) values($1,$2)`
	_, err := tx.Exec(SQL, userId, role)
	if err != nil {
		return err
	}
	return nil
}
func Delete(db *sqlx.DB, userId int) error {
	SQL := `UPDATE users
			SET archive_at = current_timestamp
			WHERE id=$1;`
	_, err := db.Exec(SQL, userId)
	if err != nil {
		return err
	}
	return nil
}
