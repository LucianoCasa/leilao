package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
)

func TestAutoCloseAuction(t *testing.T) {

	_ = os.Setenv("AUCTION_DURATION", "2s")

	ctx := context.Background()

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}

	repo := NewAuctionRepository(db)

	auctionEntity, ierr := auction_entity.CreateAuction(
		"AutoClose Test",
		"electronics",
		"Leilão de teste para fechamento automático",
		auction_entity.ProductCondition(0),
	)
	if ierr != nil {
		t.Fatalf("failed to create auction entity: %v", ierr)
	}

	if _, err := repo.CreateAuction(ctx, auctionEntity); err != nil {
		t.Fatalf("failed to persist auction: %v", err)
	}

	time.Sleep(3 * time.Second)

	got, ierr := repo.FindAuctionById(ctx, auctionEntity.Id)
	if ierr != nil {
		t.Fatalf("failed to fetch auction: %v", ierr)
	}

	if got.Status != auction_entity.Completed {
		t.Fatalf("expected status Completed, got %v", got.Status)
	}
}
