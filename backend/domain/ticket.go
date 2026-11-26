package domain

import (
	"time"
)

type Ticket struct {
	ID              uint       `gorm:"primaryKey;autoIncrement;type:int unsigned" json:"id"`
	UserID          uint       `gorm:"type:int unsigned;not null;index:idx_tickets_user_id" json:"user_id"`
	EventID         uint       `gorm:"type:int unsigned;not null;index:idx_tickets_event_id;index:idx_tickets_event_status,priority:1" json:"event_id"`
	Status          string     `gorm:"type:enum('active','cancelled','transferred');not null;default:'active';index:idx_tickets_event_status,priority:2" json:"status"`
	PurchasePrice   float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"purchase_price"`
	PurchasedAt     time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"purchased_at"`
	CancelledAt     *time.Time `gorm:"type:datetime" json:"cancelled_at,omitempty"`
	TransferredAt   *time.Time `gorm:"type:datetime" json:"transferred_at,omitempty"`
	TransferredToID *uint      `gorm:"type:int unsigned;index:idx_tickets_transferred_to" json:"transferred_to_id,omitempty"`
	CreatedAt       time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	User          User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"user"`
	Event         Event `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"event"`
	TransferredTo *User `gorm:"foreignKey:TransferredToID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"transferred_to,omitempty"`
}

// HolderID returns who currently holds the ticket. UserID is always the
// original buyer; after a transfer the ticket belongs to TransferredToID.
func (t *Ticket) HolderID() uint {
	if t.Status == "transferred" && t.TransferredToID != nil {
		return *t.TransferredToID
	}
	return t.UserID
}

// TransferredOutBy reports whether userID is the original buyer who gave
// this ticket away (they still see it, but can no longer act on it).
func (t *Ticket) TransferredOutBy(userID uint) bool {
	return t.Status == "transferred" && t.UserID == userID && t.HolderID() != userID
}
