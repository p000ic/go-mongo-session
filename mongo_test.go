package mongo

import (
	"context"
	"github.com/go-session/session"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	url    = "mongodb://127.0.0.1:27017"
	dbName = "test"
	cName  = "session"
)

func TestStore(t *testing.T) {
	cfg := NewConfig(url, dbName, cName, "", "", "")
	mStore := NewStore(cfg)
	defer mStore.Close()
	Convey("Test mongo storage operation", t, func() {
		sid := "test_mongo_store"
		store, err := mStore.Create(context.Background(), sid, 300)
		So(err, ShouldBeNil)

		err = store.Flush()
		So(err, ShouldBeNil)

		foo, ok := store.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo, ShouldBeNil)

		store.Set("foo", "bar")
		store.Set("foo2", "bar2")
		err = store.Save()
		So(err, ShouldBeNil)

		foo, ok = store.Get("foo")
		So(ok, ShouldBeTrue)
		So(foo, ShouldEqual, "bar")

		foo = store.Delete("foo")
		So(foo, ShouldEqual, "bar")

		foo, ok = store.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo, ShouldBeNil)

		foo2, ok := store.Get("foo2")
		So(ok, ShouldBeTrue)
		So(foo2, ShouldEqual, "bar2")

		err = store.Flush()
		So(err, ShouldBeNil)

		foo2, ok = store.Get("foo2")
		So(ok, ShouldBeFalse)
		So(foo2, ShouldBeNil)
	})
}

func TestManagerStore(t *testing.T) {
	cfg := NewConfig(url, dbName, cName, "", "", "")
	mStore := NewStore(cfg)
	defer func(mStore session.ManagerStore) {
		err := mStore.Close()
		if err != nil {

		}
	}(mStore)

	Convey("Test mongo-based storage management operations", t, func() {
		sid := "test_manager_store"
		store, err := mStore.Create(context.Background(), sid, 10)
		So(store, ShouldNotBeNil)
		So(err, ShouldBeNil)

		err = store.Flush()
		So(err, ShouldBeNil)

		store.Set("foo", "bar")
		err = store.Save()
		So(err, ShouldBeNil)

		store, err = mStore.Update(context.Background(), sid, 10)
		So(store, ShouldNotBeNil)
		So(err, ShouldBeNil)

		foo, ok := store.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo, ShouldBeNil)

		newsid := "test_manager_store2"
		store, err = mStore.Refresh(context.Background(), sid, newsid, 10)
		So(store, ShouldNotBeNil)
		So(err, ShouldBeNil)

		foo, ok = store.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo, ShouldBeNil)

		exists, err := mStore.Check(context.Background(), sid)
		So(exists, ShouldBeFalse)
		So(err, ShouldBeNil)

		err = mStore.Delete(context.Background(), newsid)
		So(err, ShouldBeNil)

		exists, err = mStore.Check(context.Background(), newsid)
		So(exists, ShouldBeFalse)
		So(err, ShouldBeNil)
	})
}
