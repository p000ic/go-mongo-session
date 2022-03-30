***REMOVED***

***REMOVED***
***REMOVED***
	"github.com/go-session/session"
***REMOVED***

***REMOVED***
***REMOVED***

***REMOVED***
***REMOVED***
***REMOVED***
***REMOVED***
***REMOVED***

func TestStore(t *testing.T***REMOVED*** {
	cfg := NewConfig(url, dbName, cName, "", "", ""***REMOVED***
	mStore := NewStore(cfg***REMOVED***
	defer mStore.Close(***REMOVED***
	Convey("Test mongo storage operation", t, func(***REMOVED*** {
		sid := "test_mongo_store"
		store, err := mStore.Create(context.Background(***REMOVED***, sid, 300***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		err = store.Flush(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok := store.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo, ShouldBeNil***REMOVED***

		store.Set("foo", "bar"***REMOVED***
		store.Set("foo2", "bar2"***REMOVED***
		err = store.Save(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok = store.Get("foo"***REMOVED***
		So(ok, ShouldBeTrue***REMOVED***
		So(foo, ShouldEqual, "bar"***REMOVED***

		foo = store.Delete("foo"***REMOVED***
		So(foo, ShouldEqual, "bar"***REMOVED***

		foo, ok = store.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo, ShouldBeNil***REMOVED***

		foo2, ok := store.Get("foo2"***REMOVED***
		So(ok, ShouldBeTrue***REMOVED***
		So(foo2, ShouldEqual, "bar2"***REMOVED***

		err = store.Flush(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo2, ok = store.Get("foo2"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo2, ShouldBeNil***REMOVED***
	***REMOVED******REMOVED***
***REMOVED***

func TestManagerStore(t *testing.T***REMOVED*** {
	cfg := NewConfig(url, dbName, cName, "", "", ""***REMOVED***
	mStore := NewStore(cfg***REMOVED***
	defer func(mStore session.ManagerStore***REMOVED*** {
		err := mStore.Close(***REMOVED***
	***REMOVED***

		***REMOVED***
	***REMOVED***(mStore***REMOVED***

	Convey("Test mongo-based storage management operations", t, func(***REMOVED*** {
		sid := "test_manager_store"
		store, err := mStore.Create(context.Background(***REMOVED***, sid, 10***REMOVED***
		So(store, ShouldNotBeNil***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		err = store.Flush(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		store.Set("foo", "bar"***REMOVED***
		err = store.Save(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		store, err = mStore.Update(context.Background(***REMOVED***, sid, 10***REMOVED***
		So(store, ShouldNotBeNil***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok := store.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo, ShouldBeNil***REMOVED***

		newsid := "test_manager_store2"
		store, err = mStore.Refresh(context.Background(***REMOVED***, sid, newsid, 10***REMOVED***
		So(store, ShouldNotBeNil***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok = store.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo, ShouldBeNil***REMOVED***

		exists, err := mStore.Check(context.Background(***REMOVED***, sid***REMOVED***
		So(exists, ShouldBeFalse***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		err = mStore.Delete(context.Background(***REMOVED***, newsid***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		exists, err = mStore.Check(context.Background(***REMOVED***, newsid***REMOVED***
		So(exists, ShouldBeFalse***REMOVED***
		So(err, ShouldBeNil***REMOVED***
	***REMOVED******REMOVED***
***REMOVED***
