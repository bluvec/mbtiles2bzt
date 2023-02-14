package main

import (
	"context"
	"flag"
	"fmt"
)

const (
	kTableNameTiles = "tiles"
	kTableNameMap   = "map"
)

var (
	gCtx context.Context
)

var (
	Version   string
	GitHash   string
	BuildTime string
	BuildMode string = "debug"
)

func main() {

	fmt.Println("VERSION   :", Version)
	fmt.Println("GIT HASH  :", GitHash)
	fmt.Println("BUILD TIME:", BuildTime)
	fmt.Println("BUILD MODE:", BuildMode)

	var filename string
	var tileType int

	flag.StringVar(&filename, "i", "", "input file (.mbtiles)")
	flag.IntVar(&tileType, "t", -1, "type of the tiles")
	flag.Parse()

	if filename == "" {
		fmt.Println("Error: input file is required")
		return
	}

	if tileType == -1 {
		fmt.Println("Error: type of the tiles is required")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	gCtx = ctx
	defer cancel()

	db, err := openDB(filename)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// check whether table tiles exist or not
	if exists, err := checkTableExists(db, kTableNameTiles); err != nil {
		fmt.Println("Error checking if table tiles exists:", err)
		return
	} else if !exists {
		fmt.Println("Error: table tiles does not exist")
		return
	}

	// remove all indexes
	if err := dropAllIndexes(db); err != nil {
		fmt.Println("Error removing indexes:", err)
		return
	}

	// rename table tiles to map
	if err := renameTable(db, "tiles", "map"); err != nil {
		fmt.Println("Error: ", "can not rename table")
	}

	// add column type to table map
	if err := addTypeColumnToMapWithValue(db, tileType); err != nil {
		fmt.Println("Error adding column type to table map:", err)
		return
	}

	// create index on map
	if err := createIndexOnMap(db); err != nil {
		fmt.Println("Error creating index on map:", err)
		return
	}

	// add table task to be compatible with legacy webui
	if err := createTableTask(db); err != nil {
		fmt.Println("Error creating table task:", err)
		return
	}

	fmt.Println("Good bye")

}
