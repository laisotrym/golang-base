//go:generate stringer -type RedisType;RedisZAddType -output redis_string.go
package safeweb_lib_constant

type RedisType int

const (
    RedisTypeStrings RedisType = iota
    RedisTypeLists
    RedisTypeSets
    RedisTypeHashes
    RedisTypeSortedSets
    RedisTypeBitmaps
    RedisTypeHyperLogLogs
)

type RedisZAddType int

const (
    RedisZAddTypeXN RedisZAddType = iota
    RedisZAddTypeNX
    RedisZAddTypeXX
)
