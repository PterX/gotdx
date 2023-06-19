package trading

import (
	"fmt"
	"gitee.com/quant1x/gox/errors"
	"golang.org/x/exp/slices"
	"time"
)

const (
	kTimeMinute                = "15:04"        // 分笔成交时间格式
	CN_SERVERTIME_FORMAT       = "15:04:05.000" // 服务器时间格式
	CN_SERVERTIME_SHORT_FORMAT = "15:04:05"     // 服务器时间格式
)

// 交易日时间相关常量
const (
	CN_MarketInitTime          = "09:00:00.000" // A股数据初始化时间
	CN_TradingStartTime        = "09:15:00.000" // A股数据开始时间
	CN_TradingSuspendBeginTime = "11:30:00.000" // A股午间休市开始时间
	CN_TradingSuspendEndTime   = "12:59:59.999" // A股午间休市结束时间
	CN_TradingStopTime         = "15:00:59.999" // A股数据结束时间
	CallAuctionAmBegin         = "09:15:00.000" // 早盘集合竞价开始时间
	CallAuctionAmEnd           = "09:27:59.999" // 早盘集合竞价结束时间
	CallAuctionPmBegin         = "14:57:00.000" // 尾盘集合竞价开始时间
	CallAuctionPmEnd           = "15:01:59.999" // 尾盘集合竞价结束时间
)

// 集合竞价时间相关常量
const (
	BEGIN_A_AUCTION   = "09:15:00" // A股上午集合竞价开始时间
	END_A_AUCTION     = "09:25:00" // A股上午集合竞价结束时间
	END_A_AUCTION_SPE = "09:26:00" // A股上午集合竞价结束时间过一分钟
	BEGIN_P_AUCTION   = "14:57:00" // A股下午集合竞价开始时间
	END_P_AUCTION     = "15:01:00" // A股下午集合竞价结束时间
	END_P_AUCTION_SPE = "15:02:00" // A股下午集合竞价结束时间过一分钟
)

// 分时数据相关常量
const (
	CN_DEFAULT_TOTALFZNUM = 240 // A股默认全天交易240分钟
	BEGIN_A_AM_HOUR       = 9   // A股开市-时
	BEGIN_A_AM_MINUTE     = 30  // A股开市-分
	END_A_AM_HOUR         = 11  // A股休市-时
	END_A_AM_MINUTE       = 30  // A股休市-分
	BEGIN_A_PM_HOUR       = 13  // A股开市-时
	BEGIN_A_PM_MINUTE     = 0   // A股开市-分
	END_A_PM_HOUR         = 15  // A股休市-时
	END_A_PM_MINUTE       = 0   // A股休市-分
)

type TimeRange struct {
	Begin time.Time
	End   time.Time
}

var (
	cnTimeRange   []TimeRange // 交易时间范围
	trAMBegin     time.Time   // 上午开盘时间
	trAMEnd       time.Time
	trPMBegin     time.Time
	trPMEnd       time.Time
	CN_TOTALFZNUM = 0 // A股全天交易的分钟数
)

var (
	ErrNoUpdateRequired = errors.New("No update required")
)

func init() {
	now := time.Now()
	trAMBegin = time.Date(now.Year(), now.Month(), now.Day(), BEGIN_A_AM_HOUR, BEGIN_A_AM_MINUTE, 0, 0, time.Local)
	trAMEnd = time.Date(now.Year(), now.Month(), now.Day(), END_A_AM_HOUR, END_A_AM_MINUTE, 0, 0, time.Local)
	tr_am := TimeRange{
		Begin: trAMBegin,
		End:   trAMEnd,
	}
	cnTimeRange = append(cnTimeRange, tr_am)

	trPMBegin = time.Date(now.Year(), now.Month(), now.Day(), BEGIN_A_PM_HOUR, BEGIN_A_PM_MINUTE, 0, 0, time.Local)
	trPMEnd = time.Date(now.Year(), now.Month(), now.Day(), END_A_PM_HOUR, END_A_PM_MINUTE, 0, 0, time.Local)
	tr_pm := TimeRange{
		Begin: trPMBegin,
		End:   trPMEnd,
	}
	_minutes := 0
	cnTimeRange = append(cnTimeRange, tr_pm)
	for _, v := range cnTimeRange {
		_minutes += int(v.End.Sub(v.Begin).Minutes())
	}
	CN_TOTALFZNUM = _minutes
}

func fixMinute(m time.Time) time.Time {
	return time.Date(m.Year(), m.Month(), m.Day(), m.Hour(), m.Minute(), 0, 0, time.Local)
}

