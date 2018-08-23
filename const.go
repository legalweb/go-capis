package capis

// ProductType a type of product supported by comparisonapis.com
type ProductType string

const (
	// TypeBankAccount product type.
	TypeBankAccount ProductType = "bankaccount"
	// TypeCreditCard product type.
	TypeCreditCard ProductType = "creditcard"
	// TypeMortgage product type.
	TypeMortgage ProductType = "mortgage"
)

func (t ProductType) String() string {
	return string(t)
}
