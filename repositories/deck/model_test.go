package deck

import "testing"

func TestCard_SuitName(t *testing.T) {
	card := &Card{
		Value: "10S",
	}
	if card.SuitName() != "SPADES" {
		t.Error("Expected 'SPADES' as suit name.")
	}

	card = &Card{
		Value: "AJ",
	}
	if card.SuitName() != "" {
		t.Error("Expected to get an empty suit name.")
	}
}

func TestCard_RankName(t *testing.T) {
	card := &Card{
		Value: "10S",
	}
	if card.RankName() != "10" {
		t.Error("Expected '10' as rank name.")
	}

	card = &Card{
		Value: "FC",
	}
	if card.RankName() != "" {
		t.Error("Expected to get an empty rank name.")
	}
}