// Minutes 分钟数
func Minutes(date ...string) int {
	// 最后1个交易日
	lastDay := LastTradeDate()
	// 默认是当天
	today := IndexToday()
	theDay := today
	if len(date) > 0 {
		theDay = FixTradeDate(date[0])
	}
	if theDay < today {
		return CN_TOTALFZNUM
	}
	if theDay != lastDay {
		return CN_TOTALFZNUM
	}
	tm := time.Now()
	//tm, _ = utils.ParseTime("2023-04-11 09:29:00")
	//tm, _ = utils.ParseTime("2023-04-11 09:30:00")
	//tm, _ = utils.ParseTime("2023-04-11 09:31:00")
	//tm, _ = utils.ParseTime("2023-04-11 11:31:00")
	//tm, _ = utils.ParseTime("2023-04-11 12:59:00")
	//tm, _ = utils.ParseTime("2023-04-11 13:00:00")
	//tm, _ = utils.ParseTime("2023-04-11 13:01:00")
	//tm, _ = utils.ParseTime("2023-04-11 14:59:00")
	//tm, _ = utils.ParseTime("2023-04-11 15:01:00")
	tm = fixMinute(tm)
	tr := slices.Clone(cnTimeRange)
	var last time.Time
	for _, v := range tr {
		if tm.Before(v.Begin) {
			last = v.Begin
			break
		}
		if tm.After(v.End) {
			last = v.End
			continue
		}
		//if !tm.After(v.Begin) {
		//	last = v.Begin
		//	break
		//}
		//if !tm.Before(v.End) {
		//	last = v.End
		//	continue
		//}
		last = tm
		break
	}

	m := int(last.Sub(trAMBegin).Minutes())
	if !last.Before(trPMBegin) {
		m -= int(trPMBegin.Sub(trAMEnd).Minutes())
	}
	return m
}

func IsTrading(date ...string) bool {
	lastDay := LastTradeDate()
	today := Today()
	if len(date) > 0 {
		today = FixTradeDate(date[0])
	}
	return lastDay == today
}

// CurrentlyTrading 今天的交易是否已经开始
func CurrentlyTrading(date ...string) bool {
	if DateIsTradingDay(date...) {
		now := time.Now()
		nowTime := now.Format(CN_SERVERTIME_FORMAT)
		return nowTime >= CN_TradingStartTime
	}
	return false
}

func IsTimeInRange(timeStr, startStr, endStr string) (bool, error) {
	// 将输入的字符串解析为Time类型
	timeVal, err := time.Parse(CN_SERVERTIME_SHORT_FORMAT, timeStr)
	if err != nil {
		return false, errors.New("invalid time format")
	}
	// 将起始和结束时间解析为Time类型
	startVal, err := time.Parse(CN_SERVERTIME_SHORT_FORMAT, startStr)
	if err != nil {
		return false, errors.New("invalid start time format")
	}
	endVal, err := time.Parse(CN_SERVERTIME_SHORT_FORMAT, endStr)
	if err != nil {
		return false, errors.New("invalid end time format")
	}
	// 检查输入时间是否在起始和结束时间间，包括起始和结束时间
	if !timeVal.Before(startVal) && !timeVal.After(endVal) {
		return true, nil
	}
	if startVal.Equal(endVal) && timeVal.Equal(startVal) {
		return true, nil
	}
	return false, nil
}

// CompareTime 比较两个时间字符串大小
// 如果t1 <= t2 返回true，否则返回false
// 如果格式不正确或转换错误，返回错误
func CompareTime(t1, t2 string) (bool, error) {
	_, err := time.ParseDuration(t1)
	if err == nil {
		err = errors.New("Invalid time duration string")
		return false, err
	}
	t1Time, err := time.Parse(CN_SERVERTIME_SHORT_FORMAT, t1)
	if err != nil {
		return false, err
	}
	t2Time, err := time.Parse(CN_SERVERTIME_SHORT_FORMAT, t2)
	if err != nil {
		return false, err
	}
	if t1Time.After(t2Time) || t1Time.Equal(t2Time) {
		return true, nil
	}
	return false, nil
}

