package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/go-session/session/v3"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
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
	dbConfig := qmgo.Config{
		Uri:      cfg.URL,
		Database: cfg.Source,
		Coll:     cfg.Collection,
	}
	if cfg.Auth {
		dbConfig.Auth = &qmgo.Credential{
			AuthMechanism: cfg.AuthMechanism,
			Username:      cfg.Username,
			Password:      cfg.Password,
			AuthSource:    cfg.AuthSource,
		}
	}
	opts := options.ClientOptions{
		ClientOptions: cfg.ClientOptions,
	}
	var m db
	m.ctx = ctx
	m.client, err = qmgo.Open(ctx, &dbConfig, opts)
	if err != nil {
		return nil
	}
	m.authSource = cfg.AuthSource
	mgrStore := newManagerStore(&m)
	return mgrStore
}

func newManagerStore(db *db) *managerStore {
	err := db.cloneSession()
	if err != nil {
		return nil
	}
	defer db.endSession()
	t := true
	i := int32(60)
	_ = db.client.CreateIndexes(db.ctx, []options.IndexModel{{
		Key:          []string{"expired_at"},
		IndexOptions: &mongoOpts.IndexOptions{ExpireAfterSeconds: &i}},
	})
	_ = db.client.CreateIndexes(db.ctx, []options.IndexModel{{
		Key:          []string{"sid"},
		IndexOptions: &mongoOpts.IndexOptions{Unique: &t}},
	})
	return &managerStore{
		db:      db,
		ctx:     context.Background(),
		sid:     "",
		expired: 0,
		values:  nil,
		RWMutex: sync.RWMutex{},
	}
}

func (ms *managerStore) Check(_ context.Context, sid string) (bool, error) {
	val, err := ms.db.get(sid)
	if err != nil {
		return false, err
	}
	return val != "", nil
}

func (ms *managerStore) Create(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return newStore(ctx, ms.db, sid, expired, nil), nil
}

func (ms *managerStore) Update(ctx context.Context, sid string, expired int64) (session.Store, error) {
	value, err := ms.db.get(sid)
	if err != nil {
		// log.Printf("get::%s::%s", err, sid)
		return nil, err
	} else if len(value) == 0 {
		return newStore(ctx, ms.db, sid, expired, nil), nil
	}

	values, err := ms.db.parseValue(value)
	if err != nil {
		// log.Printf("parse-value::%s::%s", err, sid)
		return nil, err
	}

	err = ms.db.save(sid, values, expired)
	if err != nil {
		// log.Printf("save::%s::%s", err, sid)
		return nil, err
	}

	return newStore(ctx, ms.db, sid, expired, values), nil
}

func (ms *managerStore) Delete(_ context.Context, sid string) error {
	err := ms.db.cloneSession()
	if err != nil {
		return err
	}
	err = ms.db.delete(sid)
	if err != nil {
		if errors.Is(err, qmgo.ErrNoSuchDocuments) {
			err = nil
			return err
		}
		return err
	}
	return nil
}

func (ms *managerStore) Refresh(ctx context.Context, oldSid, sid string, expired int64) (session.Store, error) {
	value, err := ms.db.get(oldSid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, ms.db, sid, expired, nil), nil
	}

	values, err := ms.db.parseValue(value)
	if err != nil {
		return nil, err
	}

	err = ms.db.save(sid, values, expired)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return newStore(ctx, ms.db, sid, expired, values), nil
}

func (ms *managerStore) Close() error {
	ms.db.close()
	return nil
}

func newStore(ctx context.Context, db *db, sid string, expired int64, values map[string]interface{}) *store {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &store{
		db:      db,
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

func (s *store) Get(key string) (interface{}, bool) {
	s.Lock()
	defer s.Unlock()
	value, err := s.db.get(s.sid)
	if err != nil {
		// log.Printf("get::%s::%s", err, s.sid)
		return nil, false
	}
	values, err := s.db.parseValue(value)
	if err != nil {
		// log.Printf("parse-value::%s::%s", err, s.sid)
		return nil, false
	}
	val, ok := values[key]
	return val, ok
}

func (s *store) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.values[key] = value
}

func (s *store) Delete(key string) interface{} {
	s.Lock()
	defer s.Unlock()
	v, ok := s.values[key]
	if ok {
		delete(s.values, key)
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
	s.Lock()
	defer s.Unlock()
	err := s.db.save(s.sid, s.values, s.expired)
	if err != nil {
		return err
	}
	return nil
}
