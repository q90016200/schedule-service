package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"scheduleService/dao"
	"scheduleService/model"
	"time"
)

type JobServiceStruct struct{}

//JobService 用來建構 JobService 的假建構子
func JobService() JobServiceStruct {
	return JobServiceStruct{}
}

func (r JobServiceStruct) Create(data model.Job) error {
	// 寫入資料更改建立更改時間狀態
	data.ID = primitive.NewObjectID()
	data.CreatedAt = time.Now().UTC()
	data.UpdatedAt = time.Now().UTC()
	data.Status = "working"
	// 寫入 db
	mongoConfig := dao.Config{
		Host: os.Getenv("MONGODB_HOST"),
		Port: os.Getenv("MONGODB_PORT"),
	}
	mongoClient, err := mongoConfig.Conn()
	if err != nil {
		return err
	}
	coll := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("job")
	// doc := bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}}
	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	return nil
}

func (r JobServiceStruct) Query(id string) ([]*model.Job, error) {
	var results []*model.Job

	// 搜尋 db
	mongoConfig := dao.Config{
		Host: os.Getenv("MONGODB_HOST"),
		Port: os.Getenv("MONGODB_PORT"),
	}
	mongoClient, err := mongoConfig.Conn()
	if err != nil {
		return results, err
	}
	coll := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("job")

	filter := bson.M{}
	if id != "" {
		docID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": docID}
	}
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem model.Job
		err := cur.Decode(&elem)
		if err != nil {
			return results, err
		}

		results = append(results, &elem)
	}

	//fmt.Println("result:", results)

	return results, nil
}