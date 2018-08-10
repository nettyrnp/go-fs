package storage

import (
	"time"

	"gopkg.in/mgo.v2"
	"github.com/nettyrnp/go-fs/models"
	"log"
)

var (
	Host = []string{"127.0.0.1:27017"}
	//session *mgo.Session
)
const (
	//Username   = "..."
	//Password   = "..."
	Database   = "logdb"
	Collection = "logs"
)
func init() {
	// Clear DB
	resetDB()

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

//type FormatType int8
//
//const (
//	_             FormatType = iota
//	FIRST_FORMAT
//	SECOND_FORMAT
//)
//
//func (d FormatType) String() string {
//	codes := map[FormatType]string{
//		FIRST_FORMAT: "FIRST_FORMAT",
//		SECOND_FORMAT: "SECOND_FORMAT",
//	}
//	return codes[d]
//}

func Save(records []models.LogRecord) {
	session := openSession()
	defer session.Close()

	c := session.DB(Database).C(Collection)

	for _, record := range records {
		// Insert rec
		if err := c.Insert(record); err != nil {
			panic(err)
		}
	}
	log.Printf("Saving: inserted %d records into DB", len(records))

	n, err := c.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	log.Println("Saving: total records in DB:", n)

	//// Get all
	//var games []models.LogRecord
	//err := c.Find(nil).Sort("-start").All(&games)
	//if err != nil {
	//	panic(err)
	//}
	////fmt.Println("Log records", len(games))
	////fmt.Println("Log records:", games)
}
