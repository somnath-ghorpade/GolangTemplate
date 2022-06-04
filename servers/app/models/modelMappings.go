package models

import (
	"encoding/json"
	"time"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/authmdl/jwtmdl"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/dalmdl/coremongo"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/errormdl"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/corepkgv2/loggermdl"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

var (
	// JWTKey - JWTKey for r and c
	JWTKey     = "gUkXp2s5v8y/B?E(H+MbQeThVmYq3t6w"
	ConfigPath = "../mongodb.json"
	Host       = "CoreStudio"
	Database   = "CoreStudio"
	Collection = "gotemplate"
)

const (
	// DefaultProjectPort -
	DefaultPort = "3000"
	PprofPort   = "6060"
)

type User struct {
	Id         string `json:"id" bson:"id"`
	UserName   string `json:"userName" bson:"userName"`
	Password   string `json:"password" bson:"password"`
	CreatedOn  int64  `json:"createdOn" bson:"createdOn"`
	CreatedBy  string `json:"createdBy" bson:"createdBy"`
	ModifiedOn int64  `json:"modifiedOn" bson:"modifiedOn"`
	ModifiedBy string `json:"modifiedBy" bson:"modifiedBy"`
}
type ResponseData struct {
	Result         interface{} `json:"result"`
	Error          interface{} `json:"error"`
	ResponseHeader interface{} `json:"reponseHeader"`
	ErrorCode      int         `json:"errorCode"`
	IsCompressed   bool        `json:"isCompressed"`
	ServerTime     time.Time   `json:"serverTime"`
}

func (r *User) Save() (string, error) {
	dao := coremongo.GetMongoDAO(Collection)
	InsId, err := dao.SaveData(r)
	if err != nil {
		return "error insert", errormdl.Wrap(err.Error())
	}
	return InsId, errormdl.Wrap(err.Error())
}

func (r *User) Get() error {
	dao := coremongo.GetMongoDAO(Collection)
	res, err := dao.GetData(bson.M{})
	if err != nil {
		return err
	}
	umErr := json.Unmarshal([]byte(res.Get("0").String()), r)
	return errormdl.Wrap(umErr.Error())
}

func (r *User) Update(i interface{}) error {
	if len(r.UserName) == 0 {
		return errormdl.Wrap("Cannot update non-existent document")
	}
	dao := coremongo.GetMongoDAO(Collection)
	err := dao.Update(bson.M{"UserName": r.UserName}, i)
	if err != nil {
		return errormdl.Wrap(err.Error())
	}
	return err
}

func DeleteRecord(UserName string) error {
	if len(UserName) == 0 {
		return errormdl.Wrap("Cannot delete non-existent document")
	}
	dao := coremongo.GetMongoDAO(Collection)
	err := dao.DeleteData(bson.M{"UserName": UserName})
	if err != nil {
		return errormdl.Wrap(err.Error())
	}
	return errormdl.Wrap(err.Error())
}

func NewInstance(userName, password, createdBy, modifiedBy string, createdOn, modifiedOn int64) User {
	return User{
		UserName:   userName,
		Password:   password,
		CreatedOn:  createdOn,
		CreatedBy:  createdBy,
		ModifiedOn: modifiedOn,
		ModifiedBy: modifiedBy,
	}
}

//GetResponseData - to get the return obj
func GetResponseData(result interface{}, errorData interface{}, errorCode int) ResponseData {
	return ResponseData{
		Result:     result,
		Error:      errorData,
		ErrorCode:  errorCode,
		ServerTime: time.Now(),
	}
}

func DecodeToken(tokenReq string) (p jwt.MapClaims, err error) {
	principle, err := jwtmdl.GeneratePricipleObjUsingToken(tokenReq, JWTKey)
	if err != nil {
		return nil, err
	}
	loggermdl.LogError("jwt", principle)
	return principle, err
}
