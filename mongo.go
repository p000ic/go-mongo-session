***REMOVED***

***REMOVED***
***REMOVED***
	"encoding/json"
	"sync"

	"github.com/go-session/session/v3"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	mongoOpts "go.mongodb.org/mongo-driver/mongo/options"
***REMOVED***

var (
	_             session.ManagerStore = &managerStore{***REMOVED***
	_             session.Store        = &store{***REMOVED***
	jsonMarshal                        = json.Marshal
	jsonUnmarshal                      = json.Unmarshal
***REMOVED***

// NewStore Create an instance of a mongo store
func NewStore(cfg *Config***REMOVED*** session.ManagerStore {
	var err error
	ctx := context.Background(***REMOVED***
	dbConfig := qmgo.Config{Uri: cfg.URL***REMOVED***
	if cfg.Auth {
		dbConfig.Auth = &qmgo.Credential{
			AuthMechanism: cfg.AuthMechanism,
			Username:      cfg.Username,
			Password:      cfg.Password,
			AuthSource:    cfg.AuthSource,
		***REMOVED***
	***REMOVED***
	opts := options.ClientOptions{
		ClientOptions: cfg.ClientOptions,
	***REMOVED***
	var m db
	m.ctx = ctx
	m.client, err = qmgo.NewClient(ctx, &dbConfig, opts***REMOVED***
***REMOVED***
***REMOVED*** nil
	***REMOVED***
	m.source = cfg.Source
	m.collection = cfg.Collection
	m.authSource = cfg.AuthSource
	m.database = m.client.Database(cfg.Source***REMOVED***
	mgrStore := newManagerStore(&m, cfg***REMOVED***
	return mgrStore
***REMOVED***

func newManagerStore(db *db, cfg *Config***REMOVED*** *managerStore {
	err := db.cloneSession(***REMOVED***
***REMOVED***
***REMOVED*** nil
	***REMOVED***
	defer db.endSession(***REMOVED***
	t := true
	i := int32(60***REMOVED***
	_ = db.c(cfg.Collection***REMOVED***.CreateIndexes(db.ctx, []options.IndexModel{{
		Key:          []string{"expired_at"***REMOVED***,
		IndexOptions: &mongoOpts.IndexOptions{ExpireAfterSeconds: &i***REMOVED******REMOVED***,
	***REMOVED******REMOVED***
	_ = db.c(cfg.Collection***REMOVED***.CreateIndexes(db.ctx, []options.IndexModel{{
		Key:          []string{"sid"***REMOVED***,
		IndexOptions: &mongoOpts.IndexOptions{Unique: &t***REMOVED******REMOVED***,
	***REMOVED******REMOVED***
	return &managerStore{
		db:      db,
		ctx:     context.Background(***REMOVED***,
		sid:     "",
		expired: 0,
		values:  nil,
		RWMutex: sync.RWMutex{***REMOVED***,
	***REMOVED***
***REMOVED***

func (ms *managerStore***REMOVED*** Check(_ context.Context, sid string***REMOVED*** (bool, error***REMOVED*** {
	val, err := ms.db.get(sid***REMOVED***
***REMOVED***
***REMOVED*** false, err
	***REMOVED***
	return val != "", nil
***REMOVED***

func (ms *managerStore***REMOVED*** Create(ctx context.Context, sid string, expired int64***REMOVED*** (session.Store, error***REMOVED*** {
	return newStore(ctx, ms.db, sid, expired, nil***REMOVED***, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Update(ctx context.Context, sid string, expired int64***REMOVED*** (session.Store, error***REMOVED*** {
	value, err := ms.db.get(sid***REMOVED***
***REMOVED***
		// log.Printf("get::%s::%s", err, sid***REMOVED***
***REMOVED*** nil, err
	***REMOVED*** else if len(value***REMOVED*** == 0 {
***REMOVED*** newStore(ctx, ms.db, sid, expired, nil***REMOVED***, nil
	***REMOVED***

	values, err := ms.db.parseValue(value***REMOVED***
***REMOVED***
		// log.Printf("parse-value::%s::%s", err, sid***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	err = ms.db.save(sid, values, expired***REMOVED***
***REMOVED***
		// log.Printf("save::%s::%s", err, sid***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	return newStore(ctx, ms.db, sid, expired, values***REMOVED***, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Delete(_ context.Context, sid string***REMOVED*** error {
	err := ms.db.cloneSession(***REMOVED***
***REMOVED***
***REMOVED*** err
	***REMOVED***
	err = ms.db.delete(sid***REMOVED***
***REMOVED***
		if err == qmgo.ErrNoSuchDocuments {
			err = nil
	***REMOVED*** err
		***REMOVED***
***REMOVED*** err
	***REMOVED***
	return nil
***REMOVED***

func (ms *managerStore***REMOVED*** Refresh(ctx context.Context, oldSid, sid string, expired int64***REMOVED*** (session.Store, error***REMOVED*** {
	value, err := ms.db.get(oldSid***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED*** else if value == "" {
***REMOVED*** newStore(ctx, ms.db, sid, expired, nil***REMOVED***, nil
	***REMOVED***

	values, err := ms.db.parseValue(value***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	err = ms.db.save(sid, values, expired***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	return newStore(ctx, ms.db, sid, expired, values***REMOVED***, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Close(***REMOVED*** error {
	ms.db.close(***REMOVED***
	return nil
***REMOVED***

func newStore(ctx context.Context, db *db, sid string, expired int64, values map[string]interface{***REMOVED******REMOVED*** *store {
	if values == nil {
		values = make(map[string]interface{***REMOVED******REMOVED***
	***REMOVED***
	return &store{
		db:      db,
		ctx:     ctx,
		sid:     sid,
		expired: expired,
		values:  values,
	***REMOVED***
***REMOVED***

func (s *store***REMOVED*** Context(***REMOVED*** context.Context {
	return s.ctx
***REMOVED***

func (s *store***REMOVED*** SessionID(***REMOVED*** string {
	return s.sid
***REMOVED***

func (s *store***REMOVED*** Get(key string***REMOVED*** (interface{***REMOVED***, bool***REMOVED*** {
	s.Lock(***REMOVED***
	defer s.Unlock(***REMOVED***
	value, err := s.db.get(s.sid***REMOVED***
***REMOVED***
		// log.Printf("get::%s::%s", err, s.sid***REMOVED***
***REMOVED*** nil, false
	***REMOVED***
	values, err := s.db.parseValue(value***REMOVED***
***REMOVED***
		// log.Printf("parse-value::%s::%s", err, s.sid***REMOVED***
***REMOVED*** nil, false
	***REMOVED***
	val, ok := values[key]
	return val, ok
***REMOVED***

func (s *store***REMOVED*** Set(key string, value interface{***REMOVED******REMOVED*** {
	s.Lock(***REMOVED***
	defer s.Unlock(***REMOVED***
	s.values[key] = value
***REMOVED***

func (s *store***REMOVED*** Delete(key string***REMOVED*** interface{***REMOVED*** {
	s.Lock(***REMOVED***
	defer s.Unlock(***REMOVED***
	v, ok := s.values[key]
	if ok {
		delete(s.values, key***REMOVED***
	***REMOVED***
	return v
***REMOVED***

func (s *store***REMOVED*** Flush(***REMOVED*** error {
	s.Lock(***REMOVED***
	s.values = make(map[string]interface{***REMOVED******REMOVED***
	s.Unlock(***REMOVED***
	return s.Save(***REMOVED***
***REMOVED***

func (s *store***REMOVED*** Save(***REMOVED*** error {
	s.Lock(***REMOVED***
	defer s.Unlock(***REMOVED***
	err := s.db.save(s.sid, s.values, s.expired***REMOVED***
***REMOVED***
***REMOVED*** err
	***REMOVED***
	return nil
***REMOVED***
