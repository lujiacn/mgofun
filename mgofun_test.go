package mgofun

import (
	"fmt"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	BaseModel `bson:",inline"`
	Name      string `bson:"name,omitempty"`
	Age       int    `bson:"age,omitempty"`
}

var (
	dbName = "mgofun_test"
	dial   = "localhost"
)

func TestFuncsSave(t *testing.T) {
	s, err := mgo.Dial(dial)
	if err != nil {
		panic("Cannot connect to database")
	}

	user := new(User)
	user.Name = "Tom"
	user.Age = 10
	user.Id = bson.NewObjectId()

	user2 := new(User)
	user2.Name = "Jack"
	user2.Age = 20
	user2.Id = bson.NewObjectId()

	//conduct op
	op := NewMgoFun(s, dbName, user)
	op.Save()
	op = NewMgoFun(s, dbName, user2)
	op.Save()

}

func TestFindAll(t *testing.T) {
	s, err := mgo.Dial(dial)
	if err != nil {
		panic("Cannot connect to database")
	}

	user := new(User)
	op := NewMgoFun(s, dbName, user)
	fmt.Println(op.Count())

	// for FindAll
	var users []*User
	op.FindAll(&users)
	op.Query = bson.M{"name": "Jia"}
	op.GetByQ()
	fmt.Println(user)

	// query, limit, skip
	op.Sort = "-updated_at"
	op.Limit = 10
	op.FindAll(&users)
	fmt.Println(users[0].UpdatedAt)
}
