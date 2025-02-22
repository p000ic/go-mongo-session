package mongo

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type db struct {
	ctx        context.Context
	authSource string
	client     *mongo.Client
	session    *mongo.Session
	database   *mongo.Database
	collection *mongo.Collection
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
	err := x.client.Disconnect(x.ctx)
	if err != nil {
		return
	}
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

// cloneSession - cloneSession to Database
func (x *db) cloneSession() (s *mongo.Session, err error) {
	s, err = x.client.StartSession()
	if err != nil {
		return nil, err
	}
	return
}

// get -
func (x *db) get(sid string) (value string, err error) {
	// TODO: Implement session txn logic
	var item sessionItem
	_ctx, cancel := context.WithTimeout(x.ctx, 5*time.Second)
	defer cancel()
	err = x.collection.FindOne(_ctx, bson.M{"sid": sid}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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

// save -
func (x *db) save(sid string, values map[string]interface{}, expired int64) (err error) {
	// TODO: Implement session txn logic
	_ctx, cancel := context.WithTimeout(x.ctx, 5*time.Second)
	defer cancel()
	_, err = x.collection.UpdateOne(_ctx, bson.M{"sid": sid},
		bson.M{"$set": &sessionItem{
			SID:       sid,
			Value:     values,
			ExpiredAt: time.Now().UTC().Add(time.Duration(expired) * time.Second),
		}}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

// delete -
func (x *db) delete(sid string) (err error) {
	// TODO: Implement session txn logic
	_ctx, cancel := context.WithTimeout(x.ctx, 5*time.Second)
	defer cancel()
	_, err = x.collection.DeleteOne(_ctx, bson.M{"sid": sid})
	if err != nil {
		return err
	}
	return nil
}
