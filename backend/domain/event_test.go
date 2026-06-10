package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ptr(t time.Time) *time.Time { return &t }

func TestCurrentSalePhase(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		event Event
		want  SalePhase
	}{
		{
			name:  "no presale active returns PhaseNoPresale",
			event: Event{PresaleActive: false},
			want:  PhaseNoPresale,
		},
		{
			name:  "presale active but all dates nil returns PhaseNoPresale",
			event: Event{PresaleActive: true},
			want:  PhaseNoPresale,
		},
		{
			name: "presale start date nil returns PhaseNoPresale",
			event: Event{
				PresaleActive:   true,
				GeneralSaleDate: ptr(now.Add(2 * time.Hour)),
			},
			want: PhaseNoPresale,
		},
		{
			name: "general sale date nil returns PhaseNoPresale",
			event: Event{
				PresaleActive:    true,
				PresaleStartDate: ptr(now.Add(-2 * time.Hour)),
			},
			want: PhaseNoPresale,
		},
		{
			name: "before presale start returns PhaseNotYetOpen",
			event: Event{
				PresaleActive:    true,
				PresaleStartDate: ptr(now.Add(2 * time.Hour)),
				GeneralSaleDate:  ptr(now.Add(4 * time.Hour)),
			},
			want: PhaseNotYetOpen,
		},
		{
			name: "between presale and general sale returns PhasePresale",
			event: Event{
				PresaleActive:    true,
				PresaleStartDate: ptr(now.Add(-2 * time.Hour)),
				GeneralSaleDate:  ptr(now.Add(2 * time.Hour)),
			},
			want: PhasePresale,
		},
		{
			name: "after general sale date returns PhasePublic",
			event: Event{
				PresaleActive:    true,
				PresaleStartDate: ptr(now.Add(-4 * time.Hour)),
				GeneralSaleDate:  ptr(now.Add(-2 * time.Hour)),
			},
			want: PhasePublic,
		},
		{
			name: "exactly at presale start returns PhasePresale (Before is exclusive)",
			event: Event{
				PresaleActive:    true,
				PresaleStartDate: ptr(now),
				GeneralSaleDate:  ptr(now.Add(2 * time.Hour)),
			},
			want: PhasePresale,
		},
		{
			name: "exactly at general sale date returns PhasePublic",
			event: Event{
				PresaleActive:    true,
				PresaleStartDate: ptr(now.Add(-2 * time.Hour)),
				GeneralSaleDate:  ptr(now),
			},
			want: PhasePublic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.event.CurrentSalePhase(now)
			assert.Equal(t, tt.want, got)
		})
	}
}

