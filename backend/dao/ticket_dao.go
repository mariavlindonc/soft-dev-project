package db

import (
	"backend/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketDAOImpl struct {
	db *gorm.DB
}

func NewTicketDAO(db *gorm.DB) *TicketDAOImpl {
	return &TicketDAOImpl{db: db}
}

func (d *TicketDAOImpl) Create(ticket *domain.Ticket) error {
	return d.db.Create(ticket).Error
}

func (d *TicketDAOImpl) FindByID(id uint) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := d.db.Preload("User").Preload("Event").
		Take(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (d *TicketDAOImpl) FindByUserID(userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	// Includes tickets the user bought and tickets transferred to them.
	err := d.db.Where("user_id = ? OR transferred_to_id = ?", userID, userID).
		Preload("Event").Order("purchased_at DESC").Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// validTicketCondition matches tickets that still grant entry to the event:
// a transferred ticket is valid for its new owner, so it counts like an active one.
const validTicketCondition = "event_id = ? AND status IN ('active', 'transferred')"

func (d *TicketDAOImpl) FindActiveByEvent(eventID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := d.db.Where(validTicketCondition, eventID).
		Preload("User").Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (d *TicketDAOImpl) CountActiveByEvent(eventID uint) (int, error) {
	var count int64
	err := d.db.Model(&domain.Ticket{}).
		Where(validTicketCondition, eventID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (d *TicketDAOImpl) CancelByEvent(eventID uint) error {
	return d.db.Model(&domain.Ticket{}).
		Where(validTicketCondition, eventID).
		Update("status", "cancelled").Error
}

// Save persists only the ticket's own columns. Associations are omitted
// because FindByID preloads User/Event, and GORM would otherwise overwrite
// UserID with the preloaded User.ID (breaking transfers).
func (d *TicketDAOImpl) Save(ticket *domain.Ticket) error {
	return d.db.Omit(clause.Associations).Save(ticket).Error
}

func (d *TicketDAOImpl) WithTransaction(fn func(TxContext) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

var _ TicketDAO = (*TicketDAOImpl)(nil)
