package datastore

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type testData struct {
	ID string `bson:"_id"`
	F1 string `bson:"f1"`
	F2 string `bson:"f2"`
}

func TestMongoStore(t *testing.T) { // nolint:funlen //just a test
	is := assert.New(t)
	smartDS := NewMongoSmartDatastore(&MongoSmartDatastoreConfig{
		URL:        "localhost:27017",
		Username:   "root",
		Password:   "pw",
		Database:   "testDB",
		Collection: "testcollection",
	})
	id := primitive.NewObjectID().Hex()
	data := testData{
		ID: id,
		F1: "f1-field",
		F2: "f2-filed",
	}
	t.Run("Store", func(t *testing.T) {
		err := smartDS.Store(context.Background(), data)
		is.NoError(err)
	})

	t.Run("Store-Invalid", func(t *testing.T) {
		err := smartDS.Store(context.Background(), nil)
		is.Error(err)
		is.Equal(errors.CodeUnexpected, errors.ErrCode(err))
	})

	t.Run("Store-AlreadyExists", func(t *testing.T) {
		err := smartDS.Store(context.Background(), data)
		is.Error(err)
		is.Equal(errors.CodeAlreadyExists, errors.ErrCode(err))
	})

	t.Run("Get", func(t *testing.T) {
		raw, err := smartDS.Get(context.Background(), map[string]interface{}{
			"_id": data.ID,
		})
		is.NoError(err)
		var got testData
		err = bson.Unmarshal(raw, &got)
		is.NoError(err)
		is.Equal(data.ID, got.ID)
		is.Equal(data.F1, got.F1)
		is.Equal(data.F2, got.F2)
	})

	t.Run("NotFound-Get", func(t *testing.T) {
		raw, err := smartDS.Get(context.Background(), map[string]interface{}{
			"_id": primitive.NewObjectID().Hex(),
		})
		is.Error(err)
		is.Nil(raw)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Update", func(t *testing.T) {
		nData := testData{
			ID: data.ID,
			F1: "f1-field-v2",
			F2: "f2-field-v2",
		}
		filter := map[string]interface{}{
			"_id": data.ID,
		}
		err := smartDS.Update(context.Background(), filter, nData)
		is.NoError(err)

		raw, err := smartDS.Get(context.Background(), map[string]interface{}{
			"_id": data.ID,
		})
		is.NoError(err)
		var got testData
		err = bson.Unmarshal(raw, &got)
		is.NoError(err)
		is.Equal(nData.ID, got.ID)
		is.Equal(nData.F1, got.F1)
		is.Equal(nData.F2, got.F2)
	})

	t.Run("Update-NotFound", func(t *testing.T) {
		err := smartDS.Update(context.Background(), map[string]interface{}{
			"_id": primitive.NewObjectID().Hex(),
		}, data)
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Update-Invalid", func(t *testing.T) {
		err := smartDS.Update(context.Background(), map[string]interface{}{
			"_id": primitive.NewObjectID().Hex(),
		}, nil)
		is.Error(err)
		is.Equal(errors.CodeUnexpected, errors.ErrCode(err))
	})

	t.Run("UpdateMatching", func(t *testing.T) {
		err := smartDS.UpdateMatching(context.Background(), map[string]interface{}{
			"f1": "f1-field-v2",
		}, map[string]interface{}{
			"f2": "changed-by-f1",
		})
		is.NoError(err)

		raw, err := smartDS.Get(context.Background(), map[string]interface{}{
			"_id": data.ID,
		})
		is.NoError(err)
		var got testData
		err = bson.Unmarshal(raw, &got)
		is.NoError(err)
		is.Equal("changed-by-f1", got.F2)
	})

	t.Run("UpdateMatching-Invalid", func(t *testing.T) {
		err := smartDS.UpdateMatching(context.Background(), map[string]interface{}{
			"f1": "f1-field-v2",
		}, nil)
		is.Error(err)
		is.Equal(errors.CodeUnexpected, errors.ErrCode(err))
	})

	t.Run("Delete", func(t *testing.T) {
		err := smartDS.Delete(context.Background(), map[string]interface{}{
			"_id": data.ID,
		})
		is.NoError(err)
	})

	t.Run("Delete-NotFound", func(t *testing.T) {
		err := smartDS.Delete(context.Background(), map[string]interface{}{
			"_id": data.ID,
		})
		is.Error(err)
		is.Equal(errors.CodeNotFound, errors.ErrCode(err))
	})

	t.Run("Query", func(t *testing.T) {
		n := 8
		{
			wg := sync.WaitGroup{}
			wg.Add(n)
			// setup
			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					err := smartDS.Store(context.Background(), testData{
						ID: primitive.NewObjectID().Hex(),
						F1: "query-F1",
						F2: strconv.Itoa(i),
					})
					is.NoError(err)
				}(i)
			}

			wg.Wait()
		}
		raws, err := smartDS.Query(context.Background(), "f2", map[string]interface{}{
			"f1": "query-F1",
		}, 1, int64(n/2))
		is.NoError(err)
		is.Len(raws, n/2)
	})

	t.Run("DeleteMatching", func(t *testing.T) {
		err := smartDS.DeleteMatching(context.Background(), map[string]interface{}{
			"f1": "query-F1",
		})
		is.NoError(err)
	})
}
