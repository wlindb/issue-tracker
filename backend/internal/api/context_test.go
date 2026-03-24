package api_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wlindb/issue-tracker/internal/api"
)

func Test_UserIDFromContext_TableDriven(t *testing.T) {
	expectedID := uuid.New()

	tests := []struct {
		name    string
		ctx     context.Context
		wantID  uuid.UUID
		wantErr bool
	}{
		{
			name:    "valid user ID",
			ctx:     api.WithUserID(context.Background(), expectedID),
			wantID:  expectedID,
			wantErr: false,
		},
		{
			name:    "missing user ID",
			ctx:     context.Background(),
			wantID:  uuid.Nil,
			wantErr: true,
		},
		{
			name:    "uuid.Nil user ID",
			ctx:     api.WithUserID(context.Background(), uuid.Nil),
			wantID:  uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := api.UserIDFromContext(tt.ctx)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, uuid.Nil, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantID, actual)
			}
		})
	}
}
