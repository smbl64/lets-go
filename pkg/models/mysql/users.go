package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"banisaeid.com/letsgo/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(query, name, email, hashedPassword)
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	query := `SELECT id, hashed_password from users WHERE email = ? AND active = TRUE`
	row := m.DB.QueryRow(query, email)

	var id int
	var hashedPassword []byte
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
