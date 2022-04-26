package safeweb_lib_constant

import (
    "math"
)

const (
    MaxUint = ^uint(0)
    MinUint = 0
    MaxInt  = int(MaxUint >> 1)
    MinInt  = -MaxInt - 1
    
    MaxUint8 = ^uint8(0)
    MinUint8 = 0
    MaxInt8  = int8(MaxUint8 >> 1)
    MinInt8  = -MaxInt8 - 1
    
    MaxUint16 = ^uint16(0)
    MinUint16 = 0
    MaxInt16  = int16(MaxUint16 >> 1)
    MinInt16  = -MaxInt16 - 1
    
    MaxUint32 = ^uint32(0)
    MinUint32 = 0
    MaxInt32  = int32(MaxUint32 >> 1)
    MinInt32  = -MaxInt32 - 1
    
    MaxUint64 = ^uint64(0)
    MinUint64 = 0
    MaxInt64  = int64(MaxUint64 >> 1)
    MinInt64  = -MaxInt64 - 1
    
    MaxFloat32 = math.MaxFloat32
    MinFloat32 = math.SmallestNonzeroFloat32
    
    MaxFloat64 = math.MaxFloat64
    MinFloat64 = math.SmallestNonzeroFloat64
)
