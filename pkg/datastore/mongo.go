package datastore

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type MongoRichStoreConfig struct {
	URL        string
	Username   string
	Password   string
	Database   string
	Collection string
}

func NewMongoRichStore(cfg *MongoRichStoreConfig) RichStore {
	address := fmt.Sprintf("mongodb://%s:%s@%s", cfg.Username, cfg.Password, cfg.URL)
	client, err := mongo.NewClient(
		options.Client().ApplyURI(address),
	)
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	return &mongoRichStore{
		collection: client.Database(cfg.Database).Collection(cfg.Collection),
	}
}

type mongoRichStore struct {
	collection *mongo.Collection
}

func (m *mongoRichStore) Put(ctx context.Context, in interface{}) error {
	const op = errors.Op("MongoStore.Put")

	_, err := m.collection.InsertOne(ctx, in)
	if isDuplicate(err) {
		return errors.E(op, fmt.Errorf("already exists: %w", err), errors.CodeAlreadyExists)
	} else if err != nil {
		return errors.E(op, fmt.Errorf("failed to store: %w", err), errors.CodeInternal)
	}

	return nil
}

func (m *mongoRichStore) Get(ctx context.Context, filter map[string]interface{}) ([]byte, error) {
	const op = errors.Op("MongoStore.Get")

	resp := m.collection.FindOne(ctx, filter)

	if resp.Err() == mongo.ErrNoDocuments {
		return nil, errors.E(op, fmt.Errorf("not found: %w", resp.Err()), errors.CodeNotFound)
	} else if resp.Err() != nil {
		return nil, errors.E(op, fmt.Errorf("failed to get: %w", resp.Err()), errors.CodeInternal)
	}
	raw, err := resp.DecodeBytes()
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("failed to decode: %w", err), errors.CodeInternal)
	}

	return raw, nil
}

func (m *mongoRichStore) Update(ctx context.Context, filter map[string]interface{}, in interface{}) error {
	const op = errors.Op("MongoStore.Update")

	resp, err := m.collection.UpdateOne(ctx, filter, bson.M{"$set": in})
	if err != nil {
		return errors.E(op, fmt.Errorf("failed to update: %w", err), errors.CodeInternal)
	}
	if resp.MatchedCount != 1 {
		return errors.E(op, fmt.Errorf("not found: %w", err), errors.CodeNotFound)
	}

	return nil
}

func (m *mongoRichStore) Delete(ctx context.Context, filter map[string]interface{}) error {
	const op = errors.Op("MongoStore.Delete")

	resp, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.E(op, fmt.Errorf("failed to deleted: %w", err), errors.CodeInternal)
	}

	if resp.DeletedCount != 1 {
		return errors.E(op, fmt.Errorf("not found: %w", err), errors.CodeNotFound)
	}

	return nil
}

func (m *mongoRichStore) Query(ctx context.Context, query map[string]interface{}) ([][]byte, error) {
	const op = errors.Op("MongoStore.Query")

	itr, err := m.collection.Find(ctx, query)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("failed to create query interator: %w", err), errors.CodeInternal)
	}
	defer itr.Close(ctx)

	var raws [][]byte
	for itr.Next(ctx) {
		raws = append(raws, []byte(string(itr.Current)))
	}

	return raws, nil
}

func isDuplicate(err error) bool {
	if merr, ok := err.(mongo.WriteException); ok {
		for _, e := range merr.WriteErrors {
			if e.Code == 11000 { // nolint:gomnd //status code from monog
				return true
			}
		}
	}

	return false
}
