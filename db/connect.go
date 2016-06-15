package db

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	// Session stores mongo session
	Session *mgo.Session

	// Mongo stores the mongodb connection string information
	Mongo *mgo.DialInfo
)

const (
	// MongoDBUrl is the default mongodb url that will be used to connect to the
	// database.
	MongoDBUrl = "mongodb://localhost:27017/notifications_testing"
)

// Connect connects to mongodb
func Connect() {
	uri := os.Getenv("MONGODB_URL")

	if len(uri) == 0 {
		uri = MongoDBUrl
	}

	Session, Mongo = getNewSession(uri)
	fmt.Println("Connected to", uri)
}

func getNewSession(connString string) (*mgo.Session, *mgo.DialInfo) {
	// quick hack to allow SSL based connections, may be removed in future when parseURL supports it
	// see also: https://github.com/go-mgo/mgo/issues/84
	const SSL_SUFFIX = "?ssl=true"
	useSsl := false

	if strings.HasSuffix(connString, SSL_SUFFIX) {
		connString = strings.TrimSuffix(connString, SSL_SUFFIX)
		useSsl = true
	}

	dialInfo, err := mgo.ParseURL(connString)
	if err != nil {
		panic(err)
	}

	dialInfo.Timeout = 10 * time.Second

	if useSsl {
		config := tls.Config{}
		config.InsecureSkipVerify = true

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &config)
		}
	}

	// get a mgo session
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	return session, dialInfo
}
