package storage

import (
	"time"

	"gopkg.in/mgo.v2"
	"github.com/nettyrnp/go-fs/models"
	"log"
)

var (
	Host = []string{"127.0.0.1:27017"}
)
const (
	Username   = "..."
	Password   = "..."
	Database   = "logdb"
	Collection = "logs"
)
func init() {
	// Clear DB
	resetDB()

	// Set timezone
	timezone := "UTC"
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err.Error())
	}
	time.Local = loc
}

func resetDB() {
	session := openSession()
	defer session.Close()

	err := session.DB(Database).DropDatabase()
	if err != nil {
		panic(err)
	}
}

// Uncomment the commented-out lines, when contacting remote MongoDB server
func openSession() *mgo.Session {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
		// Username: Username,
		// Password: Password,
		// Database: Database,
		// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
		// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
		// },
	})
	if err != nil {
		panic(err)
	}
	return session
}

func Save(records []models.LogRecord) {
	session := openSession()
	defer session.Close()

	c := session.DB(Database).C(Collection)

	// TODO: make batch insert
	for _, record := range records {
		if err := c.Insert(record); err != nil {
			panic(err)
		}
	}
	log.Printf("Saving: inserted %d records into DB", len(records))

	n, err := c.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	log.Println("Records in DB:", n)

}
