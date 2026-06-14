package services

import (
	"errors"
	"testing"
)

func TestValidateProductPricing(t *testing.T) {
	validDiscount := int64(750)
	zeroDiscount := int64(0)
	negativeDiscount := int64(-1)
	equalDiscount := int64(1000)
	higherDiscount := int64(1250)

	tests := []struct {
		name               string
		discountPriceCents *int64
		wantErr            bool
	}{
		{name: "no discount", discountPriceCents: nil},
		{name: "valid discount", discountPriceCents: &validDiscount},
		{name: "zero discount", discountPriceCents: &zeroDiscount, wantErr: true},
		{name: "negative discount", discountPriceCents: &negativeDiscount, wantErr: true},
		{name: "discount equals price", discountPriceCents: &equalDiscount, wantErr: true},
		{name: "discount exceeds price", discountPriceCents: &higherDiscount, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateProductPricing(1000, test.discountPriceCents)
			if test.wantErr && !errors.Is(err, ErrInvalidDiscountPrice) {
				t.Fatalf("expected ErrInvalidDiscountPrice, got %v", err)
			}

			if !test.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
