package safeweb_lib_helper

import (
    "fmt"
    "time"
    
    "github.com/spf13/viper"

    "safeweb.app/service/safeweb_lib/types/common"
)

func IsNowActive(startDate time.Time, endDate time.Time, timeRanges []safeweb_lib_types_common.TimeRange) bool {
    if IsZero(startDate, endDate) {
        return true
    }
    now := GetNow()
    hour, min, sec := now.Clock()
    nowTime := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
    active := Between(now, startDate, endDate)
    if active {
        for _, timeRange := range timeRanges {
            if nowTime < timeRange.Start || nowTime > timeRange.End {
                return false
            }
        }
        return true
    }
    return false
}

func GetNow() time.Time {
    loc, _ := time.LoadLocation(viper.GetString("TIMEZONE"))
    now := time.Now().In(loc)
    return now
}

func Between(a time.Time, start time.Time, end time.Time) bool {
    return (a.Equal(start) || a.After(start)) && (a.Equal(end) || a.Before(end))
}

func GetTomorrow() time.Time {
    return time.Now().AddDate(0, 0, 1)
}

func GetYesterday() time.Time {
    return time.Now().AddDate(0, 0, -1)
}
