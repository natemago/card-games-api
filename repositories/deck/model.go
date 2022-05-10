package deck

import (
	"time"
)

type Deck struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Shuffled  bool
	Remaining int
	Cards     []*Card
}

type Card struct {
	DeckID string `gorm:"primaryKey"`
	Value  string `gorm:"primaryKey"`
	Drawn  bool
	Idx    int
}

func (c *Card) SuitName() string {
	if c.Value == "" {
		return ""
	}
	suit := c.Value[len(c.Value)-1:]
	suitName, ok := SuitsNames[suit]
	if !ok {
		return ""
	}
	return suitName
}

func (c *Card) RankName() string {
	if c.Value == "" {
		return ""
	}
	rank := c.Value[:len(c.Value)-1]
	rankName, ok := RanksNames[rank]
	if !ok {
		return ""
	}
	return rankName
}
