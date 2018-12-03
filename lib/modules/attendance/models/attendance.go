package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

const (
	StatusNormal = 1
	StatusDeleted = 99
)

const (
	CheckInPeriodSecondly = 2	// deprecated
	CheckInPeriodMinutely = 3	// deprecated
	CheckInPeriodHourly = 4
	CheckInPeriodDaily = 5
	CheckInPeriodWeekly = 6 	// unimplement
	CheckInPeriodMonthly = 7
	CheckInPeriodYearly = 8 // deprecated
)

// its a checkin project
type AttendanceActivity struct {
	Aid int	`orm:"pk;auto"`
	Name string `orm:"size(32);unique"`

	ValidTimeStart string `orm:"type(datetime)"`
	ValidTimeEnd string `orm:"type(datetime)"`

	// unique(uid+checkInKey+checkInPeriod)
	// rule for daily work:[{"WORKUP":{"timespan":["00:00","10:00"]}},{"WORKOFF":{"timespan":["18:00","23:59"]}}]
	// rule for daily health:[{"HEALTH":{"timespan":["06:00","09:00"]}}]
	// rule for monthly report:[{"REPORTM":{"dayspan":["01","02"]}}]
	// rule for daily checkin:{"CHECKIND":{"timespan":["00:00","23:59"]}}
	CheckInRule     string `orm:"size(4095)"`           // json, rule for checkin
	CheckInPeriod   int8   `orm:"type(tinyint)"`        // Daily Hourly Monthly
	BonusNeedStep   int    `orm:"type(int);default(0)"` // step count by CheckInPeriod
	AwardPerCheckIn int    `orm:"type(int);default(0)"` // get bonus per check in
	//AwardMap string // {"1":10,"2":25,"5":100} 1day award 10, 2day award 25 // needed?
	CreatorUid int `orm:"type(int);default(0)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Status int8 `orm:"type(tinyint);default(1)"`
	JoinPrice int	`orm:"type(int);default(0)"`
	JoinedUserCount int `orm:"type(int);default(0)"`
	// loser lost all, or percent of his all
	LoserWastagePercent float32 `orm:"digits(12);decimals(4)"`

	// use Wasting Rule?
	//BonusRule string // use Bonus Rule?
	//Top N player can get Bonus?
	// leverage?

}



// unique (user+aid+(status=1))
// user - 1:N - jal
type JoinActivityLog struct {
	JalId         int       `orm:"pk;auto"`
	Aid           int       `orm:"type(int)"`
	Uid           int       `orm:"type(int)"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
	StartDate     string    `orm:"type(datetime)"` //mysql.Date?
	BonusNeedStep int       `orm:"type(int);default(0)"`
	Step          int       `orm:"type(int);default(0)"`
	LastStepDate  string    `orm:"time(datetime)"` // needed?
	IsFinish      byte      `orm:"type(tinyint);default(0)"` // is finishing, w
	RewardDispatched byte `orm:"type(tinyint);default(0)"`
	JoinUtlId int `orm:"type(int);default(0)"`
	JoinPrice int `orm:"type(int);default(0)"`
	Status byte `orm:"type(tinyint);default(1)"` // missed,expired,stopped,deleted,shared cannot restart
	//IsMissed int // is wasted
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
	CilId int `orm:"pk;auto"`
	//JalId int `orm:"-"` // needed?
	Uid int `orm:"type(int)"`
	Aid int `orm:"type(int);default(0)"`
	CheckInKeyType string `orm:"size(32)"`
	CheckInKey string `orm:"size(32)"`	// unique of a user
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Status byte `orm:"type(tinyint);default(1)"`
}

const (
	TradeTypeCharge = 1
	TradeTypeCheckInAward = 2
	TradeTypeRevertConsume = 3
	TradeTypeCheckInBonus = 4
	TradeTypeConsume = 21
	TradeTypeWastage = 22
	TradeTypeRevertCharge = 23
)

// user - 1:N - utl
type UserTradeLog struct {
	UtlId int `orm:"pk;auto"`
	Uid int `orm:"type(int)"`
	Amount int `orm:"type(int)"`
	TradeType byte `orm:"type(tinyint)"`
	PlusMinus byte `orm:"type(tinyint)"` // -1 or +1
	SourceType byte `orm:"type(tinyint);default(0)"`
	Balance int `orm:"type(int);default(0)"`
	PayStatus byte `orm:"type(tinyint);default(0)"`
	RefundStatus byte `orm:"type(tinyint);default(0)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Status byte `orm:"type(tinyint);default(1)"`
	RelatedUtlId int `orm:"type(int);default(0)"`
	WastageDetail string `orm:"type(text)"`
	ConsumeDetail string `orm:"type(text)"`
	Remark string `orm:"type(text)"`
	PayId string `orm:"type(int);default(0)"`
	PayInfo string `orm:"type(text)"`
}


// user - N:N - wsid
// UtlId - 1:N - wsid
type WastageShare struct {
	WsId int `orm:"pk;auto"`
	//UtlId int `orm:""`
	JalId int `orm:"type(int)"`
	FromUid int `orm:"type(int)"`
	ToUid int `orm:"type(int)"`
	Amount int `orm:"type(int)"`

	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Status byte `orm:"type(tinyint);default(1)"`
}

// Uid - 1:1 - UbId
type UserBalance struct {
	UbId int `orm:"pk;auto"`
	Uid int `orm:"type(int)"`
	Balance int64 `orm:"type(int);default(0)"` // cent
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Status byte `orm:"type(tinyint);default(1)"`
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
	Status byte `orm:"-"`
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
	Status byte `orm:"-"`
}

func (t *AttendanceActivity) TableEngine() string {
	return "INNODB"
}
func (t *JoinActivityLog) TableEngine() string {
	return "INNODB"
}
func (t *UserTradeLog) TableEngine() string {
	return "INNODB"
}
func (t *UserBalance) TableEngine() string {
	return "INNODB"
}
func (t *WastageShare) TableEngine() string {
	return "INNODB"
}
func (t *CheckInLog) TableEngine() string {
	return "INNODB"
}


func init () {
	orm.RegisterModelWithPrefix(
		beego.AppConfig.String("dbprefix"),
		new(AttendanceActivity),
		new(JoinActivityLog),
		new(WastageShare),
		new(UserTradeLog),
		new(CheckInLog),
		new(UserBalance))

}