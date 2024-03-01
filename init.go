package mongo

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/qiniu/qmgo"
)

type db struct {
	ctx        context.Context
	session    *qmgo.Session
	database   *qmgo.Database
	client     *qmgo.Client
	authSource string
	source     string
	collection string
}

type managerStore struct {
	sync.RWMutex
	db      *db
	ctx     context.Context
	sid     string
	expired int64
	values  map[string]interface{}
}

type store struct {
	sync.RWMutex
	db      *db
	ctx     context.Context
	sid     string
	expired int64
	values  map[string]interface{}
}

// Data items stored in mongo
type sessionItem struct {
	SID       string                 `bson:"sid"`
	Value     map[string]interface{} `bson:"value"`
	ExpiredAt time.Time              `bson:"expired_at,omitempty"`
}

// close - close mongo session
func (x *db) close() {
	err := x.client.Close(x.ctx)
	if err != nil {
		return
	}
}

// cloneSession - cloneSession to Database
func (x *db) cloneSession() (err error) {
	x.session, err = x.client.Session()
	if err != nil {
		return err
	}
	return nil
}

// endSession - endSession
func (x *db) endSession() {
	x.session.EndSession(x.ctx)
}

// c - collection
func (x *db) c(clan string) *qmgo.Collection {
	return x.database.Collection(clan)
}

// cHandler - collection handler
func (x *db) cHandler(clan string, handler func(c *qmgo.Collection)) {
	_, err := x.client.Session()
	if err != nil {
		return
	}
	defer x.session.EndSession(x.ctx)
	handler(x.database.Collection(clan))
}

// get -
func (x *db) get(sid string) (value string, err error) {
	var item sessionItem

	x.cHandler(x.collection, func(c *qmgo.Collection) {
		e := x.c(x.collection).Find(x.ctx, bson.M{"sid": sid}).One(&item)
		if e != nil {
			err = e
			return
		}
		err = nil
	})

	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			value = ""
			err = errors.New("sid does not exist [" + sid + "]")
			return
		}
		value = ""
		return
	} else if item.ExpiredAt.Before(time.Now().UTC()) {
		value = ""
		err = errors.New("sid expired [" + sid + "]")
		return
	}

	marshal, err := jsonMarshal(item.Value)
	if err != nil {
		return "", err
	}

	value = string(marshal)
	err = nil

	return
}

// parseValue -
func (x *db) parseValue(value string) (map[string]interface{}, error) {
	var values map[string]interface{}

	if len(value) > 0 {
		err := jsonUnmarshal([]byte(value), &values)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

// save -
func (x *db) save(sid string, values map[string]interface{}, expired int64) (err error) {
	x.cHandler(x.collection, func(c *qmgo.Collection) {
		_, e := x.c(x.collection).Upsert(x.ctx, bson.M{"sid": sid}, sessionItem{
			SID:       sid,
			Value:     values,
			ExpiredAt: time.Now().UTC().Add(time.Duration(expired) * time.Second),
		})
		if e != nil {
			err = e
			return
		}
		err = nil
	})

	if err != nil {
		return err
	}

	return nil
}

// delete -
func (x *db) delete(sid string) (err error) {
	x.cHandler(x.collection, func(c *qmgo.Collection) {
		e := x.c(x.collection).Remove(x.ctx, bson.M{"sid": sid})
		if e != nil {
			err = e
			return
		}
		err = nil
	})

	if err != nil {
		return err
	}

	return nil
}
