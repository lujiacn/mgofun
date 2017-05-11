package mgofun

import (
	"errors"
	"reflect"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MgoFun wrap all common functions
type MgoFun struct {
	model         interface{}
	session       *mgo.Session
	collection    *mgo.Collection
	logCollection *mgo.Collection // for change log
	Query         bson.M
	Sort          []string
	Skip          int
	Limit         int
}

//NewMgoFun initiate with input model and mgo session
func NewMgoFun(s *mgo.Session, dbName string, model interface{}) *MgoFun {
	mgoFun := &MgoFun{model: model, session: s}
	collection := Collection(s, dbName, model)
	logCollection := Collection(s, dbName, "ChangeLog")
	mgoFun.logCollection = logCollection // for change log
	mgoFun.collection = collection
	return mgoFun
}

// Collection conduct mgo.Collection
func Collection(s *mgo.Session, dbName string, m interface{}) *mgo.Collection {
	cName := getModelName(m)
	return s.DB(dbName).C(cName)
}

//getModelName reflect string name from model
func getModelName(m interface{}) string {
	var c string
	switch m.(type) {
	case string:
		c = m.(string)
	default:
		typ := reflect.TypeOf(m)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		c = typ.Name()
	}
	return c
}

//Create
func (m *MgoFun) Create() error {
	//generate new object Id
	id := reflect.ValueOf(m.model).Elem().FieldByName("Id")
	id.Set(reflect.ValueOf(bson.NewObjectId()))
	x := reflect.ValueOf(m.model).Elem().FieldByName("CreatedAt")
	x.Set(reflect.ValueOf(time.Now()))
	_, err := m.collection.Upsert(bson.M{"_id": id.Interface()}, bson.M{"$set": m.model})
	return err
}

//General Save method
func (m *MgoFun) Save() error {
	id := reflect.ValueOf(m.model).Elem().FieldByName("Id")
	x := reflect.ValueOf(m.model).Elem().FieldByName("UpdatedAt")
	x.Set(reflect.ValueOf(time.Now()))
	_, err := m.collection.Upsert(bson.M{"_id": id.Interface()}, bson.M{"$set": m.model})
	return err
}

//General Save method
func (m *MgoFun) SaveWithoutTime() error {
	id := reflect.ValueOf(m.model).Elem().FieldByName("Id")
	_, err := m.collection.Upsert(bson.M{"_id": id.Interface()}, bson.M{"$set": m.model})
	return err
}

//SaveWithLog will save old record to ChangeLog model
func (m *MgoFun) SaveWithLog(oldRecord interface{}, by, reason string) error {
	var err error
	err = m.Save()
	if err != nil {
		return err
	}
	err = m.saveLog(oldRecord, by, reason)
	if err != nil {
		return err
	}
	return nil
}

//SaveWithLog
func (m *MgoFun) saveLog(record interface{}, by, reason string) error {
	recordId := reflect.ValueOf(m.model).Elem().FieldByName("Id").Interface().(bson.ObjectId)

	cl := new(ChangeLog)
	cl.Id = bson.NewObjectId()
	cl.CreatedBy = by
	cl.CreatedAt = time.Now()
	cl.ChangeReason = reason
	cl.ModelObjId = recordId
	cl.ModelName = getModelName(record)
	cl.ModelValue = record
	_, err := m.logCollection.Upsert(bson.M{"_id": cl.Id}, bson.M{"$set": cl})
	return err

}

// Remove is softe delete
func (m *MgoFun) Remove() error {
	id := reflect.ValueOf(m.model).Elem().FieldByName("Id")
	if !id.IsValid() {
		return errors.New("No Id defined in model")
	}

	x := reflect.ValueOf(m.model).Elem().FieldByName("IsRemoved")
	if x.IsValid() {
		x.Set(reflect.ValueOf(true))
	}

	y := reflect.ValueOf(m.model).Elem().FieldByName("RemovedAt")
	if !y.IsValid() {
		return errors.New("RemovedAt not defined in model")
	}

	y.Set(reflect.ValueOf(time.Now()))
	_, err := m.collection.Upsert(bson.M{"_id": id.Interface()}, bson.M{"$set": m.model})
	return err
}

//GenQuery export mgo.Query for further usage
func (m *MgoFun) Q() *mgo.Query {
	return m.findQ()
}

//findQ conduct mgo.Query
func (m *MgoFun) findQ() *mgo.Query {
	var query *mgo.Query
	//do not query removed value
	rmQ := []interface{}{bson.M{"is_removed": bson.M{"$ne": true}}, bson.M{"IsRemoved": bson.M{"$ne": true}}}
	if m.Query != nil {
		if v, found := m.Query["$and"]; !found {
			m.Query["$and"] = rmQ
		} else {
			m.Query["$and"] = append(v.([]interface{}), rmQ...)
		}
	} else {
		m.Query = bson.M{"$and": rmQ}
	}

	query = m.collection.Find(m.Query)
	//sort
	if m.Sort != nil {
		query = query.Sort(m.Sort...)
	} else {
		query = query.Sort("-created_at", "-CreatedAt", "-UpdatedAt", "-updated_at")
	}
	//skip
	if m.Skip != 0 {
		query = query.Skip(m.Skip)
	}
	//limit
	if m.Limit != 0 {
		query = query.Limit(m.Limit)
	}
	return query
}

func (m *MgoFun) findByIdQ() *mgo.Query {
	var query *mgo.Query
	id := reflect.ValueOf(m.model).Elem().FieldByName("Id").Interface()
	query = m.collection.Find(bson.M{"_id": id})
	return query
}

//Count
func (m *MgoFun) Count() int64 {
	query := m.findQ()
	count, _ := query.Count()
	return int64(count)
}

//---------retrieve functions
// FindAll except removed, i is interface address
func (m *MgoFun) FindAll(i interface{}) error {
	query := m.findQ()
	err := query.All(i)
	return err
}

//Get will retrieve by _id
func (m *MgoFun) Get() error {
	query := m.findByIdQ()
	err := query.One(m.model)
	return err
}

//GetByQ get first one based on query, model will be updated
func (m *MgoFun) GetByQ() error {
	query := m.findQ()
	err := query.One(m.model)
	return err
}

//Select query and select columns
func (m *MgoFun) FindWithSelect(i interface{}, cols []string) error {
	sCols := bson.M{}
	for _, v := range cols {
		sCols[v] = 1
	}
	query := m.findQ().Select(sCols)
	err := query.All(i)
	return err
}

//Distinct
func (m *MgoFun) Distinct(key string, i interface{}) error {
	err := m.findQ().Distinct(key, i)
	return err
}
