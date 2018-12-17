package attendance

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/uxff/daily-attendance/lib/modules/attendance/models"
)

type PayableObject interface {
	PaySuccessCallback()
	PayFailCallback(error)
}

type PaidBill struct {
	PaymentId int64
	Amount    int
	PaidTime  time.Time
}

type BuyableObject interface {
	GetProductId() int
	GetPrice() int
	GetStoredNum() int
	GetName() string
	OnOrderCreate(Uid, UtlId, num int)
	OnBuySuccess(Uid, UtlId, num int)
	OnRefund(Uid, RefundUtlId, num int)
}

type Product struct {
	BuyableObject
}

// no need
func PrePay() (err error, paymentUserAuthUrl string) {

	return nil, ""
}

// no need
func Pay() {

}

func Charge(Uid int, amount int, remark string) (int, error) {

	utl := models.UserTradeLog{
		Uid:          Uid,
		Amount:       amount,
		TradeType:    models.TradeTypeCheckInBonus,
		SourceType:   1,
		Balance:      0,
		PayStatus:    models.PayStatusSuccess,
		PlusMinus:    models.UserTradePlus,
		Status:       models.StatusNormal,
		RefundStatus: 0,
		Remark:       remark,
	}

	ormObj := orm.NewOrm()

	utlId, err := ormObj.Insert(&utl)
	if err != nil {
		return int(utlId), err
	}

	ub := GetUserBalance(Uid)
	ub.Balance = ub.Balance + int64(amount)
	ormObj.Update(ub, "balance")

	return int(utlId), nil
}

func Consume(Uid int, p BuyableObject, num int) (int, error) {
	price := p.GetPrice()
	utl := models.UserTradeLog{
		Uid:          Uid,
		Amount:       p.GetPrice(),
		TradeType:    models.TradeTypeConsume,
		SourceType:   1,
		Balance:      0,
		PayStatus:    models.PayStatusSuccess,
		PlusMinus:    models.UserTradeMinus,
		Status:       models.StatusNormal,
		RefundStatus: 0,
		Remark:       p.GetName(),
	}

	if p.GetStoredNum() < num {
		logs.Error("no enough store of pid:%d( %s), need %d, remain %d", p.GetProductId(), p.GetName(), num, p.GetStoredNum())
		return 0, fmt.Errorf("no enough store of pid:%d( %s), need %d, remain %d", p.GetProductId(), p.GetName(), num, p.GetStoredNum())
	}

	ormObj := orm.NewOrm()

	//ub := models.UserBalance{}
	ub := GetUserBalance(Uid)
	if ub.Balance < int64(price) {
		logs.Error("no enough balance of uid:%d, need %d, remain %d", Uid, price, ub.Balance)
		return 0, fmt.Errorf("no enough balance of uid:%d, need %d, remain %d", Uid, price, ub.Balance)
	}

	ub.Balance = ub.Balance - int64(price)
	ormObj.Update(ub, "balance")

	utlId, err := ormObj.Insert(&utl)
	if err != nil {
		return int(utlId), err
	}

	p.OnBuySuccess(Uid, int(utlId), num)

	return int(utlId), nil
}

func CancelCharge() {

}

func CancelConsume() {

}
