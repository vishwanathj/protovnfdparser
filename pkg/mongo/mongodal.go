package mongo

import (
	"fmt"
	"log"

	"github.com/vishwanathj/protovnfdparser/pkg/config"

	"github.com/vishwanathj/protovnfdparser/pkg/dataaccess"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDAL stores mongo collection
type MongoDAL struct {
	session    *mgo.Session
	dbName     string
	collection *mgo.Collection
}

// NewMongoDAL creates a new MongoDAL
func NewMongoDAL(cfg config.Config) (dataaccess.VnfdRepository, error) {
	if cfg.MongoDBConfig == nil {
		panic("mongodb config missing")
	}
	dbURI := fmt.Sprintf("%s:%d", cfg.MongoDBConfig.MongoIP, cfg.MongoDBConfig.MongoPort)
	ms, err := NewSession(dbURI)
	if err != nil {
		log.Fatal("unable to connect to mongodb:", err)
	}

	collection := ms.GetCollection(cfg.MongoDBConfig.MongoDBName, cfg.MongoDBConfig.MongoColName)
	collection.EnsureIndex(vnfdMgoModelIndex())
	mongo := &MongoDAL{
		session:    ms.session,
		dbName:     cfg.MongoDBConfig.MongoDBName,
		collection: collection,
	}
	return mongo, nil
}

// c is a helper method to get a collection from the session
func (m *MongoDAL) c(collection string) *mgo.Collection {
	return m.session.DB(m.dbName).C(collection)
}

// Insert stores documents in mongo
func (m *MongoDAL) Insert(collectionName string, docs ...interface{}) error {
	return m.c(collectionName).Insert(docs)
}

// InsertVnfd stores Vnfd documents in mongo
func (m *MongoDAL) InsertVnfd(vnfd *models.Vnfd) error {
	// the conversion function below ensures that `_id` values is populated with `VNFD-{UUID} value`
	mgoModel := toVnfdMgoModel(vnfd)
	return m.collection.Insert(mgoModel)
}

func (m *MongoDAL) FindVnfdByID(vnfdID string) (*models.Vnfd, error) {
	mgoModel := vnfdMgoModel{}
	err := m.c("vnfd").Find(bson.M{VnfdID: vnfdID}).One(&mgoModel)
	if err != nil {
		//log.Error("Query Error:", err)
		return nil, err
	}
	return mgoModel.toModelVnfd(), nil
}

func (m *MongoDAL) FindVnfdByName(vnfdname string) (*models.Vnfd, error) {
	mgoModel := vnfdMgoModel{}
	err := m.c("vnfd").Find(bson.M{VnfdKey: vnfdname}).One(&mgoModel)
	if err != nil {
		//log.Error("Query Error:", err)
		return nil, err
	}
	return mgoModel.toModelVnfd(), err
}

func (m *MongoDAL) GetVnfds(start string, limit int) ([]models.Vnfd, int, error) {
	var vnfds []vnfdMgoModel

	if start == "" {
		m.collection.Find(bson.M{"name": bson.M{"$gt": ""}}).Limit(limit).All(&vnfds)
		//p.collection.Find(nil).Limit(limit).All(&vnfds)
	} else {
		//https://gist.github.com/kkdai/813f1aaf7cd58487472a
		//http://technotip.com/4101/string-comparison-mongodb/
		m.collection.Find(bson.M{"name": bson.M{"$gt": start}}).Limit(limit).All(&vnfds)
	}
	//log.WithFields(log.Fields{"VNFDS": vnfds}).Debug("GET_VNFDS")

	count, err := m.collection.Count()
	if err != nil {
		return []models.Vnfd{}, count, err
	}

	// below converts mongo model to generic model
	var modelVnfdsArr []models.Vnfd
	for i := 0; i < len(vnfds); i++ {
		modelVnfd := *vnfds[i].toModelVnfd()
		modelVnfdsArr = append(modelVnfdsArr, modelVnfd)
	}

	return modelVnfdsArr, count, nil
}
