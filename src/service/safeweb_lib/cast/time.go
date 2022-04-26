package safeweb_lib_cast

import (
    "github.com/araddon/dateparse"
    "github.com/gogo/protobuf/types"
)

func GetStringToProtoTimestamp(s string) (types.Timestamp, error) {
    if t, err := dateparse.ParseLocal(s); err != nil {
        return types.Timestamp{}, err
    } else {
        return types.Timestamp{
            Seconds: t.Unix(),
            Nanos:   int32(t.Nanosecond()),
        }, nil
    }
}
