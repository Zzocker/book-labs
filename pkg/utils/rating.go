package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Zzocker/book-labs/pkg/datastore"
	"github.com/Zzocker/book-labs/pkg/errors"
)

type NetRating struct {
	Count int64
	Value float64
}

type Rating struct {
	ID          string
	UserID      string
	RecipientID string
	Value       int
	UpdateTime  int64
}

type ratingStore struct {
	store datastore.KVStore
}

func NewRatingStore(store datastore.KVStore) *ratingStore { // nolint:revive //require only for the methods
	return &ratingStore{store: store}
}

func (r *ratingStore) Set(ctx context.Context, userID, recipientID string, value int) error {
	const op = errors.Op("RatingStore.Set")
	rating := Rating{
		ID:          getRatingKey(userID, recipientID),
		UserID:      userID,
		RecipientID: recipientID,
		Value:       value,
		UpdateTime:  time.Now().Unix(),
	}
	raw, err := rating.toBytes(op)
	if err != nil {
		return err
	}

	err = r.store.Set(ctx, rating.ID, raw)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (r *ratingStore) Get(ctx context.Context, userID, recipientID string) (*Rating, error) {
	const op = errors.Op("RatingStore.Get")
	raw, err := r.store.Get(ctx, getRatingKey(userID, recipientID))
	if err != nil {
		return nil, errors.E(op, err)
	}
	var rating Rating
	err = json.Unmarshal(raw, &rating)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("failed to unmarshal rating: %w", err), errors.CodeInternal)
	}

	return &rating, nil
}

func (r *ratingStore) Delete(ctx context.Context, userID, recipientID string) error {
	const op = errors.Op("RatingStore.Delete")
	err := r.store.Del(ctx, getRatingKey(userID, recipientID))
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (r *Rating) New(value int) {
	r.ID = getRatingKey(r.UserID, r.RecipientID)
	r.Value = value
}

func getRatingKey(userID, recipientID string) string {
	return fmt.Sprintf("%s:%s", userID, recipientID)
}

func (r *Rating) toBytes(op errors.Op) ([]byte, error) {
	raw, err := json.Marshal(r)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("failed to marshal rating: %w", err), errors.CodeInternal)
	}

	return raw, nil
}
