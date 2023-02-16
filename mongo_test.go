***REMOVED***

***REMOVED***
***REMOVED***
	"errors"
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
***REMOVED***
		store0, err := mStore.Create(context.Background(***REMOVED***, sid, 300***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo0, ok := store0.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo0, ShouldBeNil***REMOVED***

		store0.Set("foo", "bar"***REMOVED***
		store0.Set("foo2", "bar2"***REMOVED***
		err = store0.Save(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo1, ok := store0.Get("foo"***REMOVED***
		So(ok, ShouldBeTrue***REMOVED***
		So(foo1, ShouldEqual, "bar"***REMOVED***

		foo := store0.Delete("foo"***REMOVED***
		So(foo, ShouldEqual, "bar"***REMOVED***

		err = store0.Save(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok = store0.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo, ShouldBeNil***REMOVED***

		foo2, ok := store0.Get("foo2"***REMOVED***
		So(ok, ShouldBeTrue***REMOVED***
		So(foo2, ShouldEqual, "bar2"***REMOVED***

		err = store0.Flush(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo2, ok = store0.Get("foo2"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo2, ShouldBeNil***REMOVED***
	***REMOVED******REMOVED***
***REMOVED***

func TestManagerStore(t *testing.T***REMOVED*** {
	cfg := NewConfig(url, dbName, cName, "", "", ""***REMOVED***
	mStore := NewStore(cfg***REMOVED***
	defer mStore.Close(***REMOVED***
	Convey("Test mongo-based storage management operations", t, func(***REMOVED*** {
***REMOVED***
		store0, err := mStore.Create(context.Background(***REMOVED***, sid, 20***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		store0.Set("foo", "bar"***REMOVED***
		err = store0.Save(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok := store0.Get("foo"***REMOVED***
		So(ok, ShouldBeTrue***REMOVED***
		So(foo, ShouldEqual, "bar"***REMOVED***

		store0, err = mStore.Update(context.Background(***REMOVED***, sid, 10***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		err = store0.Flush(***REMOVED***
		So(err, ShouldBeNil***REMOVED***

***REMOVED***
		store2, err := mStore.Refresh(context.Background(***REMOVED***, sid, newSID, 10***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		foo, ok = store2.Get("foo"***REMOVED***
		So(ok, ShouldBeFalse***REMOVED***
		So(foo, ShouldBeNil***REMOVED***

		exists, err := mStore.Check(context.Background(***REMOVED***, sid***REMOVED***
		So(exists, ShouldBeTrue***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		err = mStore.Delete(context.Background(***REMOVED***, newSID***REMOVED***
		So(err, ShouldBeNil***REMOVED***

		exists, err = mStore.Check(context.Background(***REMOVED***, newSID***REMOVED***
		So(exists, ShouldBeFalse***REMOVED***
		So(err, ShouldResemble, errors.New("sid does not exist"***REMOVED******REMOVED***
	***REMOVED******REMOVED***
***REMOVED***
