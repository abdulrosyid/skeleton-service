package mysqldb

import "database/sql"

type RepoSQL interface {
	GetUserMysqlDbRepository() UserRepository
}

type repoSQL struct {
	db *sql.DB
}

func NewRepoSQL(db *sql.DB) RepoSQL {
	return &repoSQL{db: db}
}
