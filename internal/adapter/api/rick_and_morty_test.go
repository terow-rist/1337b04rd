package api

import (
	"1337bo4rd/internal/core/port"
	"testing"
)

func TestNewRickAndMortyImplementsAvatarProvider(t *testing.T) {
	var ap port.AvatarProvider = NewRickAndMorty()
	if ap == nil {
		t.Fatal("NewRickAndMorty returned nil")
	}
}
