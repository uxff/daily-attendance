package models

import (
	"time"
)

const (
	CheckInTypeHourly = 4
	CheckInTypeDaily = 5
	CheckInTypeMonthly = 7
)

const (
	// CheckInKeyType extendible
	CheckInKeyWorkUp = "WORKUP"
	CheckInKeyWorkOff = "WORKOFF"
	CheckInKeyHealthDaily = "HEALTHD"
)
//
type AttendanceActivity struct {
	Aid int
	Name string `orm:"size(32);unique"`

	ValidTimeStart string `orm:"-"`
	ValidTimeEnd string `orm:"-"`

	// rule for daily work:[{"WORKUP":{"timespan":["00:00","10:00"]}},{"WORKOFF":{"timespan":["18:00","23:59"]}}]
	// rule for daily health:[{"HEALTH":{"timespan":["06:00","09:00"]}}]
	// rule for monthly report:[{"REPORTM":{"dayspan":["01","02"]}}]
	CheckInRule string `orm:"-"` // json, rule for checkin
	NeedStep int `orm:"-"`
	CheckInType int `orm:"-"` // Daily Hourly Monthly
	CreatorUid int `orm:"-"`
	Created time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	JoinPrice int	`orm:"-"`
	JoinedUserCount int `orm:"-"`
	// loser lost all, or percent of his all
	LoserWastagePercent float32 `orm:"-"`
	// use Wasting Rule?
	//BonusRule string // use Bonus Rule?
	//Top N player can get Bonus?
	// leverage?

}

// unique (user+aid+(status=1))
// user - 1:N - jal
type JoinActivityLog struct {
	JalId int `orm:"-"`
	Aid int `orm:"-"`
	Uid int `orm:"-"`
	Created time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	StartDate string `orm:"-"` //mysql.Date
	NeedStep int `orm:"-"`
	Step int `orm:"-"`
	LastStepDate string `orm:"-"`
	IsFinish int `orm:"-"` // is finishing, w
	//IsMissed int // is wasted
	RewardDispatched int `orm:"-"`
	JoinUtlId int `orm:"-"`
	Status int `orm:"-"` // missed,deleted, can restart
}

const (
	CheckInKeyTypeWorkUp = "WORKUPD"
	CheckInKeyTypeWorkOff = "WORKOFFD"
	CheckInKeyTypeDailyHealth = "HEALTHD"
	CheckInKeyTypeMonthlyReport = "REPORTM"
)
// user - 1:N - CilId
// CheckInKey 1:1 CilId
type CheckInLog struct {
	CilId int `orm:"-"`
	//JalId int `orm:"-"` // needed?
	Aid int `orm:"-"`
	Uid int `orm:"-"`
	CheckInKeyType string `orm:"-"`
	CheckInKey string `orm:"-"`	// unique
	Created time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	Status int `orm:"-"`
}

const (
	TradeTypeCharge = 1
	TradeTypeCheckInAward = 2
	TradeTypeRevertConsume = 3
	TradeTypeConsume = 21
	TradeTypeWastage = 22
	TradeTypeRevertCharge = 23
)

// user - 1:N - utl
type UserTradeLog struct {
	UtlId int `orm:"-"`
	Uid int `orm:"-"`
	Amount int `orm:"-"`
	TradeType int `orm:"-"`
	SourceType int `orm:"-"`
	Balance int `orm:"-"`
	PayStatus int `orm:"-"`
	RefundStatus int `orm:"-"`
	Craeted time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	Status int `orm:"-"`
	RelatedUtlId int `orm:"-"`
	WastageDetail string `orm:"-"`
	ConsumeDetail string `orm:"-"`
	Remark string `orm:"-"`
	PayId string `orm:"-"`
	PayInfo string `orm:"-"`
}


// user - N:N - wsid
// UtlId - 1:N - wsid
type WastageShare struct {
	WsId int `orm:"-"`
	UtlId int `orm:"-"`
	FromUid int `orm:"-"`
	ToUid int `orm:"-"`
	Amount int `orm:"-"`

	Craeted time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	Status int `orm:"-"`
}

// Uid - 1:1 - UbId
type UserBalance struct {
	UbId int `orm:"-"`
	Uid int `orm:"-"`
	Balance int64 `orm:"-"` // cent
	Craeted time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	Status int `orm:"-"`
}

type RankCheckIn struct {
	RciId int `orm:"-"`
	Rank int `orm:"-"`
	Uid int `orm:"-"`
	Aid int `orm:"-"`
	// 连续check次数
	CheckInTimes int `orm:"-"`
	Created time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	Status int `orm:"-"`
}

type RankAwardAmount struct {
	RaaId int `orm:"-"`
	Rank int `orm:"-"`
	Uid int `orm:"-"`
	Aid int `orm:"-"`
	AwardAmount int `orm:"-"`
	AwardTimes int `orm:"-"`
	// 连续check次数
	CheckInTimes int `orm:"-"`
	Created time.Time `orm:"-"`
	Updated time.Time `orm:"-"`
	Status int `orm:"-"`
}
