package main

import (
	"testing"
	"time"

	// "github.com/go-playground/assert/v2"
	"github.com/sxc/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 5, 29, 19, 34, 0, 0, time.UTC),
			want: "29 May 2023 at 19:34",
		},

		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 5, 29, 19, 34, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "29 May 2023 at 18:34",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			// if hd != tt.want {
			// 	t.Errorf("got %q; want %q", hd, tt.want)
			// }
			assert.Equal(t, hd, tt.want)
		})
	}
}
