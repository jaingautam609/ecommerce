package authentication

import (
	"database/sql"
	"ecommerce/models"
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func VerifyEmailOtp(db *sqlx.DB, otp models.EmailVerify) (string, error) {
	var email string
	SQL := `SELECT user_email FROM opt where otp = $1 and user_email=$2 and expired_at > current_timestamp`
	err := db.Get(&email, SQL, otp.Otp, otp.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return email, errors.New("user not found")
		}
		return email, err
	}
	return email, nil
}
func VerifyNumberOtp(db *sqlx.DB, otp models.NumberVerify) (string, error) {
	var number string
	SQL := `SELECT phone_no FROM opt where otp = $1 and phone_no=$2 and expired_at > current_timestamp`
	err := db.Get(&number, SQL, otp.Otp, otp.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return number, errors.New("user not found")
		}
		return number, err
	}
	return number, nil
}

func StoreEmailOtp(db *sqlx.DB, otp, email string) error {
	SQL := `insert into opt(user_email,otp) values ($1,$2)`
	_, err := db.Exec(SQL, email, otp)
	if err != nil {
		return err
	}
	return nil
}
func StoreNumberOtp(db *sqlx.DB, otp, number string) error {
	SQL := `insert into opt(phone_no,otp) values ($1,$2)`
	_, err := db.Exec(SQL, number, otp)
	if err != nil {
		return err
	}
	return nil
}

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
func CheckVerified(db *sqlx.DB, info models.RegisterUser) (bool, error) {
	var flag bool
	SQL := `select count(*) >0 from users 
                   where user_email=$1 
                             and is_verified_by_email=true 
                      or is_verified_by_phone=true `
	err := db.Get(&flag, SQL, info.Email)
	return flag, err
}
func CreateUserByEmail(tx *sqlx.Tx, info models.RegisterUser) (int, error) {
	var userId int
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return userId, err
	}
	SQL := `UPDATE users
				SET user_name=$1,
				    password=$2
					WHERE user_email=$3
			RETURNING id;
`
	err = tx.QueryRowx(SQL, info.Name, hash, info.Email).Scan(&userId)
	if err != nil {
		return userId, err
	}
	return userId, nil
}
func CreateUserByNumber(tx *sqlx.Tx, info models.RegisterUser) (int, error) {
	var userId int
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return userId, err
	}
	SQL := `UPDATE users
				SET user_name=$1,
				    password=$2
					WHERE phone_no=$3
			RETURNING id;
`
	err = tx.QueryRowx(SQL, info.Name, hash, info.PhoneNumber).Scan(&userId)
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
