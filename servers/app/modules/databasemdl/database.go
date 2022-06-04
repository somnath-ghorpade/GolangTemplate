package databasemdl

import (
	"context"
	"encoding/json"
	"golangtemplate/servers/app/models"
	"net/http"
	"time"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/dalmdl/coremongo"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/loggermdl"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/utiliymdl/guidmdl"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func ReadData(c *gin.Context) {
	user := models.User{}
	dataa := user.Get()
	loggermdl.LogError("data", dataa)
	result, err := fetchCollectionData(bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.IndentedJSON(http.StatusOK, result.Value())
}

func fetchCollectionData(i interface{}) (gjson.Result, error) {
	session, sessionErr := coremongo.GetMongoConnection(models.Host)
	if sessionErr != nil {
		loggermdl.LogError("session error", sessionErr)
		return gjson.Result{}, sessionErr
	}
	collection := session.Database(models.Database).Collection(models.Collection)
	cur, err := collection.Find(context.Background(), i, &options.FindOptions{})
	if err != nil {
		loggermdl.LogError(err)
		return gjson.Result{}, err
	}
	defer cur.Close(context.Background())
	var results []interface{}
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			loggermdl.LogError(err)
			return gjson.Result{}, err
		}
		results = append(results, result)
	}
	ba, marshErr := json.Marshal(results)
	if marshErr != nil {
		loggermdl.LogError(marshErr)
		return gjson.Result{}, marshErr
	}
	rs := gjson.ParseBytes(ba)
	return rs, err
}

func InsertData(c *gin.Context) {
	user := models.User{}
	user.CreatedOn = time.Now().Unix()
	user.ModifiedOn = time.Now().Unix()
	err := c.Bind(&user)
	if err != nil {
		loggermdl.LogError("error while bind data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	res, err := insertCollectionData(user)
	if err != nil {
		loggermdl.LogError("error while inserting data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.IndentedJSON(http.StatusOK, res.InsertedID)
}

func insertCollectionData(result models.User) (insertResult *mongo.InsertOneResult, err error) {
	session, sessionErr := coremongo.GetMongoConnection(models.Host)
	if sessionErr != nil {
		loggermdl.LogError("session error", sessionErr)
		return nil, sessionErr
	}
	collection := session.Database(models.Database).Collection(models.Collection)
	insertRes, err := collection.InsertOne(context.TODO(), result)
	if err != nil {
		loggermdl.LogError(err)
		return nil, err
	}
	return insertRes, err
}

func UpdateData(c *gin.Context) {
	user := models.User{}
	user.CreatedOn = time.Now().Unix()
	user.ModifiedOn = time.Now().Unix()
	user.ModifiedBy = "somnathg@mkcl.org"
	user.Id = guidmdl.GetGUID()
	err := c.Bind(&user)
	if err != nil {
		loggermdl.LogError("error while bind data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	result, err := updateCollectionData(user)
	if err != nil {
		loggermdl.LogError("error while updating data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.String(http.StatusOK, "result", result)
}

func updateCollectionData(user models.User) (updateRes *mongo.UpdateResult, err error) {
	session, sessionErr := coremongo.GetMongoConnection(models.Host)
	if sessionErr != nil {
		loggermdl.LogError("session error", sessionErr)
		return nil, sessionErr
	}
	collection := session.Database(models.Database).Collection(models.Collection)
	filter := bson.M{"userName": user.UserName}
	update := bson.M{"$set": bson.M{"password": user.Password, "dd": user.ModifiedBy, "modifiedOn": user.ModifiedOn}}
	opts := options.Update().SetUpsert(false)
	result, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		loggermdl.LogError(err)
		return nil, err
	}
	return result, err
}

func DeleteData(c *gin.Context) {
	userName := "somnathg"
	result, err := deleteCollectionData(userName)
	if err != nil {
		loggermdl.LogError("error while deleting data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.String(http.StatusOK, "result", result)
}

func deleteCollectionData(userName string) (deleteRes *mongo.DeleteResult, err error) {
	session, sessionErr := coremongo.GetMongoConnection(models.Host)
	if sessionErr != nil {
		loggermdl.LogError("session error", sessionErr)
		return nil, sessionErr
	}
	collection := session.Database(models.Database).Collection(models.Collection)
	filter := bson.M{"userName": userName}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		loggermdl.LogError(err)
		return nil, err
	}
	return result, err
}
