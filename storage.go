package main

import (
    "database/sql"
    "sync"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

type SQLStorage struct {
    db *sql.DB
    dbMutex sync.Mutex
}

func (s *SQLStorage) AddDomain (d string) {
    s.dbMutex.Lock()
    defer s.dbMutex.Unlock()

    stmt, err := s.db.Prepare("insert or ignore into domains (domain, added) VALUES (?, ?)")
    _, err = stmt.Exec(d, time.Now())
    if err != nil {
        sugar.Errorw("error inserting domain into db", err)
    }
}

// NewSQLClient initializes new sqlite database with domains table
func NewSQLClient() *SQLStorage {
    db, err := sql.Open("sqlite3", "./certa.sqlite")
    if err != nil {
        sugar.Fatal(err)
    }
    db.SetMaxOpenConns(1)

    // setup new domains table if fresh db
    // only accept unique domain inserts
    sqlStmt := `
    create table if not exists domains(id integer primary key, domain text not null, added timestamp not null, unique(domain));
    `
    _, err = db.Exec(sqlStmt)
    if err != nil {
        sugar.Fatal(err)
    }

    return &SQLStorage{db, sync.Mutex{}}
}