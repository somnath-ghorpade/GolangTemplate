package localmongo

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReadData(c *gin.Context) {
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
	defer func() {
		cur.Close(context.Background())
		session.Disconnect(context.TODO())
	}()
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
	user.Id = guidmdl.GetGUID()
	err := c.Bind(&user)
	if err != nil {
		loggermdl.LogError("error while bind data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	mongoDAO := coremongo.GetMongoDAO(models.Collection)
	result, err := mongoDAO.SaveData(user)
	if err != nil {
		loggermdl.LogError("error while inserting data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}

func UpdateData(c *gin.Context) {
	user := models.User{}
	user.CreatedOn = time.Now().Unix()
	user.ModifiedOn = time.Now().Unix()
	user.ModifiedBy = "somnathg@mkcl.org"
	err := c.Bind(&user)
	if err != nil {
		loggermdl.LogError("error while bind data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	mongoDAO := coremongo.GetMongoDAO(models.Collection)
	filter := bson.M{"id": user.Id}
	update := bson.M{"password": user.Password, "modifiedBy": user.ModifiedBy, "modifiedOn": user.ModifiedOn}
	err = mongoDAO.Update(filter, update)
	if err != nil {
		loggermdl.LogError("error while updating data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.IndentedJSON(http.StatusOK, "SUCCESS")
}

func DeleteData(c *gin.Context) {
	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		loggermdl.LogError("error while bind data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	filter := bson.M{"id": user.Id}
	dao := coremongo.GetMongoDAO(models.Collection)
	err = dao.DeleteData(filter)
	if err != nil {
		loggermdl.LogError("error while deleting data", err)
		c.AbortWithStatusJSON(200, models.GetResponseData(nil, err.Error(), 417))
		c.Abort()
		return
	}
	c.IndentedJSON(http.StatusOK, "SUCCESS")
}
