package domain

import (
	"time"

	"gorm.io/gorm"
)

type SalePhase string

const (
	PhaseNotYetOpen SalePhase = "not_yet_open"
	PhasePresale    SalePhase = "presale"
	PhasePublic     SalePhase = "public"
	PhaseNoPresale  SalePhase = "no_presale"
)

type Event struct {
	ID               uint           `gorm:"primaryKey;autoIncrement;type:int unsigned" json:"id"`
	Title            string         `gorm:"type:varchar(200);not null" json:"title"`
	Description      *string        `gorm:"type:text" json:"description,omitempty"`
	ImageURL         *string        `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	Category         *string        `gorm:"type:varchar(100);index:idx_events_category" json:"category,omitempty"`
	Location         *string        `gorm:"type:varchar(300)" json:"location,omitempty"`
	EventDate        time.Time      `gorm:"not null;index:idx_events_event_date;index:idx_events_status_date,priority:2" json:"event_date"`
	DurationMinutes  int            `gorm:"not null;default:0" json:"duration_minutes"`
	Capacity         int            `gorm:"not null;default:0" json:"capacity"`
	TicketsSold      int            `gorm:"not null;default:0" json:"tickets_sold"`
	Price            float64        `gorm:"type:decimal(10,2);not null;default:0.00" json:"price"`
	Status           string         `gorm:"type:enum('active','presale','sold_out','cancelled');not null;default:'active';index:idx_events_status_date,priority:1" json:"status"`
	PresaleActive    bool           `gorm:"not null;default:0" json:"presale_active"`
	PresaleCode      *string        `gorm:"type:varchar(100)" json:"-"`
	PresaleStartDate *time.Time    `gorm:"type:datetime" json:"presale_start_date,omitempty"`
	GeneralSaleDate  *time.Time    `gorm:"type:datetime" json:"general_sale_date,omitempty"`
	CreatedByID      uint           `gorm:"type:int unsigned;not null;index:idx_events_created_by" json:"created_by_id"`
	CreatedAt        time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index:idx_events_deleted_at" json:"deleted_at,omitempty"`

	CreatedBy User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
}

func (e *Event) CurrentSalePhase(now time.Time) SalePhase {
	if !e.PresaleActive || e.PresaleStartDate == nil || e.GeneralSaleDate == nil {
		return PhaseNoPresale
	}
	if now.Before(*e.PresaleStartDate) {
		return PhaseNotYetOpen
	}
	if now.Before(*e.GeneralSaleDate) {
		return PhasePresale
	}
	return PhasePublic
}
