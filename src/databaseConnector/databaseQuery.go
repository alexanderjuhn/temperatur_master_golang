package databaseConnector

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"strings"

	_ "github.com/lib/pq"
)

var (
	db    *sql.DB
	debug bool
)

type RoomData struct {
	Id         int
	Room       string
	RoomId     int
	Temperatur float32
	Humidity   float32
	RecordDate string
	RecordTime string
	Pressure   float32
}

func InsertRoom(roomName string) {
	db := GetConnection()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	stmt := "INSERT INTO statusdata.room (id, name) SELECT nextval('statusdata.room_seq'), $1 FROM statusdata.room WHERE NOT EXISTS (SELECT 1 FROM statusdata.room WHERE name=$2) limit 1"
	_, e := tx.Exec(stmt, roomName, roomName)
	CheckInsertError(e, tx)

	err = tx.Commit()
	CheckError(err)
}

func InsertRoomData(roomData RoomData) {
	db := GetConnection()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	stmt := "INSERT INTO statusdata.roomdata(id,room_id,temperature,humidity,datecreated,pressure) VALUES(nextval('statusdata.roomdata_seq'),$1,$2,$3, now(), $4);"
	_, e := tx.Exec(stmt, roomData.RoomId, roomData.Temperatur, roomData.Humidity, roomData.Pressure)
	CheckInsertError(e, tx)

	err = tx.Commit()
	CheckError(err)
}

func ReadRoom(roomName string) int {
	db = GetConnection()
	if debug {
		log.Printf("get id for room ", roomName)
	}
	rows, err := db.Query(`SELECT id FROM statusdata.room WHERE name = $1`, roomName)
	CheckError(err)
	defer rows.Close()
	var roomId int
	rows.Next()
	err = rows.Scan(&roomId)
	CheckError(err)
	return roomId
}

func CheckInsertError(err error, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

// Return open connection
// If no connection is currently open then open a new connection
func GetConnection() *sql.DB {
	if db == nil {
		ReadConfig()
		db = GetDatabaseConnection()
		return db
	} else {
		return db
	}
}

func ProcessValue(value string) {
	value = strings.Replace(value, "'", "\"", -1)

	data := RoomData{}
	if debug {
		log.Println(value)
		log.Println(data)
	}
	json.Unmarshal([]byte(value), &data)
	data.RoomId = ReadRoom(data.Room)
	if debug {
		log.Println(value)
		log.Println(data)
	}
	InsertRoomData(data)
}
