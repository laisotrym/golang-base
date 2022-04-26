package config

import (
    "time"
)

func SetTimeZone(timezone *string) {
    if timezone == nil {
        timezone = new(string)
        *timezone = "UTC"
    }

    loc, err := time.LoadLocation(*timezone)
    if err != nil {
        panic(err.Error())
    }
    time.Local = loc
}
