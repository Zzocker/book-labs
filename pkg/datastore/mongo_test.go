package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type testData struct {
	ID string `bson:"_id"`
	F1 string `bson:"f1"`
	F2 string `bson:"f2"`
}

// nolint:funlen //this is a test
func TestMongoRichStore(t *testing.T) {
	is := assert.New(t)

	store := NewMongoRichStore(&MongoRichStoreConfig{
		URL:        "localhost:27017",
		Username:   "root",
		Password:   "pw",
		Database:   "test",
		Collection: "test-collection",
	})

	data1 := testData{
		ID: "id1",
		F1: "id1-f1",
		F2: "id1-f2",
	}

	ctx := context.Background()
	t.Run("Put", func(t *testing.T) {
		err := store.Put(ctx, data1)
		is.NoError(err)
	})

	t.Run("Put-Duplicate", func(t *testing.T) {
		err := store.Put(ctx, data1)
		is.Error(err)
		is.Equal(errors.CodeAlreadyExists, errors.ErrCode(err))
	})

	t.Run("Get", func(t *testing.T) {
		raw, err := store.Get(ctx, map[string]interface{}{
			"_id": data1.ID,
		})
		is.NoError(err)

		var data testData
		err = bson.Unmarshal(raw, &data)
		is.NoError(err)
		is.Equal(data1, data)

		_, err = store.Get(ctx, map[string]interface{}{
			"_id": "not-found",
		})
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Update", func(t *testing.T) {
		err := store.Update(ctx, map[string]interface{}{
			"_id": data1.ID,
		}, map[string]interface{}{
			"f1": "f1-new",
			"f2": "f2-new",
		})

		is.NoError(err)

		raw, err := store.Get(ctx, map[string]interface{}{
			"_id": data1.ID,
		})
		is.NoError(err)

		var data testData
		err = bson.Unmarshal(raw, &data)
		is.NoError(err)

		is.Equal("f1-new", data.F1)
		is.Equal("f2-new", data.F2)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete(ctx, map[string]interface{}{
			"_id": data1.ID,
		})
		is.NoError(err)

		err = store.Delete(ctx, map[string]interface{}{
			"_id": data1.ID,
		})
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Query", func(t *testing.T) {
		n := 8
		for i := 0; i < n; i++ {
			err := store.Put(ctx, testData{
				ID: fmt.Sprintf("%d", i),
				F1: "query-f1",
				F2: fmt.Sprintf("id-%d-f2", i),
			})
			is.NoError(err)
		}

		raws, err := store.Query(ctx, map[string]interface{}{
			"f1": "query-f1",
		})
		is.NoError(err)
		is.Len(raws, n)

		for _, raw := range raws {
			var data testData
			err = bson.Unmarshal(raw, &data)
			is.NoError(err)

			err = store.Delete(ctx, map[string]interface{}{
				"_id": data.ID,
			})
			is.NoError(err)
		}
	})
}
