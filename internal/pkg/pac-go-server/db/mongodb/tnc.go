package mongodb

import (
	"context"
	"fmt"

	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetTermsAndConditionsByUserID gets the terms and conditions status for a user.
func (db *MongoDB) GetTermsAndConditionsByUserID(id string) (*models.TermsAndConditions, error) {
	terms := []models.TermsAndConditions{}

	if id == "" {
		return nil, fmt.Errorf("user id is required")
	}

	filter := bson.D{{Key: "user_id", Value: id}}

	collection := db.Database.Collection("tnc")
	ctx, cancel := context.WithTimeout(context.Background(), dbContextTimeout)
	defer cancel()
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting terms and coditions entries: %w", err)
	}
	defer cur.Close(ctx)

	if err = cur.All(context.TODO(), &terms); err != nil {
		return nil, fmt.Errorf("error fetching terms and coditions entries: %w", err)
	}

	if len(terms) == 0 {
		return nil, fmt.Errorf("no terms and coditions entries found for user id: %s, %w", id, mongo.ErrNoDocuments)
	}

	return &terms[0], nil
}

// AcceptTermsAndConditions updates the terms and conditions for a user. Creates a new entry if one does not exist or updates the existing one.
func (db *MongoDB) AcceptTermsAndConditions(terms *models.TermsAndConditions) error {
	collection := db.Database.Collection("tnc")
	ctx, cancel := context.WithTimeout(context.Background(), dbContextTimeout)
	defer cancel()
	filter := bson.D{{"user_id", terms.UserID}}
	update := bson.D{{"$set", bson.D{{"user_id", terms.UserID}, {"accepted", terms.Accepted}, {"accepted_at", terms.AcceptedAt}}}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error updating terms and coditions: %w", err)
	}

	return nil
}