// GetTodayTimeByString 返回当天指定时刻的时间
func GetTodayTimeByString(timeStr string) (time.Time, error) {
	layout := time.DateTime
	todayStr := fmt.Sprintf("%d-%02d-%02d %s", time.Now().Year(), time.Now().Month(), time.Now().Day(), timeStr)
	today, err := time.ParseInLocation(layout, todayStr, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return today, nil
}

type TimeStatus = int

const (
	//BeforeLastTradingDay TimeStatus = 1 << iota // 缓存非交易日, 可以更新

	ExchangeClosing   TimeStatus = -2 // 收盘收, 交易停止
	ExchangePreMarket TimeStatus = -1 // 盘前
	ExchangeSuspend   TimeStatus = 0  // 休市中, 交易暂停
	ExchangeTrading   TimeStatus = 1  // 交易中
)

// 检查时间
//
//	默认检查当前时间是否可以...
func checkTradingTimestamp(lastModified ...time.Time) (beforeLastTradeDay, isHoliday, beforeInitTime, cacheAfterInitTime, trading bool, status TimeStatus) {
	lastDay := LastTradeDate()
	timestamp := time.Now()
	if len(lastModified) > 0 {
		timestamp = lastModified[0]
	}
	status = ExchangeClosing
	// 1. 缓存时间无效
	modDate := timestamp.Format(TradingDayDateFormat)
	// 1.1 非交易日, 缓存在最后一个交易日前, 可更新
	if modDate < lastDay {
		beforeLastTradeDay = true
		return
	}
	// 2 缓存日期和最后一个交易日相同
	now := time.Now()
	today := now.Format(TradingDayDateFormat)
	// 2.1 当前日期非最后一个交易日, 也就是节假日了
	if today != lastDay {
		// 节假日
		isHoliday = true
		return
	}
	// 3. 交易日, A股市场初始化前
	currentTimestamp := now.Format(CN_SERVERTIME_FORMAT)
	if currentTimestamp < CN_MarketInitTime {
		beforeInitTime = true
		return
	}
	status = ExchangePreMarket
	// 4. 交易日, A股市场初始化后
	modTimestamp := timestamp.Format(CN_SERVERTIME_FORMAT)
	if modTimestamp >= CN_MarketInitTime {
		cacheAfterInitTime = true
	}
	// 5. 交易日, A股市场实时数据后
	if currentTimestamp >= CN_TradingStartTime && currentTimestamp <= CN_TradingStopTime {
		status = ExchangeTrading
		trading = true
		if currentTimestamp >= CN_TradingSuspendBeginTime && currentTimestamp <= CN_TradingSuspendEndTime {
			status = ExchangeSuspend
		}
	}
	return
}

// CanUpdate 数据是否可以更新
func CanUpdate(lastModified ...time.Time) (updated bool) {
	beforeLastTradeDay, isHoliday, beforeInitTime, cacheAfterInitTime, _, _ := checkTradingTimestamp(lastModified...)
	if beforeLastTradeDay {
		return true
	}
	if isHoliday {
		return false
	}
	if beforeInitTime {
		return false
	}
	return cacheAfterInitTime
}

// CanInitialize 数据是否初始化(One-time update)
func CanInitialize(lastModified ...time.Time) (toInit bool) {
	beforeLastTradeDay, isHoliday, beforeInitTime, cacheAfterInitTime, _, _ := checkTradingTimestamp(lastModified...)
	if beforeLastTradeDay {
		return true
	}
	if isHoliday {
		return false
	}
	if beforeInitTime {
		return false
	}
	return !cacheAfterInitTime
}

// CanUpdateInRealtime 能否实时更新
func CanUpdateInRealtime(lastModified ...time.Time) (updateInRealTime bool, status int) {
	_, _, _, _, updateInRealTime, status = checkTradingTimestamp(lastModified...)
	return
}

// CheckCallAuctionTime 检查当前时间是否集合竞价阶段
func CheckCallAuctionTime(timestamp time.Time) (canUpdate bool) {
	return CheckCallAuctionOpen(timestamp) || CheckCallAuctionClose(timestamp)
}

// CheckCallAuctionOpen 检查当前时间是否集合竞价阶段
func CheckCallAuctionOpen(timestamp time.Time) (canUpdate bool) {
	tm := timestamp.Format(CN_SERVERTIME_FORMAT)
	if tm >= CallAuctionAmBegin && tm < CallAuctionAmEnd {
		return true
	}
	return false
}

// CheckCallAuctionClose 检查当前时间是否集合竞价阶段
func CheckCallAuctionClose(timestamp time.Time) (canUpdate bool) {
	tm := timestamp.Format(CN_SERVERTIME_FORMAT)
	if tm >= CallAuctionPmBegin && tm < CallAuctionPmEnd {
		return true
	}
	return false
}
