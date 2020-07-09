package util

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MgoClient = &MongoObj{}
)

type MongoObj struct {
	client *mongo.Client
}

func MongoInit() {
	var err error
	mongoConn := beego.AppConfig.String("mongo.conn")
	mongoAuthSource := beego.AppConfig.String("mongo.authsource")
	mongouser := beego.AppConfig.String("mongo.user")
	mongopwd := beego.AppConfig.String("mongo.password")
	uri := "mongodb://" + mongoConn
	clientOptions := options.Client().ApplyURI(uri)

	//认证参数设置，否则连不上
	opts := &options.ClientOptions{}
	opts.SetAuth(options.Credential{
		AuthSource: mongoAuthSource,
		Username:   mongouser,
		Password:   mongopwd})

	// Connect to MongoDB
	MgoClient.client, err = mongo.Connect(context.TODO(), clientOptions, opts)

	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	err = MgoClient.client.Ping(context.TODO(), nil)
}

func (m *MongoObj) InsertOne(db, tableName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := m.client.Database(db).Collection(tableName)
	insertResult, err := collection.InsertOne(context.TODO(), document)
	return insertResult, err
}

func (m *MongoObj) CountDocuments(db, tableName string, filter interface{}) (int64, error) {
	collection := m.client.Database(db).Collection(tableName)
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (m *MongoObj) Find(db, tableName string, document, filter interface{}, opts ...*options.FindOptions) (int64, error) {
	collection := m.client.Database(db).Collection(tableName)
	total, _ := m.CountDocuments(db, tableName, filter)
	fmt.Println(total)
	cur, err := collection.Find(context.TODO(), filter, opts...)
	if err != nil {
		return 0, err
	}
	if err := cur.Err(); err != nil {
		return 0, err
	}
	err = cur.All(context.TODO(), document)
	defer cur.Close(context.TODO())
	return total, err
}

func (m *MongoObj) FindOne(db, tableName string, filter, document interface{}, opts ...*options.FindOneOptions) error {
	collection := m.client.Database(db).Collection(tableName)
	err := collection.FindOne(context.TODO(), filter, opts...).Decode(document)
	return err
}

func (m *MongoObj) UpdateOne(db, tableName string, filter, document interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection := m.client.Database(db).Collection(tableName)
	res, err := collection.UpdateOne(context.Background(), filter, document, opts...)
	return res, err
}
