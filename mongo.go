package goroku

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

const mongoDBName = "meetapp"

type mongodb string

var databaseName string

// MustMongoDB Must is a helper that wraps a call to a function returning (*mgo.Session, error) and panics if the error is non-nil
func MustMongoDB(ctx context.Context) *mgo.Session {
	db, ok := MongoDB(ctx)
	if !ok {
		panic("not found mongoDB")
	}
	return db
}

// MongoDB returns the connected redis client
func MongoDB(ctx context.Context) (*mgo.Session, bool) {
	key := mongodb(mongoDBName)
	session, ok := ctx.Value(key).(*mgo.Session)
	return session, ok
}

// MongoDBName return mongoDB name
func MongoDBName() string {
	return databaseName
}

// WithMockMongoDB mock mongoDB name
func WithMockMongoDB() {
	databaseName = "test_" + mongoDBName
}

// OpenMongoDB open MongoDB connections in the context's default
func OpenMongoDB(ctx context.Context) context.Context {
	uri, dbName := getHerokuMongoURI()
	databaseName = dbName

	sesh, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, mongodb(mongoDBName), sesh)
	return ctx
}

func getHerokuMongoURI() (uri string, dbName string) {
	// default
	uri = fmt.Sprintf("%s:%d", "localhost", 27017)
	dbName = mongoDBName

	mongoURI := os.Getenv("MONGOLAB_URI")
	if mongoURI == "" {
		fmt.Println("local: mongoDB", uri, dbName)
		return
	}
	mongoInfo, err := url.Parse(mongoURI)
	if err != nil {
		return
	}

	uri = mongoURI
	dbName = strings.Replace(mongoInfo.Path, "/", "", 1)
	return
}

// CloseMongoDB closes mongoDB connections in the context's
func CloseMongoDB(ctx context.Context) context.Context {
	sesh, _ := MongoDB(ctx)
	if sesh == nil {
		fmt.Println("not found mongoDB")
	}
	sesh.Close()
	ctx = context.WithValue(ctx, mongodb(mongoDBName), nil)
	return ctx
}
