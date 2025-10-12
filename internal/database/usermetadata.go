package database

type UserMetadata struct {
	Username    string
	AvatarUrl   string
	LastLoginTs int
	CreateTs    int
	UpdateTs    int
}

func (um *UserMetadata) GetByID() error {
	query := `SELECT username, last_login_ts, create_ts, update_ts FROM user_metadata WHERE username = ?`
	row := Get().conn.QueryRow(query, um.Username)
	err := row.Scan(&um.Username, &um.LastLoginTs, &um.CreateTs, &um.UpdateTs)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMetadata) Save() error {
	query := `INSERT INTO user_metadata (username, last_login_ts, create_ts, update_ts) VALUES (?, ?, ?, ?)`
	_, err := Get().conn.Exec(query, um.Username, um.LastLoginTs, um.CreateTs, um.UpdateTs)
	return err
}

func (um *UserMetadata) Update() error {
	return nil
}

func (um *UserMetadata) Delete() error {
	return nil
}

func (um *UserMetadata) GetFirst() error {
	query := `SELECT username, avatar_url, last_login_ts, create_ts, update_ts FROM user_metadata LIMIT 1`
	row := Get().conn.QueryRow(query)
	err := row.Scan(&um.Username, &um.AvatarUrl, &um.LastLoginTs, &um.CreateTs, &um.UpdateTs)
	if err != nil {
		return err
	}
	return nil
}
