package mongo

import (
	"context"
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
		sid := "test_mongo_store0"
		store0, err := mStore.Create(context.Background(), sid, 300)
		So(err, ShouldBeNil)

		foo0, ok := store0.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo0, ShouldBeNil)

		store0.Set("foo", "bar")
		store0.Set("foo2", "bar2")
		err = store0.Save()
		So(err, ShouldBeNil)

		foo1, ok := store0.Get("foo")
		So(ok, ShouldBeTrue)
		So(foo1, ShouldEqual, "bar")

		foo := store0.Delete("foo")
		So(foo, ShouldEqual, "bar")

		foo, ok = store0.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo, ShouldBeNil)

		foo2, ok := store0.Get("foo2")
		So(ok, ShouldBeTrue)
		So(foo2, ShouldEqual, "bar2")

		err = store0.Flush()
		So(err, ShouldBeNil)

		foo2, ok = store0.Get("foo2")
		So(ok, ShouldBeFalse)
		So(foo2, ShouldBeNil)
	})
}

func TestManagerStore(t *testing.T) {
	cfg := NewConfig(url, dbName, cName, "", "", "")
	mStore := NewStore(cfg)
	defer mStore.Close()
	Convey("Test mongo-based storage management operations", t, func() {
		sid := "test_manager_store1"
		store0, err := mStore.Create(context.Background(), sid, 100)
		So(store0, ShouldNotBeNil)
		So(err, ShouldBeNil)

		store0.Set("foo", "bar")
		err = store0.Save()
		So(err, ShouldBeNil)

		foo, ok := store0.Get("foo")
		So(ok, ShouldBeTrue)
		So(foo, ShouldEqual, "bar")

		store0, err = mStore.Update(context.Background(), sid, 10)
		So(store0, ShouldNotBeNil)
		So(err, ShouldBeNil)

		err = store0.Flush()
		So(err, ShouldBeNil)

		newSID := "test_manager_store2"
		store2, err := mStore.Refresh(context.Background(), sid, newSID, 10)
		So(store2, ShouldNotBeNil)
		So(err, ShouldBeNil)

		foo, ok = store2.Get("foo")
		So(ok, ShouldBeFalse)
		So(foo, ShouldBeNil)

		exists, err := mStore.Check(context.Background(), sid)
		So(exists, ShouldBeFalse)
		So(err, ShouldBeNil)

		err = mStore.Delete(context.Background(), newSID)
		So(err, ShouldBeNil)

		exists, err = mStore.Check(context.Background(), newSID)
		So(exists, ShouldBeFalse)
		So(err, ShouldBeNil)
	})
}
