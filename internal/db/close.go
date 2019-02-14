package db

func (db *DB) Close() error {
	if db == nil || db.db == nil {
		return ErrDatabaseClosed
	}
	return db.db.Close()
}
