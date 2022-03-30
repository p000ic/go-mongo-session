package mongo

import (
	"context"
	"sync"
	"time"

	"github.com/qiniu/qmgo"
)

type managerStore struct {
	ctx        context.Context
	session    *qmgo.Session
	source     *qmgo.Database
	client     *qmgo.Client
	authDBName string
	dbname     string
	cname      string
}

type store struct {
	sync.RWMutex
	s       *managerStore
	ctx     context.Context
	sid     string
	expired int64
	values  map[string]interface{}
}

// Data items stored in mongo
type sessionItem struct {
	SID       string    `bson:"sid"`
	Value     string    `bson:"value"`
	ExpiredAt time.Time `bson:"expired_at,omitempty"`
}

// close - close mongo session
func (ms *managerStore) close() {
	err := ms.client.Close(ms.ctx)
	if err != nil {
		return
	}
}

// cloneSession - cloneSession to Database
func (ms *managerStore) cloneSession() error {
	var err error
	ms.session, err = ms.client.Session()
	if err != nil {
		return err
	}
	return nil
}

// c - collection
func (ms *managerStore) c(clan string) *qmgo.Collection {
	return ms.source.Collection(clan)
}

// cHandler - collection handler
func (ms *managerStore) cHandler(clan string, handler func(c *qmgo.Collection)) {
	_, err := ms.client.Session()
	if err != nil {
		return
	}
	defer ms.session.EndSession(ms.ctx)
	handler(ms.source.Collection(clan))
}
