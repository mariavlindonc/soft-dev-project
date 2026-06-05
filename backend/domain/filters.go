package domain

import "time"

// EventFilters is a shared query-parameter type used by both the
// service and DAO layers to filter event listings. It lives in
// domain to avoid an import cycle between services and dao.
type EventFilters struct {
	Category string
	DateFrom *time.Time
	DateTo   *time.Time
}
