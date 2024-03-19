***REMOVED***

***REMOVED***
***REMOVED***
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/qiniu/qmgo"
***REMOVED***

type db struct {
	ctx        context.Context
	session    *qmgo.Session
	database   *qmgo.Database
	client     *qmgo.QmgoClient
	authSource string
	source     string
	collection string
***REMOVED***

type managerStore struct {
	sync.RWMutex
	db      *db
	ctx     context.Context
	sid     string
	expired int64
	values  map[string]interface{***REMOVED***
***REMOVED***

type store struct {
	sync.RWMutex
	db      *db
	ctx     context.Context
	sid     string
	expired int64
	values  map[string]interface{***REMOVED***
***REMOVED***

// Data items stored in mongo
type sessionItem struct {
	SID       string                 `bson:"sid"`
	Value     map[string]interface{***REMOVED*** `bson:"value"`
	ExpiredAt time.Time              `bson:"expired_at,omitempty"`
***REMOVED***

// close - close mongo session
func (x *db***REMOVED*** close(***REMOVED*** {
	err := x.client.Close(x.ctx***REMOVED***
***REMOVED***
***REMOVED***
	***REMOVED***
***REMOVED***

// cloneSession - cloneSession to Database
func (x *db***REMOVED*** cloneSession(***REMOVED*** (err error***REMOVED*** {
	x.session, err = x.client.Session(***REMOVED***
***REMOVED***
***REMOVED*** err
	***REMOVED***
	return nil
***REMOVED***

// endSession - endSession
func (x *db***REMOVED*** endSession(***REMOVED*** {
	x.session.EndSession(x.ctx***REMOVED***
***REMOVED***

// c - collection
func (x *db***REMOVED*** c(clan string***REMOVED*** *qmgo.Collection {
	return x.database.Collection(clan***REMOVED***
***REMOVED***

// cHandler - collection handler
func (x *db***REMOVED*** cHandler(clan string, handler func(c *qmgo.Collection***REMOVED******REMOVED*** {
	_, err := x.client.Session(***REMOVED***
***REMOVED***
***REMOVED***
	***REMOVED***
	defer x.session.EndSession(x.ctx***REMOVED***
	handler(x.database.Collection(clan***REMOVED******REMOVED***
***REMOVED***

// get -
func (x *db***REMOVED*** get(sid string***REMOVED*** (value string, err error***REMOVED*** {
	var item sessionItem
	_ctx, cancel := context.WithTimeout(x.ctx, 5*time.Second***REMOVED***
	defer cancel(***REMOVED***
	// x.cHandler(x.collection, func(c *qmgo.Collection***REMOVED*** {
	err = x.client.Find(_ctx, bson.M{"sid": sid***REMOVED******REMOVED***.One(&item***REMOVED***
	//   if e != nil {
	//     err = e
	//     return
	//   ***REMOVED***
	//   err = nil
	// ***REMOVED******REMOVED***

***REMOVED***
		if errors.Is(err, qmgo.ErrNoSuchDocuments***REMOVED*** {
			value = ""
			err = errors.New("sid does not exist [" + sid + "]"***REMOVED***
	***REMOVED***
		***REMOVED***
		value = ""
***REMOVED***
	***REMOVED*** else if item.ExpiredAt.Before(time.Now(***REMOVED***.UTC(***REMOVED******REMOVED*** {
		value = ""
		err = errors.New("sid expired [" + sid + "]"***REMOVED***
***REMOVED***
	***REMOVED***

	marshal, err := jsonMarshal(item.Value***REMOVED***
***REMOVED***
***REMOVED*** "", err
	***REMOVED***

	value = string(marshal***REMOVED***
	err = nil

	return
***REMOVED***

// parseValue -
func (x *db***REMOVED*** parseValue(value string***REMOVED*** (map[string]interface{***REMOVED***, error***REMOVED*** {
	var values map[string]interface{***REMOVED***

	if len(value***REMOVED*** > 0 {
		err := jsonUnmarshal([]byte(value***REMOVED***, &values***REMOVED***
	***REMOVED***
	***REMOVED*** nil, err
		***REMOVED***
	***REMOVED***

	return values, nil
***REMOVED***

// save -
func (x *db***REMOVED*** save(sid string, values map[string]interface{***REMOVED***, expired int64***REMOVED*** (err error***REMOVED*** {
	_ctx, cancel := context.WithTimeout(x.ctx, 5*time.Second***REMOVED***
	defer cancel(***REMOVED***
	// x.cHandler(x.collection, func(c *qmgo.Collection***REMOVED*** {
	_, err = x.client.Upsert(_ctx, bson.M{"sid": sid***REMOVED***, sessionItem{
		SID:       sid,
		Value:     values,
		ExpiredAt: time.Now(***REMOVED***.UTC(***REMOVED***.Add(time.Duration(expired***REMOVED*** * time.Second***REMOVED***,
	***REMOVED******REMOVED***
	//   if e != nil {
	//     err = e
	//     return
	//   ***REMOVED***
	//   err = nil
	// ***REMOVED******REMOVED***

***REMOVED***
***REMOVED*** err
	***REMOVED***

	return nil
***REMOVED***

// delete -
func (x *db***REMOVED*** delete(sid string***REMOVED*** (err error***REMOVED*** {
	_ctx, cancel := context.WithTimeout(x.ctx, 5*time.Second***REMOVED***
	defer cancel(***REMOVED***
	// x.cHandler(x.collection, func(c *qmgo.Collection***REMOVED*** {
	err = x.client.Remove(_ctx, bson.M{"sid": sid***REMOVED******REMOVED***
	//   if e != nil {
	//     err = e
	//     return
	//   ***REMOVED***
	//   err = nil
	// ***REMOVED******REMOVED***

***REMOVED***
***REMOVED*** err
	***REMOVED***

	return nil
***REMOVED***
