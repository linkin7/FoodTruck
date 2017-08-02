// package mysql implements a sql database or the UserDbManager
// interface (common/data_manager_interface.go). During inserting a new
// user data it automatically creates a unique user id.

package mysql

import (
        "database/sql"
        "fmt"
        "log"
)

type Users struct {
        db *sql.DB
}

func New(dataSourceName string) *Users {
        log.Printf("Opening mysql database: %s\n", dataSourceName)
        db, err := sql.Open("mysql", dataSourceName)
        if err != nil {
                log.Fatalf("Could not open db: %v", err)
                return nil
        }
        return &Users{
                db: db,
        }
}

func (us *Users) AddUser(name, pw, cuisine string) int64 {
        q := fmt.Sprintf("INSERT INTO users (name, password, cuisine) VALUES ('%s', '%s', '%s')", name, pw, cuisine)
        rows, err := us.db.Query(q)
        if err != nil {
                log.Printf("%v", err)
                return -1
        }
        defer rows.Close()

        return us.UserID(name)
}

func (us *Users) ValidateUser(name, pw string) bool {
        q := fmt.Sprintf("SELECT password FROM users WHERE name = '%s'", name)
        rows, err := us.db.Query(q)
        if err != nil {
                log.Printf("%v", err)
                return false
        }
        defer rows.Close()

        for rows.Next() {
                var password string
                if err := rows.Scan(&password); err != nil {
                        log.Printf("%v", err)
                        return false
                }
                return pw == password
        }
        return false
}

func (us *Users) UserID(name string) int64 {
        q := fmt.Sprintf("SELECT id FROM users WHERE name = '%s'", name)
        rows, err := us.db.Query(q)
        if err != nil {
                log.Printf("%v", err)
                return -1
        }
        defer rows.Close()

        for rows.Next() {
                // TODO: need to make sure whether its int or int64
                var uid int64
                if err := rows.Scan(&uid); err != nil {
                        log.Printf("%v", err)
                        return -1
                }
                return uid
        }
        return -1
}

func (us *Users) CuisineType(id int64) string {
        q := fmt.Sprintf("SELECT cuisine FROM users WHERE is = %v", id)
        rows, err := us.db.Query(q)
        if err != nil {
                log.Printf("%v", err)
                return "None"
        }
        defer rows.Close()

        for rows.Next() {
                // TODO: need to make sure whether its int or int64
                var cuisine string
                if err := rows.Scan(&cuisine); err != nil {
                        log.Printf("%v", err)
                        return "None"
                }
                return cuisine
        }
        return "None"
}