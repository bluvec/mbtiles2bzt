package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type MapTask struct {
	bun.BaseModel `bun:"table:task"`

	ID       string  `bun:"id,pk"`
	Count    int     `bun:"count"`
	Version  float64 `bun:"version"`
	Language string  `bun:"language"`
	Date     string  `bun:"date"`
	MaxLevel int     `bun:"max_level"`
	MinLevel int     `bun:"min_level"`
}

func openDB(p string) (*bun.DB, error) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.Remove(p); err != nil {
			return nil, err
		}
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, p)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func showTables(db *bun.DB) error {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		fmt.Println("Table Name:", name)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func checkTableExists(db *bun.DB, tableName string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func dropAllIndexes(db *bun.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index'")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var indexName string
		if err := rows.Scan(&indexName); err != nil {
			tx.Rollback()
			return err
		}

		_, err := tx.Exec(fmt.Sprintf("DROP INDEX %s", indexName))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func renameTable(db *bun.DB, oldName string, newName string) error {
	_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldName, newName))
	return err
}

func createIndexOnMap(db *bun.DB) error {
	_, err := db.Exec("CREATE INDEX idx_map ON map (zoom_level, tile_column, tile_row)")
	return err
}

func createTableTask(db *bun.DB) error {
	if _, err := db.NewCreateTable().Model((*MapTask)(nil)).IfNotExists().Exec(gCtx); err != nil {
		log.Println("maptile: create table (MapMetaData) err:", err)
		return err
	}
	task := &MapTask{
		ID:       "1",
		Count:    1,
		Version:  1.0,
		Language: "en",
		Date:     "2021-01-01",
		MaxLevel: 16,
		MinLevel: 8,
	}
	if _, err := db.NewInsert().Model(task).Exec(gCtx); err != nil {
		log.Println("maptile: insert into table (MapTask) err:", err)
		return err
	}
	return nil
}

func addTypeColumnToMapWithValue(db *bun.DB, value int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("ALTER TABLE map ADD COLUMN type INTEGER")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE map SET type = ?", value)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
