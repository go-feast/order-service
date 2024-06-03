package order

import "gorm.io/gorm"

type GormOrderModel struct { //nolint:govet
	gorm.Model
	Order
}
