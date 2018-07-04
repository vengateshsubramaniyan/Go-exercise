package main

import (
	"Go-exercise/phone_number_normalizer/model"
	"fmt"
	"regexp"
)

const (
	user     = "postgres"
	password = "postgres"
	host     = "127.0.0.1"
	port     = 5432
)

func main() {
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s", host, port, user, password)
	errorCheck(phonedb.ResetDB("postgres", connectStr))
	db, err := phonedb.Migrate("postgres", connectStr)
	defer db.Close()
	errorCheck(err)
	errorCheck(db.InsertRecord("1234567890"))
	errorCheck(db.InsertRecord("123 456 7891"))
	errorCheck(db.InsertRecord("(123) 456 7892"))
	errorCheck(db.InsertRecord("(123) 456-7893"))
	errorCheck(db.InsertRecord("123-456-7894"))
	errorCheck(db.InsertRecord("123-456-7890"))
	errorCheck(db.InsertRecord("1234567892"))
	errorCheck(db.InsertRecord("(123)456-7892"))
	phonelist, err := db.GetAllRecords()
	errorCheck(err)

	for _, p := range phonelist {
		normalizedValue := normalize(p.Value)
		if normalizedValue == p.Value {
			fmt.Printf("%s -----> No change required\n", p.Value)
			nP, err := db.GetOneRecord(p)
			errorCheck(err)
			if nP != nil {
				err = db.DeletePhoneRecord(*nP)
				errorCheck(err)
			}
		} else {
			fmt.Printf("%s -----> Updating or Removing\n", p.Value)
			nP, err := db.GetOneRecord(phonedb.Phone{ID: p.ID, Value: normalizedValue})
			errorCheck(err)
			if nP == nil {
				err = db.UpdatePhoneRecord(phonedb.Phone{ID: p.ID, Value: normalizedValue})
				errorCheck(err)
			} else {
				err = db.DeletePhoneRecord(phonedb.Phone{ID: p.ID, Value: normalizedValue})
				errorCheck(err)
			}
		}
	}
}

func normalize(number string) string {
	reg, err := regexp.Compile("[^0-9]")
	errorCheck(err)
	return reg.ReplaceAllString(number, "")
}

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}
