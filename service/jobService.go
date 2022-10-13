package service

import (
	"fmt"
	"os"
	"scheduleService/dao/mysql"
	"scheduleService/model"
)

type JobServiceStruct struct{}

//var mongoConfig mongodb.Config
var mysqlConfig mysql.Config

// JobService 用來建構 JobService 的假建構子
func JobService() JobServiceStruct {
	//switch os.Getenv("DATABASE") {
	//case "mongodb":
	//	mongoConfig = mongodb.Config{
	//		Host: os.Getenv("MONGODB_HOST"),
	//		Port: os.Getenv("MONGODB_PORT"),
	//	}
	//case "mysql":
	mysqlConfig = mysql.Config{
		UserName: os.Getenv("MYSQL_USER"),
		PassWord: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		DataBase: os.Getenv("MYSQL_DATABASE"),
	}
	//}

	return JobServiceStruct{}
}

func (r JobServiceStruct) Create(data model.Job) error {
	data.Status = "working"
	//switch os.Getenv("DATABASE") {
	//case "mongodb":
	//	// 寫入資料更改建立更改時間狀態
	//	data.ID = primitive.NewObjectID()
	//	data.CreatedAt = time.Now().UTC()
	//	data.UpdatedAt = time.Now().UTC()
	//
	//	// 寫入 db
	//	mongoClient, err := mongoConfig.Conn()
	//	if err != nil {
	//		return err
	//	}
	//	coll := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("job")
	//	// doc := bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}}
	//	_, err = coll.InsertOne(context.TODO(), data)
	//	if err != nil {
	//		panic(err)
	//	}
	//case "mysql":
	db, err := mysqlConfig.Conn()
	if err != nil {
		return err
	}

	result := db.Omit("ID").Create(&data)
	if result.Error != nil {
		panic(err)
	}
	//}

	return nil
}

func (r JobServiceStruct) Query(id string) ([]*model.Job, error) {
	var results []*model.Job

	//switch os.Getenv("DATABASE") {
	//case "mongodb":
	//	// 搜尋 db
	//	mongoClient, err := mongoConfig.Conn()
	//	if err != nil {
	//		return results, err
	//	}
	//	coll := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("job")
	//
	//	filter := bson.M{}
	//	if id != "" {
	//		docID, _ := primitive.ObjectIDFromHex(id)
	//		filter = bson.M{"_id": docID}
	//	}
	//	cur, err := coll.Find(context.TODO(), filter)
	//	if err != nil {
	//		return results, err
	//	}
	//
	//	for cur.Next(context.TODO()) {
	//		// create a value into which the single document can be decoded
	//		var elem model.Job
	//		err := cur.Decode(&elem)
	//		if err != nil {
	//			return results, err
	//		}
	//
	//		results = append(results, &elem)
	//	}
	//case "mysql":
	db, err := mysqlConfig.Conn()
	if err != nil {
		fmt.Println("111:", err)
		return results, err
	}

	//db.Unscoped().Find(&results)
	db.Find(&results)

	//}

	//fmt.Println("result:", results)

	return results, nil
}

func (r JobServiceStruct) Update(id string, data model.Job) error {
	//// 搜尋 db
	//mongoClient, err := mongoConfig.Conn()
	//if err != nil {
	//	return err
	//}
	//coll := mongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("job")
	//docID, _ := primitive.ObjectIDFromHex(id)
	//filter := bson.M{
	//	"_id": docID,
	//}
	//update := bson.D{{"$set", data}}
	//_, err = coll.UpdateOne(context.TODO(), filter, update)

	// mysql
	db, err := mysqlConfig.Conn()
	result := db.Model(&model.Job{}).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		err = result.Error
	}

	return err
}

func (r JobServiceStruct) Delete(id string) error {
	// mysql
	db, err := mysqlConfig.Conn()
	result := db.Where("id = ?", id).Delete(&model.Job{})

	if result.Error != nil {
		err = result.Error
	}

	return err
}
