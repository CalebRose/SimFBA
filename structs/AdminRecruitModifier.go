package structs

import "github.com/jinzhu/gorm"

type AdminRecruitModifier struct {
	gorm.Model
	ModifierOne       int
	ModifierTwo       float64
	WeeksOfRecruiting int
}

func (ARM *AdminRecruitModifier) SetModifierOne(val int) {
	ARM.ModifierOne = val
}

func (ARM *AdminRecruitModifier) SetModifierTwo(val float64) {
	ARM.ModifierTwo = val
}

func (ARM *AdminRecruitModifier) SetWeek(val int) {
	if val > 20 {
		return
	}
	ARM.WeeksOfRecruiting = val
}
