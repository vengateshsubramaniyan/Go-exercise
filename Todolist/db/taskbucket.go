package taskbucket

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var bucketName = "tasks"
var completedBucket = "completedlist"
var db *bolt.DB

//Task struct is used to store the key, value pair of tasks in the Todolist.
type Task struct {
	Key   int
	Value string
}

//Init function to intialize db connection.
func Init(dbName string) {
	var err error
	db, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	ErrorReporter(err)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	ErrorReporter(err)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(completedBucket))
		return err
	})
	ErrorReporter(err)
}

//AddTask is used to add task to the bolt.db
func AddTask(task string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		id, _ := b.NextSequence()
		bID := itob(id)
		err := b.Put(bID, []byte(task))
		return err
	})
	ErrorReporter(err)
}

//ListTask is used to list the tasks from the Todolist
func ListTask() []Task {
	var tasks []Task
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{int(btoi(k)), string(v)})
		}
		return nil
	})
	return tasks
}

//DoTask is used to mark the task as completed
func DoTask(key int) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete(itob(uint64(key)))
	})
	ErrorReporter(err)
}

//AddToCompletedList add task to the completed list.
func AddToCompletedList(task string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(completedBucket))
		id, _ := b.NextSequence()
		bID := itob(id)
		err := b.Put(bID, []byte(task))
		return err
	})
	ErrorReporter(err)
}

//ListCompletedTask list the completed task from the Todolist.
func ListCompletedTask() []Task {
	var tasks []Task
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(completedBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{int(btoi(k)), string(v)})
		}
		return nil
	})
	return tasks
}

//ErrorReporter to handle the
func ErrorReporter(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func btoi(val []byte) uint64 {
	return binary.BigEndian.Uint64(val)
}
func itob(val uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, val)
	return b
}
