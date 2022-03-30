***REMOVED***

***REMOVED***
***REMOVED***
	"encoding/json"
	"github.com/qiniu/qmgo/options"
	"log"
	"time"

	"github.com/go-session/session"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
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
			AuthSource:    cfg.DB,
		***REMOVED***
	***REMOVED***
	var m managerStore
	m.ctx = ctx
	m.client, err = qmgo.NewClient(ctx, &dbConfig***REMOVED***
***REMOVED***
***REMOVED*** nil
	***REMOVED***
	m.dbname = cfg.DB
	m.cname = cfg.Collection
	m.authDBName = cfg.Source
	m.source = m.client.Database(cfg.DB***REMOVED***
	mgrStore := newManagerStore(&m, cfg***REMOVED***
	return mgrStore
***REMOVED***

// NewStoreWithSession Create an instance of a mongo store
func NewStoreWithSession(m *managerStore, cfg *Config***REMOVED*** session.ManagerStore {
	return newManagerStore(m, cfg***REMOVED***
***REMOVED***

func newManagerStore(m *managerStore, cfg *Config***REMOVED*** *managerStore {
	err := m.cloneSession(***REMOVED***
***REMOVED***
***REMOVED*** nil
	***REMOVED***
	//t := true
	i := int32(60***REMOVED***
	_ = m.c(cfg.Collection***REMOVED***.CreateIndexes(m.ctx, []options.IndexModel{{
		Key:          []string{"expired_at"***REMOVED***,
		IndexOptions: &mongoOpts.IndexOptions{ExpireAfterSeconds: &i***REMOVED******REMOVED***,
	***REMOVED******REMOVED***
	_ = m.c(cfg.Collection***REMOVED***.CreateIndexes(m.ctx, []options.IndexModel{{
		Key:          []string{"sid"***REMOVED***,
		IndexOptions: &mongoOpts.IndexOptions{***REMOVED******REMOVED***,
	***REMOVED******REMOVED***
	return m
***REMOVED***

func (ms *managerStore***REMOVED*** getValue(sid string***REMOVED*** (value string, err error***REMOVED*** {
	var item sessionItem
	err = ms.c(ms.cname***REMOVED***.Find(ms.ctx, sessionItem{SID: sid***REMOVED******REMOVED***.One(item***REMOVED***
***REMOVED***
		if err == qmgo.ErrNoSuchDocuments {
			value = ""
			err = nil
	***REMOVED***
		***REMOVED***
		value = ""
***REMOVED***
	***REMOVED*** else if item.ExpiredAt.Before(time.Now(***REMOVED***.UTC(***REMOVED******REMOVED*** {
		value = ""
		err = nil
***REMOVED***
	***REMOVED***
	value = item.Value
	err = nil
	return
***REMOVED***

func (ms *managerStore***REMOVED*** parseValue(value string***REMOVED*** (map[string]interface{***REMOVED***, error***REMOVED*** {
	var values map[string]interface{***REMOVED***
	if len(value***REMOVED*** > 0 {
		err := jsonUnmarshal([]byte(value***REMOVED***, &values***REMOVED***
	***REMOVED***
	***REMOVED*** nil, err
		***REMOVED***
	***REMOVED***

	return values, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Check(_ context.Context, sid string***REMOVED*** (bool, error***REMOVED*** {
	val, err := ms.getValue(sid***REMOVED***
***REMOVED***
***REMOVED*** false, err
	***REMOVED***
	return val != "", nil
***REMOVED***

func (ms *managerStore***REMOVED*** Create(ctx context.Context, sid string, expired int64***REMOVED*** (session.Store, error***REMOVED*** {
	return newStore(ctx, ms, sid, expired, nil***REMOVED***, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Update(ctx context.Context, sid string, expired int64***REMOVED*** (session.Store, error***REMOVED*** {
	err := ms.cloneSession(***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	value, err := ms.getValue(sid***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED*** else if value == "" {
		log.Printf("%s", value***REMOVED***
***REMOVED*** newStore(ctx, ms, sid, expired, nil***REMOVED***, nil
	***REMOVED***

	err = ms.c(ms.cname***REMOVED***.UpdateOne(ms.ctx, sessionItem{SID: sid***REMOVED***, bson.M{
		"$set": sessionItem{
			ExpiredAt: time.Now(***REMOVED***.UTC(***REMOVED***.Add(time.Duration(expired***REMOVED*** * time.Second***REMOVED***,
		***REMOVED***,
	***REMOVED******REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	values, err := ms.parseValue(value***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***
	log.Printf("%+v", values***REMOVED***
	return newStore(ctx, ms, sid, expired, values***REMOVED***, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Delete(_ context.Context, sid string***REMOVED*** error {
	err := ms.cloneSession(***REMOVED***
***REMOVED***
***REMOVED*** err
	***REMOVED***
	err = ms.c(ms.cname***REMOVED***.Remove(ms.ctx, sessionItem{SID: sid***REMOVED******REMOVED***
***REMOVED***
		if err == qmgo.ErrNoSuchDocuments {
			err = nil
	***REMOVED*** err
		***REMOVED***
***REMOVED*** err
	***REMOVED***
	return nil
***REMOVED***

func (ms *managerStore***REMOVED*** Refresh(ctx context.Context, oldsid, sid string, expired int64***REMOVED*** (session.Store, error***REMOVED*** {
	value, err := ms.getValue(oldsid***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED*** else if value == "" {
***REMOVED*** newStore(ctx, ms, sid, expired, nil***REMOVED***, nil
	***REMOVED***

	err = ms.cloneSession(***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***
	_, err = ms.c(ms.cname***REMOVED***.InsertOne(ms.ctx, sessionItem{
		SID:       sid,
		Value:     value,
		ExpiredAt: time.Now(***REMOVED***.UTC(***REMOVED***.Add(time.Duration(expired***REMOVED*** * time.Second***REMOVED***,
	***REMOVED******REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***
	err = ms.c(ms.cname***REMOVED***.Remove(ms.ctx, sessionItem{SID: oldsid***REMOVED******REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	values, err := ms.parseValue(value***REMOVED***
***REMOVED***
***REMOVED*** nil, err
	***REMOVED***

	return newStore(ctx, ms, sid, expired, values***REMOVED***, nil
***REMOVED***

func (ms *managerStore***REMOVED*** Close(***REMOVED*** error {
	ms.close(***REMOVED***
	return nil
***REMOVED***

func newStore(ctx context.Context, s *managerStore, sid string, expired int64, values map[string]interface{***REMOVED******REMOVED*** *store {
	if values == nil {
		values = make(map[string]interface{***REMOVED******REMOVED***
	***REMOVED***
	return &store{
		s:       s,
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

func (s *store***REMOVED*** Set(key string, value interface{***REMOVED******REMOVED*** {
	s.Lock(***REMOVED***
	s.values[key] = value
	s.Unlock(***REMOVED***
***REMOVED***

func (s *store***REMOVED*** Get(key string***REMOVED*** (interface{***REMOVED***, bool***REMOVED*** {
	s.RLock(***REMOVED***
	val, ok := s.values[key]
	s.RUnlock(***REMOVED***
	return val, ok
***REMOVED***

func (s *store***REMOVED*** Delete(key string***REMOVED*** interface{***REMOVED*** {
	s.RLock(***REMOVED***
	v, ok := s.values[key]
	s.RUnlock(***REMOVED***
	if ok {
		s.Lock(***REMOVED***
		delete(s.values, key***REMOVED***
		s.Unlock(***REMOVED***
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
	var value string
	s.RLock(***REMOVED***
	if len(s.values***REMOVED*** > 0 {
		buf, err := jsonMarshal(s.values***REMOVED***
	***REMOVED***
			s.RUnlock(***REMOVED***
	***REMOVED*** err
		***REMOVED***
		value = string(buf***REMOVED***
	***REMOVED***
	s.RUnlock(***REMOVED***
	_, err := s.s.c(s.s.cname***REMOVED***.Upsert(s.ctx, sessionItem{SID: s.sid***REMOVED***, sessionItem{
		SID:       s.sid,
		Value:     value,
		ExpiredAt: time.Now(***REMOVED***.UTC(***REMOVED***.Add(time.Duration(s.expired***REMOVED*** * time.Second***REMOVED***,
	***REMOVED******REMOVED***
***REMOVED***
***REMOVED*** err
	***REMOVED***
	return nil
***REMOVED***
