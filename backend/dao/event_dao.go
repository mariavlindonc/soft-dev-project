package db

import (
	"backend/domain"

	"gorm.io/gorm"
)

type EventDAOImpl struct {
	db *gorm.DB
}

func NewEventDAO(db *gorm.DB) *EventDAOImpl {
	return &EventDAOImpl{db: db}
}

func (d *EventDAOImpl) FindAll(filters domain.EventFilters) ([]domain.Event, error) {
	query := d.db.Model(&domain.Event{}).Preload("CreatedBy")

	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.DateFrom != nil {
		query = query.Where("event_date >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("event_date <= ?", *filters.DateTo)
	}

	var events []domain.Event
	err := query.Order("event_date ASC").Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (d *EventDAOImpl) FindByID(id uint) (*domain.Event, error) {
	var event domain.Event
	err := d.db.Preload("CreatedBy").Take(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (d *EventDAOImpl) Create(event *domain.Event) error {
	return d.db.Create(event).Error
}

func (d *EventDAOImpl) Update(event *domain.Event) error {
	return d.db.Save(event).Error
}

func (d *EventDAOImpl) Delete(id uint) error {
	return d.db.Delete(&domain.Event{}, id).Error
}

var _ EventDAO = (*EventDAOImpl)(nil)
