package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/* common status defined for each entity */
const (
	StatusNormal  = 1
	StatusDeleted = 99
)

const (
	CheckInPeriodSecondly = 2 // deprecated
	CheckInPeriodMinutely = 3 // deprecated
	CheckInPeriodHourly   = 4
	CheckInPeriodDaily    = 5
	CheckInPeriodWeekly   = 6 // unimplement
	CheckInPeriodMonthly  = 7
	CheckInPeriodYearly   = 8 // deprecated
)

var CheckInPeriodMap = map[int8]string{
	CheckInPeriodSecondly: "秒",
	CheckInPeriodMinutely: "分钟",
	CheckInPeriodHourly:   "小时",
	CheckInPeriodDaily:    "天",
	CheckInPeriodWeekly:   "周",
	CheckInPeriodMonthly:  "月",
}

// its a checkin project
type AttendanceActivity struct {
	Aid  int    `orm:"pk;auto"`
	Name string `orm:"size(32);unique"`

	ValidTimeStart string `orm:"size(10)"`
	ValidTimeEnd   string `orm:"size(10)"`

	// unique(uid+checkInKey+checkInPeriod)
	// rule for daily work:[{"WORKUP":{"timespan":["00:00","10:00"]}},{"WORKOFF":{"timespan":["18:00","23:59"]}}]
	// rule for daily health:[{"HEALTH":{"timespan":["06:00","09:00"]}}]
	// rule for monthly report:[{"REPORTM":{"dayspan":["01","02"]}}]
	// rule for daily checkin:{"CHECKIND":{"timespan":["00:00","23:59"]}}
	CheckInRule     string `orm:"size(255)"`            // json, rule for checkin
	CheckInPeriod   int8   `orm:"type(tinyint)"`        // Daily Hourly Monthly
	BonusNeedStep   int    `orm:"type(int);default(0)"` // step count by CheckInPeriod
	AwardPerCheckIn int    `orm:"type(int);default(0)"` // get bonus per check in
	//AwardMap string // {"1":10,"2":25,"5":100} 1day award 10, 2day award 25 // needed?
	CreatorUid int       `orm:"type(int);default(0)"`
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
	Status     int8      `orm:"type(tinyint);default(1)"`
	JoinPrice  int       `orm:"type(int);default(0)"`
	// loser lost all, or percent of his all
	//LoserWastagePercent float32 `orm:"digits(12);decimals(4)"`

	JoinedUserCount int `orm:"type(int);default(0)"` //累计
	JoinedAmount    int `orm:"type(int);default(0)"` //累计 status(all) JoinedAmount>MissedAmount
	MissedUserCount int `orm:"type(int);default(0)"` //累计
	AllMissedAmount int `orm:"type(int);default(0)"` //累计 status(missed,stopped,shared)
	SharedAmount    int `orm:"type(int);default(0)"` //累计 status(shared)
	UnsharedAmount  int `orm:"type(int);default(0)"` //实时值 status(missed,stopped) MissedAmount=UnsharedAmount+SharedAmount

	Desc string `orm:"size(255)"`

	// use Wasting Rule?
	//BonusRule string // use Bonus Rule?
	//Top N player can get Bonus?
	// leverage?

}

const (
	JalStatusInited   = 1
	JalStatusAchieved = 2 // more than 5 days
	JalStatusMissed   = 3 // disachived, missed, same as wasted
	JalStatusStopped  = 4 // stopped by user manual
	JalStatusShared   = 5 //
)

var JalStatusMap = map[int8]string{
	JalStatusInited:   "坚持中",
	JalStatusAchieved: "达标分红中",
	JalStatusMissed:   "坚持未遂",
	JalStatusStopped:  "停止",
	JalStatusShared:   "被瓜分",
}

