//go:build linux || freebsd

package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateEventFilterVolume verifies the VOLUME filter logic fixed in
// #29398, where `podman events --filter volume=<id>` silently returned no
// events because only a name-prefix match was performed with no ID fallback.
func TestGenerateEventFilterVolume(t *testing.T) {
	const (
		volName = "my-data-volume"
		volID   = "a1b2c3d4e5f6"
	)

	tests := []struct {
		name        string
		filterValue string
		event       Event
		want        bool
	}{
		{
			name:        "exact volume name match",
			filterValue: volName,
			event:       Event{Type: Volume, Name: volName, ID: volID},
			want:        true,
		},
		{
			name:        "partial name prefix match",
			filterValue: "my-data",
			event:       Event{Type: Volume, Name: volName, ID: volID},
			want:        true,
		},
		{
			name:        "ID prefix match",
			filterValue: volID[:6],
			event:       Event{Type: Volume, Name: volName, ID: volID},
			want:        true,
		},
		{
			name:        "full ID match",
			filterValue: volID,
			event:       Event{Type: Volume, Name: volName, ID: volID},
			want:        true,
		},
		{
			name:        "wrong volume name does not match",
			filterValue: "other-volume",
			event:       Event{Type: Volume, Name: volName, ID: volID},
			want:        false,
		},
		{
			name:        "wrong event type does not match",
			filterValue: volName,
			event:       Event{Type: Container, Name: volName, ID: volID},
			want:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filterFn, err := generateEventFilter("volume", tc.filterValue)
			require.NoError(t, err)
			assert.Equal(t, tc.want, filterFn(&tc.event))
		})
	}
}

// TestGenerateEventFilterVolumeUnknownReturnsError verifies that an
// unrecognised filter key still returns an error (regression guard).
func TestGenerateEventFilterVolumeUnknownReturnsError(t *testing.T) {
	_, err := generateEventFilter("nosuchfilter", "value")
	require.Error(t, err)
}
