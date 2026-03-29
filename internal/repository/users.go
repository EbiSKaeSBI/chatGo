package repository

import "golang.org/x/crypto/bcrypt"

func (r *Repository) CreateUser(username, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	_, err = r.db.Exec("INSERT INTO users(username,password_hash) VALUES($1,$2)", username, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindUser(username string) (userId int64, err error) {
	row := r.db.QueryRow("SELECT id FROM users WHERE username = $1", username)
	err = row.Scan(&userId)
	return userId, err
}
