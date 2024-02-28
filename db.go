package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MySQL")

	return &MySQLStorage{
		db: db,
	}
}

func (s *MySQLStorage) Init() (*sql.DB, error) {
	// initialize the tables
	if err := s.createProjectsTable(); err != nil {
		return nil, err
	}

	if err := s.createUsersTable(); err != nil {
		return nil, err
	}

	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *MySQLStorage) createProjectsTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createTasksTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		status ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE') NOT NULL DEFAULT 'TODO',
		projectId INT NOT NULL,
		createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		assignedToID INT NOT NULL,

		FOREIGN KEY (projectId) REFERENCES projects(id),
		FOREIGN KEY (assignedToID) REFERENCES users(id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)

	return err
}

func (s *MySQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		firstName VARCHAR(255) NOT NULL,
		lastName VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		UNIQUE KEY (email)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)

	return err
}
