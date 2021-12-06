package Utils

import (
	"fmt"
	"ngo/Constant"
	"time"
)

const (
	nano  int = 19
	micro int = 16
	mill  int = 13
	sec   int = 10
)

type TimeResult struct {
	UtcTime  string
	UnixTime int64
}

func UnixTimeCovert(unixTime int64) TimeResult {
	var t time.Time
	switch len(fmt.Sprintf("%d", unixTime)) {
	case nano:
		t = time.Unix(unixTime/1e9, 0)
		break
	case micro:
		t = time.Unix(unixTime/1e6, 0)
		break
	case mill:
		t = time.Unix(unixTime/1e3, 0)
		break
	case sec:
		t = time.Unix(unixTime, 0)
		break
	default:
		t = time.Unix(0, 0)
	}
	return TimeResult{UtcTime: t.UTC().Format(Constant.TimeUtcFormat), UnixTime: t.Unix()}
}

func UnixTimeOffset(unixTime int64, offset int64) int64 {
	switch len(fmt.Sprintf("%d", unixTime)) {
	case nano:
		return unixTime + int64(time.Nanosecond*time.Duration(offset)*1e9)
	case micro:
		return unixTime + int64(time.Millisecond*time.Duration(offset)*1e6)
	case mill:
		return unixTime + int64(time.Millisecond*time.Duration(offset)*1e3)
	case sec:
		return unixTime + int64(time.Second*time.Duration(offset))
	default:
		return 0
	}
}

func UnixPointTime(unixTime int64, period int64) int64 {
	if period == 0 || period%60 != 0 {
		return unixTime
	}
	inTime, err := time.Parse(Constant.TimeUtcFormat, UnixTimeCovert(unixTime).UtcTime)
	if err != nil {
		return unixTime
	}
	min := int(period / 60)
	if min == 1 {
		sec := inTime.Second() / int(period/15) * int(period/15)
		t := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), inTime.Minute(), sec, 0, time.UTC)
		return t.UnixNano() / 1e6
	} else if min < 60 && 60%min == 0 {
		sec := inTime.Second() / int(period/15) * int(period/15)
		t := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), sec/60, sec%60, 0, time.UTC)
		return t.UnixNano() / 1e6
	} else {
		return unixTime
	}
}

func NextTurnTime(inTime time.Time, period int64) time.Time {
	//60 300 900
	if period%60 != 0 {
		return time.Time{}
	}
	min := int(period / 60)
	if min == 1 {
		date := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), inTime.Minute(), 0, 0, time.UTC)
		return date.Add(time.Second * time.Duration(period))
	} else if min < 60 && 60%min == 0 {
		m := inTime.Minute() / min * min
		date := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), m, 0, 0, time.UTC)
		return date.Add(time.Second * time.Duration(period))
	} else {
		return time.Time{}
	}
}

func NextTurnDuration(inTime time.Time, period int64) time.Duration {
	return NextTurnTime(inTime, period).Sub(inTime)
}
