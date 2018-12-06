package attendance

import "github.com/uxff/daily-attendance/lib/modules/attendance/models"

type BuyableAttendanceActivity struct {
	*models.AttendanceActivity
	BuyableObject
}

func (b *BuyableAttendanceActivity) GetProductId() int {
	return b.AttendanceActivity.Aid
}

func (b *BuyableAttendanceActivity) GetPrice() int {
	return b.AttendanceActivity.JoinPrice
}
func (b *BuyableAttendanceActivity) GetName() string {
	return b.AttendanceActivity.Name
}

func (b *BuyableAttendanceActivity) GetStoreNum() int {
	return 1
}
func (b *BuyableAttendanceActivity) OnOrderCreate(Uid, UtlId, num int) {
}

func (b *BuyableAttendanceActivity) OnBuySuccess(Uid, UtlId, num int) {
	//userJoinActivityStart
}
func (b *BuyableAttendanceActivity) OnRefund(Uid, RefundUtlId, num int) {
}

func ActivityToProduct(act *models.AttendanceActivity) BuyableObject {
	return &BuyableAttendanceActivity{AttendanceActivity:act}
}


