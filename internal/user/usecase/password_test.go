package usecase

import (
	"github.com/stretchr/testify/assert"
	"go-todolist-sber/internal/apperror"
	"testing"
)

func TestArgon_GenerateHashFromPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     string
		wantErr  error
	}{
		{
			name:     "ok",
			password: "Test12414",
			want:     "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
			wantErr:  nil,
		},
		{
			name:     "ok",
			password: "newAccount",
			want:     "47e153353c4007136518e13fe79610e19e0705162c31ee8c99831c6e8ce28172",
			wantErr:  nil,
		},
		{
			name:     "password not valid",
			password: "123",
			want:     "",
			wantErr:  apperror.ErrDataNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			a := NewArgonPassword("")
			hash, err := a.GenerateHashFromPassword(tt.password)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, hash)
		})
	}
}

func TestArgon_VerifyPassword(t *testing.T) {
	tests := []struct {
		name         string
		hashPassword string
		password     string
		wantErr      error
	}{
		{
			name:         "not equal",
			hashPassword: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
			password:     "test",
			wantErr:      apperror.ErrHashPasswordsNotEqual,
		},
		{
			name:         "ok",
			hashPassword: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
			password:     "Test12414",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			a := NewArgonPassword("")
			err := a.VerifyPassword(tt.hashPassword, tt.password)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
