package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"gaproject.terminator8000.net/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"fname"`
	Surname   string    `json:"lname"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Role      string    `json:"user_role"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserInfoModel struct {
	DB *sql.DB
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (m UserInfoModel) Insert(user *User) error {
	query := `
INSERT INTO user_info (fname, lname, email, password_hash, user_role, activated, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, version`
	args := []any{user.Name, user.Surname, user.Email, user.Password.hash, user.Role, user.Activated, time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_info_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserInfoModel) Get(id int64) (*User, error) {
	query := `
SELECT id, created_at, updated_at, fname, lname, email, password_hash, user_role, activated, version
FROM user_info
WHERE id = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Password.hash,
		&user.Role,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (m UserInfoModel) GetAll() ([]*User, error) {
	query := `
SELECT id, created_at, updated_at, fname, lname, email, password_hash, user_role, activated, version
FROM user_info`
	var users []*User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Password.hash,
			&user.Role,
			&user.Activated,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m UserInfoModel) GetByEmail(email string) (*User, error) {
	query := `
SELECT id, created_at, updated_at, fname, lname, email, password_hash, user_role, activated, version
FROM user_info
WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Password.hash,
		&user.Role,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserInfoModel) Update(user *User) error {
	query := `
UPDATE user_info
SET fname = $1, lname = $2, email = $3, password_hash = $4, activated = $5, updated_at = $6, version = version + 1
WHERE id = $7 AND version = $8
RETURNING version`
	args := []any{
		user.Name,
		user.Surname,
		user.Email,
		user.Password.hash,
		user.Activated,
		time.Now(),
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_info_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserInfoModel) Delete(id int64) error {
	query := `
    DELETE FROM user_info
    WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m UserInfoModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
SELECT user_info.id, user_info.created_at, user_info.fname, user_info.lname, user_info.email, user_info.password_hash, user_info.activated, user_info.version
FROM user_info
INNER JOIN tokens
ON user_info.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > now()`

	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserInfoModel) GetForAllToken() ([]User, error) {
	query := `
SELECT user_info.id, user_info.created_at, user_info.fname, user_info.lname, user_info.email, user_info.password_hash, user_info.activated, user_info.version
FROM user_info
INNER JOIN tokens
ON user_info.id = tokens.user_id
WHERE  tokens.expiry < now() AND user_info.activated = false`

	//args := []any{}
	var users []User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Password.hash,
			&user.Activated,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(user.Surname != "", "name", "must be provided")
	v.Check(len(user.Surname) <= 500, "name", "must not be more than 500 bytes long")
	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
