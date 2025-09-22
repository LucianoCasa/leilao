package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) (string, *internal_error.InternalError) {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return "", internal_error.NewInternalServerError("Error trying to insert auction")
	}

	ar.scheduleAutoClose(auctionEntity.Id, auctionEntity.Timestamp)

	return auctionEntity.Id, nil
}

func getAuctionDuration() time.Duration {
	val := os.Getenv("AUCTION_DURATION")
	d, err := time.ParseDuration(val)
	if err != nil || d <= 0 {
		return 5 * time.Minute
	}
	return d
}

func (ar *AuctionRepository) scheduleAutoClose(auctionID string, startedAt time.Time) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered in scheduleAutoClose", nil)
			}
		}()

		duration := getAuctionDuration()
		sleepFor := startedAt.Add(duration).Sub(time.Now())
		if sleepFor > 0 {
			time.Sleep(sleepFor)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"_id": auctionID, "status": auction_entity.Active}
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

		res, err := ar.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Error trying to auto-close auction", err)
			return
		}
		if res.ModifiedCount > 0 {
			logger.Info("Auction closed automatically")
		}
	}()
}
