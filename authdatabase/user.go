package authdatabase

import (
	//"bufio"
	//"bytes"
	//"crypto/sha256"
	//"database/sql"
	//"errors"
	//"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	dbh "github.com/fnuerpod/odb-obsauth/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
	//"log"
	//"strconv"
)

func (db *MCAuthDB_sqlite3) Getuser(name string) (UserInfo, bool) {
	db.logger.Debug.Println("Getting user from database (by their Discord ID)...")

	stmt, err := db.handler.FetchRowsFromTableStmt("users", []dbh.FilterItem{
		{"username", dbh.CMP_EQ, nil},
	})

	defer stmt.Close()

	var uname, passhash string
	var perm int

	if err != nil {
		goto errorState
	}

	err = stmt.QueryRow(name).Scan(&uname, &passhash, &perm)

errorState:
	if err != nil {

		return UserInfo{}, false

		db.logger.Debug.Println("Failed to obtain user from database. Invalid user(?)")
	} else {
		db.logger.Debug.Println("Managed to obtain user from database.")
	}

	return UserInfo{
		Name:      uname,
		PassHash:  passhash,
		PermLevel: perm,
	}, true

}
