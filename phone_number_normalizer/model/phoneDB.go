package phonedb

import (
	"database/sql"
	"fmt"
	//To hide the sql driver from user
	_ "github.com/lib/pq"
)

const dbname = "phone_number_db"

//DB struct access db field from sql package
type DB struct {
	db *sql.DB
}

//Phone struct is to map value from the table
type Phone struct {
	ID    int
	Value string
}

//Close is used to close the db connection
func (db *DB) Close() error {
	return db.db.Close()
}

//GetAllRecords return all records from the db
func (db *DB) GetAllRecords() ([]Phone, error) {
	var ret []Phone
	statement := "SELECT * FROM phone_numbers"
	rows, err := db.db.Query(statement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p Phone
		err := rows.Scan(&p.ID, &p.Value)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

//UpdatePhoneRecord updates the record
func (db *DB) UpdatePhoneRecord(p Phone) error {
	statement := "UPDATE phone_numbers SET value = $2 WHERE id = $1"
	_, err := db.db.Exec(statement, p.ID, p.Value)
	return err
}

//GetOneRecord returns matched record from table
func (db *DB) GetOneRecord(in Phone) (*Phone, error) {
	var p Phone
	statement := "SELECT * FROM phone_numbers WHERE id != $1 AND value = $2"
	row := db.db.QueryRow(statement, in.ID, in.Value)
	err := row.Scan(&p.ID, &p.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//InsertRecord inserts number to the table
func (db *DB) InsertRecord(value string) error {
	statement := "INSERT INTO phone_numbers (value) VALUES ($1)"
	_, err := db.db.Exec(statement, value)
	return err
}

//DeletePhoneRecord deletes record from the table
func (db *DB) DeletePhoneRecord(p Phone) error {
	statement := "DELETE FROM phone_numbers WHERE id = $1"
	_, err := db.db.Exec(statement, p.ID)
	return err
}

//ResetDB is used to reset database
func ResetDB(driverName string, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	defer db.Close()
	if err != nil {
		return err
	}
	err = resetPhoneDB(db)
	if err != nil {
		return err
	}
	return nil
}

//Migrate is used to create the table.
func Migrate(driverName string, dataSource string) (*DB, error) {
	dataSource = fmt.Sprintf("%s dbname=%s", dataSource, dbname)
	db, err := sql.Open(driverName, dataSource)
	ret := DB{db}
	if err != nil {
		return &ret, err
	}
	err = createTableIfNotExists(db)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func createPhoneDB(db *sql.DB) error {
	statement := "CREATE DATABASE " + dbname
	_, err := db.Exec(statement)
	return err
}

func resetPhoneDB(db *sql.DB) error {
	statement := "DROP DATABASE IF EXISTS " + dbname
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return createPhoneDB(db)
}

func createTableIfNotExists(db *sql.DB) error {
	statement := `CREATE TABLE IF NOT EXISTS phone_numbers 
	(id SERIAL PRIMARY KEY,
		value VARCHAR(255))`
	_, err := db.Exec(statement)
	return err
}
