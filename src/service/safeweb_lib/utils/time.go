package safeweb_lib_utils

import (
    "database/sql"
    "time"
)

const (
    ImportLayout = "2006-01-02 15:04" // yyyy-MM-dd HH:mm
)

func InTimeRange(time, from, to *time.Time, isAllowNil bool) bool {
    if time == nil {
        return false
    }
    if isAllowNil {
        if from == nil && to == nil {
            return false
        } else if from == nil && to != nil {
            return time.Equal(*to) || time.Before(*to)
        } else if from != nil && to == nil {
            return time.Equal(*from) || time.After(*from)
        } else {
            return (time.Equal(*from) || time.After(*from)) &&
                (time.Equal(*to) || time.Before(*to))
        }
    } else {
        if from == nil || to == nil {
            return false
        } else {
            return (time.Equal(*from) || time.After(*from)) &&
                (time.Equal(*to) || time.Before(*to))
        }
    }
}

func FormatTime(timeObj interface{}, layout string) string {
    switch timeObj.(type) {
    case *time.Time:
        time := timeObj.(*time.Time)
        if time != nil {
            return time.Format(layout)
        }
    case *sql.NullTime:
        time := timeObj.(*sql.NullTime)
        if time != nil && time.Valid {
            return time.Time.Format(layout)
        }
    }
    return ""
}

func ParseTime(layout string, str string) *time.Time {
    time, err := time.Parse(layout, str)
    if err != nil {
        return nil
    } else {
        return &time
    }
}

func IsValidImportTime(timeStr string) bool {
    _, err := time.Parse(ImportLayout, timeStr)
    return err == nil
}

func TimeVNYmdHis(t time.Time) (string, error) {
    t = t.Add(7 * time.Hour)
    const TimeFormat = "20060102150405"
    return t.Format(TimeFormat), nil
}

func SafeUnix(t *time.Time) int64 {
    if t != nil {
        return t.Unix()
    }
    
    return 0
}

func GetSqlNullTime(time *time.Time) *sql.NullTime {
    if time == nil {
        return &sql.NullTime{Valid: false}
    } else {
        return &sql.NullTime{Time: *time, Valid: true}
    }
}
