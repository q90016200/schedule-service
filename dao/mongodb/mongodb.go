package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Config struct {
	AuthMechanism string
	UserName      string
	PassWord      string
	Host          string
	Port          string
	ReplicaSet    string
}

func (config *Config) Conn() (dbClient *mongo.Client, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 180*time.Second)

	// if !reflect.DeepEqual(config, Config{}) {
	//	authMechanism = config.AuthMechanism
	//	host = config.Host
	//	port = config.Port
	//	userName = config.UserName
	//	passWord = config.PassWord
	//	replicaSet = config.ReplicaSet
	// }

	hostList := strings.Split(config.Host, ";")
	url := "mongodb://" + strings.Join(hostList, ":"+config.Port+",") + ":" + config.Port

	if config.ReplicaSet != "" {
		url = url + "/?replicaSet=" + config.ReplicaSet
	}

	opts := options.Client().ApplyURI(url)
	if config.AuthMechanism == "PLAIN" {
		credential := options.Credential{
			AuthMechanism: config.AuthMechanism,
			Username:      config.UserName,
			Password:      config.PassWord,
		}
		opts = opts.SetAuth(credential)
	} else if config.AuthMechanism == "SCRAM" {
		credential := options.Credential{
			Username: config.UserName,
			Password: config.PassWord,
		}
		opts = opts.SetAuth(credential)
	}

	dbClient, err = mongo.NewClient(opts)

	if err != nil {
		return nil, errors.New("Mongo NewClient Error: " + err.Error())
	}

	err = dbClient.Connect(ctx)
	if err != nil {
		return nil, errors.New("Mongo Connect Error: " + err.Error())
	}

	return
}

func Reset(config *Config) (dbClient *mongo.Client, err error) {
	dbClient, err = config.Conn()
	return
}

// IsAlive 確認連線
func IsAlive(dbClient *mongo.Client) bool {
	ctx, _ := context.WithTimeout(context.Background(), 180*time.Second)
	err := dbClient.Ping(ctx, nil)
	//defer Disconnect(dbClient)
	if err != nil {
		return false
	}

	return true
}
