package repository

func (r *Repository) SaveMessage(userId int, text string) error {
	_, err := r.db.Exec("INSERT INTO messages (user_id, text) values ($1, $2)", userId, text)
	if err != nil {
		return err
	}
	return nil
}

//func (r *Repository) GetMessage() (string, error) {
//	_, err := r.db.Query("SELECT * FROM messages LIMIT 1")
//	if err != nil {
//		return "", err
//	}
//	return "row", nil
//}
