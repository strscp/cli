package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHasAllFlag_NotSet(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("all", false, "fetch all")
	assert.False(t, hasAllFlag(cmd))
}

func TestHasAllFlag_Set(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("all", false, "fetch all")
	_ = cmd.Flags().Set("all", "true")
	assert.True(t, hasAllFlag(cmd))
}

func TestHasAllFlag_NoFlag(t *testing.T) {
	cmd := &cobra.Command{}
	assert.False(t, hasAllFlag(cmd))
}

func TestWatchIntervalParsing(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"10s", false},
		{"1m", false},
		{"5m", false},
		{"30s", false},
		{"3s", false},  // valid duration but below minimum
		{"bad", true},  // invalid duration
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := time.ParseDuration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
