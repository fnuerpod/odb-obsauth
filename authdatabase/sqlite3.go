package authdatabase

import (
	"database/sql"
	"path/filepath"

	"github.com/fnuerpod/odb-obsauth/config"
	"github.com/fnuerpod/odb-obsauth/logging"

	dbh "github.com/fnuerpod/odb-obsauth/sqlite3dbh"
	_ "github.com/mattn/go-sqlite3"
)

type MCAuthDB_sqlite3 struct {
	handler *dbh.DBHandler
	logger  *logging.Logger
}

func createDB(handler *dbh.DBHandler) error {

	// TODO(ultrabear) some columns here use the 32 bit sql INT type to store unix timestamps
	// Anyone familiar with the 2038 problem will know that this will overflow in 2038
	// To solve this, all unix timestamps should be moved to BIGINT or T_INT64 and all databases should be migrated
	// This TODO is not immediate but will destroy the site if it is not patched

	var err error

	// permLevels
	// 1 - stream view
	// 2 - stream publish
	// 3 - both

	if _, err = handler.MakeNewTable("users", [][2]string{
		{"username", dbh.T_STR + dbh.T_PK},
		{"passHash", dbh.T_STR},
		{"permissionLevel", dbh.T_INT},
	}); err != nil {
		goto errorState
	}

errorState:
	return err
}

type UserInfo struct {
	Name      string
	PassHash  string
	PermLevel int
}

type UserError struct {
	err    string
	exists bool
	dberr  bool
}

func (UE UserError) Error() string    { return UE.err }
func (UE UserError) UserExists() bool { return UE.exists }
func (UE UserError) DBError() bool    { return UE.dberr }

// db initialiser
func InitSQLite3DB(logger *logging.Logger) *MCAuthDB_sqlite3 {
	diskdb, err := sql.Open("sqlite3", filepath.Join(config.GetDataDir(), "site.db"))
	logger.Debug.Println("Database opened, creating tables if they don't exist...")

	authDB := MCAuthDB_sqlite3{
		handler: dbh.OpenFromDB(diskdb),
		logger:  logger,
	}

	err = createDB(authDB.handler)

	if err != nil {
		logger.Err.Println("Error occurred while initialising on-disk database - more information below.")
	}

	logger.Debug.Println("Database initialised OK.")

	return &authDB
}
