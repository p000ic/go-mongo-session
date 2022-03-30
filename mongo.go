package mongo

import (
	"context"
	"encoding/json"
	"github.com/qiniu/qmgo/options"
	"log"
	"time"

	"github.com/go-session/session"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	mongoOpts "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_             session.ManagerStore = &managerStore{}
	_             session.Store        = &store{}
	jsonMarshal                        = json.Marshal
	jsonUnmarshal                      = json.Unmarshal
)

// NewStore Create an instance of a mongo store
func NewStore(cfg *Config) session.ManagerStore {
	var err error
	ctx := context.Background()
	dbConfig := qmgo.Config{Uri: cfg.URL}
	if cfg.Auth {
		dbConfig.Auth = &qmgo.Credential{
			AuthMechanism: cfg.AuthMechanism,
			Username:      cfg.Username,
			Password:      cfg.Password,
			AuthSource:    cfg.DB,
		}
	}
	var m managerStore
	m.ctx = ctx
	m.client, err = qmgo.NewClient(ctx, &dbConfig)
	if err != nil {
		return nil
	}
	m.dbname = cfg.DB
	m.cname = cfg.Collection
	m.authDBName = cfg.Source
	m.source = m.client.Database(cfg.DB)
	mgrStore := newManagerStore(&m, cfg)
	return mgrStore
}

// NewStoreWithSession Create an instance of a mongo store
func NewStoreWithSession(m *managerStore, cfg *Config) session.ManagerStore {
	return newManagerStore(m, cfg)
}

func newManagerStore(m *managerStore, cfg *Config) *managerStore {
	err := m.cloneSession()
	if err != nil {
		return nil
	}
	//t := true
	i := int32(60)
	_ = m.c(cfg.Collection).CreateIndexes(m.ctx, []options.IndexModel{{
		Key:          []string{"expired_at"},
		IndexOptions: &mongoOpts.IndexOptions{ExpireAfterSeconds: &i}},
	})
	_ = m.c(cfg.Collection).CreateIndexes(m.ctx, []options.IndexModel{{
		Key:          []string{"sid"},
		IndexOptions: &mongoOpts.IndexOptions{}},
	})
	return m
}

func (ms *managerStore) getValue(sid string) (value string, err error) {
	var item sessionItem
	err = ms.c(ms.cname).Find(ms.ctx, sessionItem{SID: sid}).One(item)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			value = ""
			err = nil
			return
		}
		value = ""
		return
	} else if item.ExpiredAt.Before(time.Now().UTC()) {
		value = ""
		err = nil
		return
	}
	value = item.Value
	err = nil
	return
}

func (ms *managerStore) parseValue(value string) (map[string]interface{}, error) {
	var values map[string]interface{}
	if len(value) > 0 {
		err := jsonUnmarshal([]byte(value), &values)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (ms *managerStore) Check(_ context.Context, sid string) (bool, error) {
	val, err := ms.getValue(sid)
	if err != nil {
		return false, err
	}
	return val != "", nil
}

func (ms *managerStore) Create(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return newStore(ctx, ms, sid, expired, nil), nil
}

func (ms *managerStore) Update(ctx context.Context, sid string, expired int64) (session.Store, error) {
	err := ms.cloneSession()
	if err != nil {
		return nil, err
	}

	value, err := ms.getValue(sid)
	if err != nil {
		return nil, err
	} else if value == "" {
		log.Printf("%s", value)
		return newStore(ctx, ms, sid, expired, nil), nil
	}

	err = ms.c(ms.cname).UpdateOne(ms.ctx, sessionItem{SID: sid}, bson.M{
		"$set": sessionItem{
			ExpiredAt: time.Now().UTC().Add(time.Duration(expired) * time.Second),
		},
	})
	if err != nil {
		return nil, err
	}

	values, err := ms.parseValue(value)
	if err != nil {
		return nil, err
	}
	log.Printf("%+v", values)
	return newStore(ctx, ms, sid, expired, values), nil
}

func (ms *managerStore) Delete(_ context.Context, sid string) error {
	err := ms.cloneSession()
	if err != nil {
		return err
	}
	err = ms.c(ms.cname).Remove(ms.ctx, sessionItem{SID: sid})
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			err = nil
			return err
		}
		return err
	}
	return nil
}

func (ms *managerStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error) {
	value, err := ms.getValue(oldsid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, ms, sid, expired, nil), nil
	}

	err = ms.cloneSession()
	if err != nil {
		return nil, err
	}
	_, err = ms.c(ms.cname).InsertOne(ms.ctx, sessionItem{
		SID:       sid,
		Value:     value,
		ExpiredAt: time.Now().UTC().Add(time.Duration(expired) * time.Second),
	})
	if err != nil {
		return nil, err
	}
	err = ms.c(ms.cname).Remove(ms.ctx, sessionItem{SID: oldsid})
	if err != nil {
		return nil, err
	}

	values, err := ms.parseValue(value)
	if err != nil {
		return nil, err
	}

	return newStore(ctx, ms, sid, expired, values), nil
}

func (ms *managerStore) Close() error {
	ms.close()
	return nil
}

func newStore(ctx context.Context, s *managerStore, sid string, expired int64, values map[string]interface{}) *store {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &store{
		s:       s,
		ctx:     ctx,
		sid:     sid,
		expired: expired,
		values:  values,
	}
}

func (s *store) Context() context.Context {
	return s.ctx
}

func (s *store) SessionID() string {
	return s.sid
}

func (s *store) Set(key string, value interface{}) {
	s.Lock()
	s.values[key] = value
	s.Unlock()
}

func (s *store) Get(key string) (interface{}, bool) {
	s.RLock()
	val, ok := s.values[key]
	s.RUnlock()
	return val, ok
}

func (s *store) Delete(key string) interface{} {
	s.RLock()
	v, ok := s.values[key]
	s.RUnlock()
	if ok {
		s.Lock()
		delete(s.values, key)
		s.Unlock()
	}
	return v
}

func (s *store) Flush() error {
	s.Lock()
	s.values = make(map[string]interface{})
	s.Unlock()
	return s.Save()
}

func (s *store) Save() error {
	var value string
	s.RLock()
	if len(s.values) > 0 {
		buf, err := jsonMarshal(s.values)
		if err != nil {
			s.RUnlock()
			return err
		}
		value = string(buf)
	}
	s.RUnlock()
	_, err := s.s.c(s.s.cname).Upsert(s.ctx, sessionItem{SID: s.sid}, sessionItem{
		SID:       s.sid,
		Value:     value,
		ExpiredAt: time.Now().UTC().Add(time.Duration(s.expired) * time.Second),
	})
	if err != nil {
		return err
	}
	return nil
}
