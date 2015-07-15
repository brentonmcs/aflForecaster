package aflShared

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type queryFn func(*mgo.Database) interface{}
type processQuery func(*mgo.Query)
type processPipeQuery func(*mgo.Pipe)

//UpdateRecord wraps the Upsert Mongo method and applys the $set property
func UpdateRecord(c *mgo.Collection, selector interface{}, model interface{}) int {
	info, err := c.Upsert(selector, bson.M{"$set": &model})
	HandleError(err)
	return info.Updated
}

//Find helper method for querying mongo
func Find(tableName string, queryBson interface{}, queryFn processQuery) interface{} {
	return newConnect(func(db *mgo.Database) interface{} {
		queryFn(db.C(tableName).Find(queryBson))
		return nil
	})
}

//Pipe helper method for querying mongo
func Pipe(tableName string, queryBson interface{}, queryFn processPipeQuery) interface{} {
	return newConnect(func(db *mgo.Database) interface{} {
		queryFn(db.C(tableName).Pipe(queryBson))
		return nil
	})
}

//Update helper for updating a collection in mongo
func Update(tableName string, index []string, query interface{}, updateModel interface{}) {
	newConnect(func(db *mgo.Database) interface{} {
		c := db.C(tableName)
		AddIndex(c, index)
		UpdateRecord(c, query, &updateModel)
		return nil
	})
}

//RemoveAll helper for revmoing a collection in mongo
func RemoveAll(tableName string, query interface{}) {
	newConnect(func(db *mgo.Database) interface{} {
		db.C(tableName).RemoveAll(query)
		return nil
	})
}

//AddIndex adds the Index as the collection is saved
func AddIndex(collection *mgo.Collection, keys []string) {
	index := mgo.Index{
		Key:        keys,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	HandleError(collection.EnsureIndex(index))
}

//NewConnect connects to the database and executes the queryFn
func newConnect(query queryFn) interface{} {

	session, db := connect()
	defer session.Close()
	return query(db)
}

func connect() (*mgo.Session, *mgo.Database) {

	mongoConnectionString := os.Getenv("MONGOLAB_URI")

	if mongoConnectionString == "" {
		mongoConnectionString = "localhost"
	}

	mongoDb := os.Getenv("MONGOLAB_DB")
	if mongoDb == "" {
		mongoDb = "aflForecaster"
	}

	session, err := mgo.Dial(mongoConnectionString)
	HandleError(err)

	if false {
		mgo.SetDebug(true)

		var aLogger *log.Logger
		aLogger = log.New(os.Stderr, "", log.LstdFlags)
		mgo.SetLogger(aLogger)
	}

	session.SetMode(mgo.Monotonic, true)
	return session, session.DB(mongoDb)
}
