package db

import (
	"crypto/tls"
	"errors"
	"gopkg.in/mgo.v2"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Source : https://github.com/hashicorp/vault/blob/master/plugins/database/mongodb/connection_producer.go#L101
*/
func parseMongoURL(rawURL string) (*mgo.DialInfo, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	info := mgo.DialInfo{
		Addrs:    strings.Split(url.Host, ","),
		Database: strings.TrimPrefix(url.Path, "/"),
		Timeout:  10 * time.Second,
	}

	if url.User != nil {
		info.Username = url.User.Username()
		info.Password, _ = url.User.Password()
	}

	query := url.Query()
	for key, values := range query {
		var value string
		if len(values) > 0 {
			value = values[0]
		}

		switch key {
		case "authSource":
			info.Source = value
		case "authMechanism":
			info.Mechanism = value
		case "gssapiServiceName":
			info.Service = value
		case "replicaSet":
			info.ReplicaSetName = value
		case "maxPoolSize":
			poolLimit, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("bad value for maxPoolSize: " + value)
			}
			info.PoolLimit = poolLimit
		case "ssl":
			// Unfortunately, mgo doesn't support the ssl parameter in its MongoDB URI parsing logic, so we have to handle that
			// ourselves. See https://github.com/go-mgo/mgo/issues/84
			ssl, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New("bad value for ssl: " + value)
			}
			if ssl {
				info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
					return tls.Dial("tcp", addr.String(), &tls.Config{})
				}
			}
		case "connect":
			if value == "direct" {
				info.Direct = true
				break
			}
			if value == "replicaSet" {
				break
			}
			fallthrough
		default:
			return nil, errors.New("unsupported connection URL option: " + key + "=" + value)
		}
	}

	return &info, nil
}

func GetMongoSession() *mgo.Session {
	uri := os.Getenv("MONGO_URI")
	if len(uri) <= 0 {
		panic("MONGO_URI variable not set!")
	}

	info, err := parseMongoURL(uri)
	if err != nil {
		panic(err)
	}

	session, err := mgo.DialWithInfo(info)

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
