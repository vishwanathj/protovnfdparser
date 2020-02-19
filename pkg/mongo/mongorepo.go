package mongo

import (
	"fmt"
	"log"

	"github.com/vishwanathj/protovnfdparser/pkg/config"

	"github.com/vishwanathj/protovnfdparser/pkg/datarepo"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoRepo stores mongo collection
type MongoRepo struct {
	session    *mgo.Session
	dbName     string
	collection *mgo.Collection
}

// NewMongoRepo creates a new MongoRepo
func NewMongoRepo(cfg config.Config) (datarepo.VnfdRepository, error) {
	if cfg.MongoDBConfig == nil {
		panic("mongodb config missing")
	}
	dbURI := fmt.Sprintf("%s:%d", cfg.MongoDBConfig.MongoIP, cfg.MongoDBConfig.MongoPort)
	ms, err := NewSession(dbURI)
	if err != nil {
		log.Fatal("unable to connect to mongodb:", err)
	}

	collection := ms.GetCollection(cfg.MongoDBConfig.MongoDBName, cfg.MongoDBConfig.MongoColName)
	err = collection.EnsureIndex(vnfdMgoModelIndex())
	if err != nil {
		log.Fatal("collection.EnsureIndex failed")
		panic(err)
	}
	mongo := &MongoRepo{
		session:    ms.session,
		dbName:     cfg.MongoDBConfig.MongoDBName,
		collection: collection,
	}
	return mongo, nil
}

// c is a helper method to get a collection from the session
func (m *MongoRepo) c(collection string) *mgo.Collection {
	return m.session.DB(m.dbName).C(collection)
}

// InsertVnfd stores Vnfd documents in mongo
func (m *MongoRepo) InsertVnfd(vnfd *models.Vnfd) error {
	// the conversion function below ensures that `_id` values is populated with `VNFD-{UUID} value`
	mgoModel := toVnfdMgoModel(vnfd)
	return m.collection.Insert(mgoModel)
}

func (m *MongoRepo) FindVnfdByID(vnfdID string) (*models.Vnfd, error) {
	mgoModel := vnfdMgoModel{}
	//err := m.c("vnfd").Find(bson.M{VnfdID: vnfdID}).One(&mgoModel)
	collName := m.collection.Name
	err := m.c(collName).Find(bson.M{VnfdID: vnfdID}).One(&mgoModel)
	if err != nil {
		//log.Error("Query Error:", err)
		return nil, err
	}
	return mgoModel.toModelVnfd(), nil
}

func (m *MongoRepo) FindVnfdByName(vnfdname string) (*models.Vnfd, error) {
	mgoModel := vnfdMgoModel{}
	collName := m.collection.Name
	err := m.c(collName).Find(bson.M{VnfdKey: vnfdname}).One(&mgoModel)
	if err != nil {
		//log.Error("Query Error:", err)
		return nil, err
	}
	return mgoModel.toModelVnfd(), err
}

func (m *MongoRepo) GetVnfds(start string, limit int, sort string) ([]models.Vnfd, int, error) {
	var vnfds []vnfdMgoModel
	var err error

	if len(sort) != 0 {
		//https://gist.github.com/kkdai/813f1aaf7cd58487472a
		//http://technotip.com/4101/string-comparison-mongodb/
		err = m.collection.Find(bson.M{"name": bson.M{"$gt": start}}).Sort(sort).Limit(limit).All(&vnfds)
	} else {
		err = m.collection.Find(bson.M{"name": bson.M{"$gt": start}}).Limit(limit).All(&vnfds)
	}

	if err != nil {
		log.Fatal(err)
		panic(err)
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
