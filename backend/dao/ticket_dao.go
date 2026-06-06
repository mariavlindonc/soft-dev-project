package db

import (
	"backend/domain"

	"gorm.io/gorm"
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
	err := d.db.Where("user_id = ?", userID).
		Preload("Event").Order("purchased_at DESC").Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (d *TicketDAOImpl) FindActiveByEvent(eventID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := d.db.Where("event_id = ? AND status = 'active'", eventID).
		Preload("User").Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (d *TicketDAOImpl) CountActiveByEvent(eventID uint) (int, error) {
	var count int64
	err := d.db.Model(&domain.Ticket{}).
		Where("event_id = ? AND status = 'active'", eventID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (d *TicketDAOImpl) CancelByEvent(eventID uint) error {
	return d.db.Model(&domain.Ticket{}).
		Where("event_id = ? AND status = 'active'", eventID).
		Update("status", "cancelled").Error
}

func (d *TicketDAOImpl) Save(ticket *domain.Ticket) error {
	return d.db.Save(ticket).Error
}

func (d *TicketDAOImpl) WithTransaction(fn func(TxContext) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

var _ TicketDAO = (*TicketDAOImpl)(nil)
