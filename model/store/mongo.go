package store

import (
	"fmt"
	"os"
	"time"

	"github.com/mgutz/logxi/v1"
	"gopkg.in/mgo.v2"
)

const (
	envMongoAuth    = "MANAGEME_MONGO_AUTH"
	envMongoHost    = "MANAGEME_MONGO_HOST"
	envDatabaseName = "MANAGEME_MONGO_DATABASE"
)

var (
	logger = log.New("store")

	mongoAuth, mongoHost, databaseName string

	defMongoAuth    = "mdbmanageme:manageme"
	defMongoHost    = "localhost:27017"
	defDatabaseName = "manageme"

	mongo *mgo.Session
)

type MongoStore struct {
	s *mgo.Session
}

func getMongoURL() string {
	return fmt.Sprintf("mongodb://%v@%v/%v", mongoAuth, mongoHost, databaseName)
}

// InitMongoSession resets the mongo session pointer with updated connection info
func InitMongoSession() error {
	// To avoid a socket leak
	if mongo != nil {
		CleanupMongoSession()
	}

	// Reset connection params
	warn := func(env, def string) {
		logger.Warn(fmt.Sprintf("env.%s not defined, defaulting to %v", env, def))
	}
	if mongoAuth = os.Getenv(envMongoAuth); len(mongoAuth) == 0 {
		warn(envMongoAuth, defMongoAuth)
		mongoAuth = defMongoAuth
	}
	if mongoHost = os.Getenv(envMongoHost); len(mongoHost) == 0 {
		warn(envMongoHost, defMongoHost)
		mongoHost = defMongoHost
	}
	if databaseName = os.Getenv(envDatabaseName); len(databaseName) == 0 {
		warn(envDatabaseName, defDatabaseName)
		databaseName = defDatabaseName
	}

	// Establish new session
	logger.Debug("init mongo", "host", mongoHost, "url", getMongoURL())
	var err error
	if mongo, err = mgo.Dial(getMongoURL()); err != nil {
		return err
	}

	// Ensure indicies
	ensureUserIndex()

	return nil
}

// CleanupMongoSession closes the current session and sets the pointer to nil
func CleanupMongoSession() {
	if mongo == nil {
		return
	}
	mongo.Close()
	time.Sleep(time.Second)

	mongo = nil
}

// Nuke destroys the database if it is in a test environment
func Nuke() error {
	if testing := os.Getenv("TESTING"); testing == "true" && mongo != nil {
		return mongo.DB(databaseName).DropDatabase()
	}
	return fmt.Errorf("env.TESTING must be set to true")
}

// NewMongoStore returns an instance of the store with a copied mongo session
// error is 500 if mongo ping fails
func NewMongoStore() (*MongoStore, error) {
	if err := mongo.Ping(); err != nil {
		return nil, err
	}
	return &MongoStore{s: mongo.Copy()}, nil
}

// Cleanup closes the mongo session of this store object
func (m *MongoStore) Cleanup() {
	m.s.Close()
}

// GetDatabase returns a pointer to an mgo database object
func (m *MongoStore) GetDatabase() *mgo.Database {
	return m.s.DB(databaseName)
}
