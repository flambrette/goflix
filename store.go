package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/glebarez/go-sqlite"
)

type Store interface {
	Open() error
	Close() error

	GetMovies() ([]*Movie, error)
	GetMovie(id int64) (*Movie, error)
	CreateMovie(m *Movie) (error)
	UpdateMovie(m *Movie) (error)
	DeleteMovie(id int64) (error)
	AuthenticateUser(username string, password string) (bool, error)
}

type dbStore struct {
	db *sqlx.DB
}

var schema = `CREATE TABLE IF NOT EXISTS movie
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	release_date TEXT,
	duration INTEGER,
	trailer_url TEXT
);

CREATE TABLE IF NOT EXISTS user
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	password TEXT
);

INSERT INTO user (username, password) SELECT 'admin', 'admin' 
WHERE NOT EXISTS (SELECT 1 FROM user WHERE username = 'admin');
`

func (store *dbStore) Open() error{
	db, err := sqlx.Connect("sqlite", "goflix.db")
	if err != nil {
		return err
	}
	log.Println("Connect to DB")
	db.MustExec(schema)
	store.db = db
	return nil
}

func (store *dbStore) Close() error{
	return store.db.Close()
}

func (store *dbStore) GetMovies() ([]*Movie, error){
	var movies []*Movie
	err := store.db.Select(&movies, "SELECT * FROM movie")
	if err != nil {
		return movies, err
	}
	return movies, nil
}

func (store *dbStore) GetMovie(id int64) (*Movie, error){
	var movie = &Movie{}
	err := store.db.Get(movie, "SELECT * FROM movie where id = $1", id)
	if err != nil {
		return movie, err
	}
	return movie, nil
}

func (store *dbStore) CreateMovie(m *Movie) (error){
	res, err := store.db.Exec("INSERT INTO movie (title, release_date, duration, trailer_url) VALUES(?,?,?,?)",
		m.Title,m.ReleaseDate,m.Duration, m.TrailerUrl)
	if err != nil {
		return err
	}
	m.Id, err = res.LastInsertId()
	return err
}

func (store *dbStore) UpdateMovie(m *Movie) (error){
	_, err := store.db.Exec("UPDATE movie SET title = ?, release_date = ?, duration = ?, trailer_url = ? WHERE id = ?",
		m.Title,m.ReleaseDate,m.Duration, m.TrailerUrl, m.Id)
	if err != nil {
		return err
	}
	return err
}

func (store *dbStore) DeleteMovie(id int64) (error){
	_, err := store.db.Exec("DELETE FROM movie WHERE id = ?", id)
	if err != nil {
		return err
	}
	
	return nil
}

func (store *dbStore) AuthenticateUser(username string, password string) (bool, error){
	var count int
	err := store.db.Get(&count, "SELECT count(id) FROM user where username =$1 and password =$2", username, password)
	if err != nil {
		return false, err
	}
	
	return count == 1, nil
}