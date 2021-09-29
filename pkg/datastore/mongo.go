package datastore

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type mongoStore struct {
	collection *mongo.Collection
}

type MongoSmartDatastoreConfig struct {
	Username   string
	Password   string
	URL        string
	Database   string
	Collection string
}

const mongoPingTimeoutS = 5 // sec

func NewMongoSmartDatastore(cfg *MongoSmartDatastoreConfig) SmartDataStore {
	adrs := fmt.Sprintf("mongodb://%s:%s@%s", cfg.Username, cfg.Password, cfg.URL)
	client, err := mongo.NewClient(options.Client().ApplyURI(adrs))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), mongoPingTimeoutS*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return &mongoStore{
		collection: client.Database(cfg.Database).Collection(cfg.Collection),
	}
}

func (m *mongoStore) Store(ctx context.Context, in interface{}) error {
	const op = errors.Op("MongoSmartDatastore.Store")
	_, err := m.collection.InsertOne(ctx, in)
	if m.isDuplicate(err) {
		return errors.E(
			op,
			fmt.Errorf("duplicate entry : %w", err),
			errors.CodeAlreadyExists,
			errors.SeverityDebug,
		)
	} else if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed to store key : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return nil
}

func (m *mongoStore) Get(ctx context.Context, filter map[string]interface{}) ([]byte, error) {
	const op = errors.Op("MongoSmartDatastore.Get")
	res := m.collection.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return nil, errors.E(
			op,
			fmt.Errorf("entry not found : %w", res.Err()),
			errors.CodeNotFound,
			errors.SeverityDebug,
		)
	}
	raw, err := res.DecodeBytes()
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to decode document : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return raw, nil
}

func (m *mongoStore) Update(ctx context.Context, filter map[string]interface{}, in interface{}) error {
	const op = errors.Op("MongoSmartDatastore.Update")
	resp, err := m.collection.UpdateOne(ctx, filter, bson.M{"$set": in})
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed update document : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}
	if resp.MatchedCount != 1 {
		return errors.E(
			op,
			fmt.Errorf("entry not found"),
			errors.CodeNotFound,
			errors.SeverityDebug,
		)
	}

	return nil
}

func (m *mongoStore) UpdateMatching(ctx context.Context, query map[string]interface{}, in interface{}) error {
	const op = errors.Op("MongoSmartDatastore.UpdateMatching")
	_, err := m.collection.UpdateMany(ctx, query, bson.M{"$set": in})
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed update document : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return nil
}

func (m *mongoStore) Delete(ctx context.Context, filter map[string]interface{}) error {
	const op = errors.Op("MongoSmartDatastore.Delete")
	resp, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed delete document : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}
	if resp.DeletedCount != 1 {
		return errors.E(
			op,
			fmt.Errorf("document not found"),
			errors.CodeNotFound,
			errors.SeverityDebug,
		)
	}

	return nil
}

func (m *mongoStore) Query(ctx context.Context, sortingKey string, query map[string]interface{}, pageNumber, perPage int64) ([][]byte, error) {
	const op = errors.Op("MongoSmartDatastore.Query")
	skip := (pageNumber - 1) * perPage
	if skip < 0 {
		skip = 0
	}
	opts := options.FindOptions{
		Limit: &perPage,
		Skip:  &skip,
		Sort:  bson.M{sortingKey: 1},
	}
	cur, err := m.collection.Find(ctx, query, &opts)
	if err != nil {
		return nil, errors.E(
			op,
			fmt.Errorf("failed to make query : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}
	defer cur.Close(ctx)
	raws := make([][]byte, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		raws = append(raws, cur.Current)
	}

	return raws, nil
}

func (m *mongoStore) DeleteMatching(ctx context.Context, query map[string]interface{}) error {
	const op = errors.Op("MongoSmartDatastore.DeleteMatching")
	_, err := m.collection.DeleteMany(ctx, query)
	if err != nil {
		return errors.E(
			op,
			fmt.Errorf("failed delete document : %w", err),
			errors.CodeUnexpected,
			errors.SeverityError,
		)
	}

	return nil
}

func (m *mongoStore) isDuplicate(err error) bool {
	if mErr, ok := err.(mongo.WriteException); ok {
		for _, e := range mErr.WriteErrors {
			if e.Code == 11000 { // nolint:gomnd //status code from mongo
				return true
			}
		}
	}

	return false
}