// unique (user+aid+(status=1))
// user - 1:N - jal
type JoinActivityLog struct {
	JalId int                 `orm:"pk;auto"`
	Aid   *AttendanceActivity `orm:"rel(fk);default(0);null"`
	//Aidd           int       `orm:"type(int)"`
	Uid           int       `orm:"type(int)"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
	BonusNeedStep int       `orm:"type(int);default(0)"`
	Step          int       `orm:"type(int);default(0)"`
	IsFinish      int8      `orm:"type(tinyint);default(0)"` // is finishing, w
	JoinUtlId     int       `orm:"type(int);default(0)"`
	JoinPrice     int       `orm:"type(int);default(0)"`
	Status        int8      `orm:"type(tinyint);default(1)"` // missed,expired,stopped,deleted,shared cannot restart
	BonusTotal    int       `orm:"type(int);default(0)"`
	StartDate     string    `orm:"size(10)"`   //mysql.Date?
	LastStepDate  string    `orm:"size(10)"`   // needed?
	Schedule      string    `orm:"type(text)"` // json of
	//Schedulemap      map[string]string `orm:"type(json);defualt('')"`
	//IsMissed int // is wasted

}

const (
	CheckInKeyTypeWorkUp        = "WORKUPD"
	CheckInKeyTypeWorkOff       = "WORKOFFD"
	CheckInKeyTypeDailyHealth   = "HEALTHD"
	CheckInKeyTypeMonthlyReport = "REPORTM"
)

// user - 1:N - CilId
// CheckInKey 1:1 CilId
type CheckInLog struct {
	CilId          int       `orm:"pk;auto"`
	JalId          int       `orm:"type(int)"` // needed? need tobe struct JoinActivityLog?
	Uid            int       `orm:"type(int)"`
	Aid            int       `orm:"type(int);default(0)"`
	CheckInKeyType string    `orm:"size(32)"`
	CheckInKey     string    `orm:"size(32)"` // unique of a user
	Created        time.Time `orm:"auto_now_add;type(datetime)"`
	Updated        time.Time `orm:"auto_now;type(datetime)"`
	Status         int8      `orm:"type(tinyint);default(1)"`
	//Map            map[string]string `orm:"type(json);defualt('')"`
}

const (
	TradeTypeCharge        = 1
	TradeTypeCheckInAward  = 2
	TradeTypeRevertConsume = 3
	TradeTypeCheckInBonus  = 4
	TradeTypeRegisterAward = 5
	TradeTypeConsume       = 21
	TradeTypeWastage       = 22
	TradeTypeRevertCharge  = 23
)

var TradeTypeMap = map[int8]string{
	TradeTypeCharge:        "充值",
	TradeTypeCheckInAward:  "打卡奖励",
	TradeTypeRevertConsume: "撤销消费",
	TradeTypeCheckInBonus:  "打卡分红",
	TradeTypeRegisterAward: "注册奖励",
	TradeTypeConsume:       "消费",
	TradeTypeWastage:       "被瓜分押金",
	TradeTypeRevertCharge:  "撤销充值",
}

const (
	PayStatusNone    = 1
	PayStatusSuccess = 2
	PayStatusFail    = 3
)

const (
	UserTradePlus  = 1
	UserTradeMinus = -1
)

// user - 1:N - utl
type UserTradeLog struct {
	UtlId         int       `orm:"pk;auto"`
	Uid           int       `orm:"type(int)"`
	Amount        int       `orm:"type(int)"`
	TradeType     int8      `orm:"type(tinyint)"`
	PlusMinus     int8      `orm:"type(tinyint)"` // -1 or +1
	SourceType    int8      `orm:"type(tinyint);default(0)"`
	Balance       int64     `orm:"type(int);default(0)"`
	PayStatus     int8      `orm:"type(tinyint);default(0)"`
	RefundStatus  int8      `orm:"type(tinyint);default(0)"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
	Status        int8      `orm:"type(tinyint);default(1)"`
	RelatedUtlId  int       `orm:"type(int);default(0)"`
	WastageDetail string    `orm:"type(text)"`
	ConsumeDetail string    `orm:"type(text)"`
	Remark        string    `orm:"type(text)"`
	PayId         string    `orm:"type(int);default(0)"`
	PayInfo       string    `orm:"type(text)"`
}

// user - N:N - wsid
// UtlId - 1:N - wsid
type WastageShare struct {
	WsId        int `orm:"pk;auto"`
	WastedJalId int `orm:"type(int)"`
	ToJalId     int `orm:"type(int)"`
	FromUid     int `orm:"type(int)"`
	ToUid       int `orm:"type(int)"`
	UtlId       int `orm:"type(int)"`
	Amount      int `orm:"type(int)"`

	Aid *AttendanceActivity `orm:"rel(fk);default(0);null"`

	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Status  int8      `orm:"type(tinyint);default(1)"`
}

// Uid - 1:1 - UbId
type UserBalance struct {
	UbId         int       `orm:"pk;auto"`
	Uid          int       `orm:"type(int)"`
	Balance      int64     `orm:"type(bigint);default(0)"` // cent
	TotalIncome  int64     `orm:"type(bigint);default(0)"`
	TotalExpense int64     `orm:"type(bigint);default(0)"`
	Created      time.Time `orm:"auto_now_add;type(datetime)"`
	Updated      time.Time `orm:"auto_now;type(datetime)"`
	Status       int8      `orm:"type(tinyint);default(1)"`
}

type RankCheckIn struct {
	RciId int `orm:"-"`
	Rank  int `orm:"-"`
	Uid   int `orm:"-"`
	Aid   int `orm:"-"`
	// 连续check次数
	CheckInTimes int       `orm:"-"`
	Created      time.Time `orm:"-"`
	Updated      time.Time `orm:"-"`
	Status       int8      `orm:"-"`
}

type RankAwardAmount struct {
	RaaId       int `orm:"-"`
	Rank        int `orm:"-"`
	Uid         int `orm:"-"`
	Aid         int `orm:"-"`
	AwardAmount int `orm:"-"`
	AwardTimes  int `orm:"-"`
	// 连续check次数
	CheckInTimes int       `orm:"-"`
	Created      time.Time `orm:"-"`
	Updated      time.Time `orm:"-"`
	Status       int8      `orm:"-"`
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

func init() {
	orm.RegisterModelWithPrefix(
		beego.AppConfig.String("dbprefix"),
		new(AttendanceActivity),
		new(JoinActivityLog),
		new(WastageShare),
		new(UserTradeLog),
		new(CheckInLog),
		new(UserBalance))

}
