***REMOVED***

***REMOVED***
***REMOVED***
	"sync"
	"time"

	"github.com/qiniu/qmgo"
***REMOVED***

type managerStore struct {
	ctx        context.Context
	session    *qmgo.Session
	source     *qmgo.Database
	client     *qmgo.Client
	authDBName string
	dbname     string
	cname      string
***REMOVED***

type store struct {
	sync.RWMutex
	s       *managerStore
	ctx     context.Context
	sid     string
	expired int64
	values  map[string]interface{***REMOVED***
***REMOVED***

// Data items stored in mongo
type sessionItem struct {
	SID       string    `bson:"sid"`
	Value     string    `bson:"value"`
	ExpiredAt time.Time `bson:"expired_at,omitempty"`
***REMOVED***

// close - close mongo session
func (ms *managerStore***REMOVED*** close(***REMOVED*** {
	err := ms.client.Close(ms.ctx***REMOVED***
***REMOVED***
***REMOVED***
	***REMOVED***
***REMOVED***

// cloneSession - cloneSession to Database
func (ms *managerStore***REMOVED*** cloneSession(***REMOVED*** error {
	var err error
	ms.session, err = ms.client.Session(***REMOVED***
***REMOVED***
***REMOVED*** err
	***REMOVED***
	return nil
***REMOVED***

// c - collection
func (ms *managerStore***REMOVED*** c(clan string***REMOVED*** *qmgo.Collection {
	return ms.source.Collection(clan***REMOVED***
***REMOVED***

// cHandler - collection handler
func (ms *managerStore***REMOVED*** cHandler(clan string, handler func(c *qmgo.Collection***REMOVED******REMOVED*** {
	_, err := ms.client.Session(***REMOVED***
***REMOVED***
***REMOVED***
	***REMOVED***
	defer ms.session.EndSession(ms.ctx***REMOVED***
	handler(ms.source.Collection(clan***REMOVED******REMOVED***
***REMOVED***
