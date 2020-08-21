package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var sqlDb *sql.DB
type RowID int64

func ParseRowID(rowID string) (RowID, error) {
	rowId, err := strconv.ParseInt(rowID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse row id: %v", err)
	}

	return RowID(rowId), nil
}

func ensureDatabaseExists(dbPath string) (string, error) {
	fullPath, err := filepath.Abs(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to find abs path for %s: %v", dbPath, err)
	}

	// todo: this seem unnecessarily complex. really??
	info, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to find info on the db file: %v", err)
	} else if info != nil && info.IsDir() {
		return "", fmt.Errorf("database path '%s' is a directory", fullPath)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create '%s': %v", fullPath, err)
	}

	defer file.Close()

	return fullPath, nil
}

var tableSchemas = map[string]string{
	"users": usersTable,
	"vehicles": vehiclesTable,
}

func OpenDatabase(dbPath string) {
	var err error

	dbPath, err = ensureDatabaseExists(dbPath)
	log.Printf("Database path: %s", dbPath)

	sqlDb, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open '%s': %v", dbPath, err)
	}

	tables, err := getTableList()
	if err != nil {
		log.Fatalf("Failed to get table list: %v", err)
	}

	for tableName, tableSchema := range tableSchemas {
		if _, ok := tables[tableName]; !ok {
			err = createTable(tableSchema)
			if err != nil {
				log.Fatalf("Failed to create %s table: %v", tableName, err)
			}
		}
	}
}

func createTable(query string) error {
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func getTableList() (map[string]bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table'`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}
	row, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer row.Close()

	tables := make(map[string]bool)

	for row.Next() {
		var name string
		err = row.Scan(&name)
		if err != nil {
			return nil, err
		}

		tables[name] = true
	}

	return tables, nil
}

func CloseDatabase() error {
	fmt.Println("closing database")
	err := sqlDb.Close()
	if err != nil {
		return err
	}

	return nil
}