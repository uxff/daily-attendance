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

func Charge() {

}

func Award(Uid int, amount int, tradeType int8, remark string) (int, error) {

	logs.Info("WILL AWARD: uid:%d amount:%d tradeType:%d", Uid, amount, tradeType)
	ub := GetUserBalance(Uid)
	utl := models.UserTradeLog{
		Uid:          Uid,
		Amount:       amount,
		TradeType:    tradeType,
		SourceType:   1,
		Balance:      ub.Balance + int64(amount),
		PayStatus:    models.PayStatusSuccess,
		PlusMinus:    models.UserTradePlus,
		Status:       models.StatusNormal,
		RefundStatus: 0,
		Remark:       remark,
	}

	ormObj := orm.NewOrm()

	utlId, err := ormObj.Insert(&utl)
	if err != nil {
		logs.Error("when award error:%v uid:%d tradeType:%d amount:%d", err, Uid, amount, tradeType)
		return int(utlId), err
	}

	ub.Balance = ub.Balance + int64(amount)
	_, err = ormObj.Update(ub, "balance")
	if err != nil {
		logs.Error("when award error:%v uid:%d tradeType:%d amount:%d", err, Uid, amount, tradeType)
		//return int(utlId), err
	}

	return int(utlId), nil
}

// means buy
func Consume(Uid int, p BuyableObject, num int, remark string) (int, error) {
	price := p.GetPrice()
	ub := GetUserBalance(Uid)
	utl := models.UserTradeLog{
		Uid:          Uid,
		Amount:       p.GetPrice(),
		TradeType:    models.TradeTypeConsume,
		SourceType:   1,
		Balance:      ub.Balance - int64(price),
		PayStatus:    models.PayStatusSuccess,
		PlusMinus:    models.UserTradeMinus,
		Status:       models.StatusNormal,
		RefundStatus: 0,
		Remark:       remark,
	}

	if p.GetStoredNum() < num {
		logs.Error("no enough store of pid:%d( %s), need %d, remain %d", p.GetProductId(), p.GetName(), num, p.GetStoredNum())
		return 0, fmt.Errorf("no enough store of pid:%d( %s), need %d, remain %d", p.GetProductId(), p.GetName(), num, p.GetStoredNum())
	}

	ormObj := orm.NewOrm()

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
