package constants

// DiscountType defines the supported discount calculation methods.
type DiscountType string

const (
	// DiscountTypePercentage reduces the amount by a percentage of the base price.
	DiscountTypePercentage DiscountType = "percentage"
	// DiscountTypeFixed reduces the amount by a fixed monetary value.
	DiscountTypeFixed DiscountType = "fixed"
)

func (d DiscountType) IsValidEnum() bool {
	switch d {
	case DiscountTypePercentage, DiscountTypeFixed:
		return true
	default:
		return false
	}
}
