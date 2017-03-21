package db

import (
	"gopkg.in/mgo.v2"
	"os"
)

func GetMongoSession() *mgo.Session {
	mongoUri := os.Getenv("MONGO_URI")
	if len(mongoUri) <= 0 {
		panic("MONGO_URI variable not set!")
	}
	session, err := mgo.Dial(mongoUri)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func GetMongoCollection(session *mgo.Session, name string) *mgo.Collection {
	return session.DB("").C(name)
}

type QueryDef func(session *mgo.Session) (result interface{}, err error)
type Command func(session *mgo.Session) (err error)

func Query(session *mgo.Session, query QueryDef) (interface{}, error) {
	newSession := session.Copy()
	defer newSession.Close()
	return query(newSession)
}

func Execute(session *mgo.Session, command Command) error {
	newSession := session.Copy()
	defer newSession.Close()
	return command(newSession)
}
