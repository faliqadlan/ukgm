package utils_test

import (
	"fmt"
	"testing"

	"github.com/faliqadlan/ukgm/internal/app/utils"
)

func TestRateDivider(t *testing.T) {
	itemsAmmount := 6
	index := 6

	rate := utils.RateDivider(index, itemsAmmount)
	fmt.Println(rate)
	return
}
