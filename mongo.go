package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-session/session/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mongoOpts "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	_             session.ManagerStore = &managerStore{}
	_             session.Store        = &store{}
	jsonMarshal                        = json.Marshal
	jsonUnmarshal                      = json.Unmarshal
	minPoolSize                        = uint64(5)
	maxPoolSize                        = uint64(99)
	connTimeout                        = time.Second * 10 // 10 sec
	maxIdleTime                        = time.Second * 1  // 1 sec
)

// NewStore Create an instance of a mongo store
func NewStore(cfg *Config) session.ManagerStore {
	var err error
	var m db
	m.ctx = context.Background()
	opts := mongoOpts.Client()
	opts.SetConnectTimeout(connTimeout)
	opts.SetMaxConnIdleTime(maxIdleTime)
	opts.SetMinPoolSize(minPoolSize)
	opts.SetMaxPoolSize(maxPoolSize)
	if cfg.Auth {
		opts.Auth = &mongoOpts.Credential{
			AuthMechanism: cfg.AuthMechanism,
			Username:      cfg.Username,
			Password:      cfg.Password,
			AuthSource:    cfg.AuthSource,
		}
	}
	var startedCommands sync.Map
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			startedCommands.Store(evt.RequestID, evt.Command)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			startedCommands.Delete(evt.RequestID)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			if cmd, ok := startedCommands.Load(evt.RequestID); ok {
				log.Printf("cmd: %v failure-resp: %v", cmd, evt.Failure)
				startedCommands.Delete(evt.RequestID)
			}
		},
	}
	opts.ApplyURI(cfg.URL).SetMonitor(cmdMonitor)
	m.client, err = mongo.Connect(opts)
	if err != nil {
		return nil
	}
	m.authSource = cfg.AuthSource
	m.database = m.client.Database(cfg.Database)
	m.collection = m.database.Collection(cfg.Collection)
	mgrStore := newManagerStore(&m)
	return mgrStore
}

func newManagerStore(m *db) *managerStore {
	s, err := m.cloneSession()
	if err != nil {
		return nil
	}
	defer m.endSession(s)
	i := int32(0)
	var indexModel []mongo.IndexModel
	indexModel = append(indexModel, mongo.IndexModel{
		Keys: bson.D{{"expired_at", 1}},
		Options: mongoOpts.Index().
			SetName("_" + m.collection.Name() + "_expired_at_").
			SetExpireAfterSeconds(i),
	})
	indexModel = append(indexModel, mongo.IndexModel{
		Keys: bson.D{{"sid", 1}},
		Options: mongoOpts.Index().
			SetName("_" + m.collection.Name() + "_sid_").
			SetUnique(true),
	})
	_, err = m.collection.Indexes().CreateMany(m.ctx, indexModel)
	if err != nil {
		return nil
	}

	return &managerStore{
		db:      m,
		ctx:     context.Background(),
		sid:     "",
		expired: int64(i),
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
	err := ms.db.delete(sid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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
