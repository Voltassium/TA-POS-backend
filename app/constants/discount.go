package constants

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
)

func (d DiscountType) IsValidEnum() bool {
	switch d {
	case DiscountTypePercentage, DiscountTypeFixed:
		return true
	default:
		return false
	}
}
