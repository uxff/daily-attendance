package attendance

import (
	"time"
)

//
type AttendanceActivity struct {
	Aid int
	Name string `orm:"size(32);unique"`

	ValidTimeStart string `orm:""`
	ValidTimeEnd string `orm:""`

	// rule for daily work:[{"WORK_UP":{"timespan":["00:00","10:00"]}},{"WORK_OFF":{"timespan":["18:00","23:59"]}}]
	// rule for daily health:[{"HEALTH":{"timespan":["06:00","09:00"]}}]
	// rule for monthly report:[{"MONTH_UPPER":{"dayspan":["01","02"]}}]
	Rule string `orm:""` // json
	//DailyNeedTimes int `orm:""`
	NeedStep int `orm:""`
	CheckinType int `orm:""` // Daily Hourly Monthly
	Creator int `orm:""`
	Created time.Time `orm:""`
	Updated time.Time `orm:""`
	JoinUidCount int `orm:""`
	//ShareNeedTimes int //
	LoserWastagePercent float32
}

// user+aid - 1:N - jal
type JoinActivityLog struct {
	JalId int
	Aid int

	Uid int
	Created time.Time
	Updated time.Time `orm:""`
	StartDate string //mysql.Date
	NeedStep int
	Step int
	IsFinish int
	RewardDispatched int
	Status int // is deleted, can restart
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
	CilId int
	JalId int
	Aid int
	Uid int
	CheckInKeyType string
	CheckInKey string	// unique
	Created time.Time
	Updated time.Time
	Status int
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
	UtlId int
	Uid int
	Amount int
	TradeType int
	SourceType int
	Balance int
	PayStatus int
	Craeted time.Time
	Updated time.Time
	Status int
	OriginUtlId int
	WastageDetail string
	ConsumeDetail string
	Remark string
	PayInfo string
}


// user - N:N - wsid
// UtlId - 1:N - wsid
type WastageShare struct {
	WsId int
	UtlId int
	FromUid int
	ToUid int
	Amount int

	Craeted time.Time
	Updated time.Time
	Status int
}

// Uid - 1:1 - UbId
type UserBalance struct {
	UbId int
	Uid int
	Balance int64 // cent
	Craeted time.Time
	Updated time.Time
	Status int
}

type RankCheckIn struct {
	RciId int
	Rank int
	Uid int
	Aid int
	// 连续check次数
	CheckInTimes int
	Created time.Time
	Updated time.Time
	Status int
}

type RankAwardAmount struct {
	RaaId int
	Rank int
	Uid int
	Aid int
	AwardAmount int
	AwardTimes int
	// 连续check次数
	CheckInTimes int
	Created time.Time
	Updated time.Time
	Status int
}
