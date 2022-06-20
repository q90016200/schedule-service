package service

import (
	"context"
	"fmt"
	"os"
	"scheduleService/dao"
	"scheduleService/model"

	"go.mongodb.org/mongo-driver/bson"
)

type JobServiceStruct struct{}

//JobService 用來建構 JobService 的假建構子
func JobService() JobServiceStruct {
	return JobServiceStruct{}
}

func (r JobServiceStruct) Create(data model.Job) error {
	fmt.Println(data)
	mongoConfig := dao.Config{
		Host: os.Getenv("MONGODB_HOST"),
		Port: os.Getenv("MONGODB_PORT"),
	}
	mongoClient, err := mongoConfig.Conn()
	if err != nil {
		return err
	}
	coll := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("job")
	doc := bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	return nil
}
