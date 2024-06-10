package services

import (
	"testing"

	"github.com/sshaparenko/donation-service/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestGetAllItems_Success(t *testing.T) {
	var result []int = services.GetAllItems()
	assert.ElementsMatch(t, result, []int{1, 2, 3})
}
