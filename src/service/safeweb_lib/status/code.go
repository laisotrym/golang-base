package safeweb_lib_status

import (
    rpc "google.golang.org/grpc/codes"
)

type Code int

const (
    Success                 Code = 200
    Unauthorized            Code = 410
    InvalidFormat           Code = 411
    InvalidData             Code = 412
    InvalidOrder            Code = 420
    OrderFailed             Code = 421
    InvalidTransaction      Code = 430
    FailedTransaction       Code = 431
    DoubtingTransaction     Code = 432
    TimeoutTransaction      Code = 433
    CanceledTransaction     Code = 434
    WaitingBuyerTransaction Code = 435
    ProcessingTransaction   Code = 436
    RefundedTransaction     Code = 437
    InternalError           Code = 500
    Maintenance             Code = 503
    Unspecified             Code = 600
)

var codeToStr = map[Code]string{
    Success:                 "200",
    Unauthorized:            "410",
    InvalidFormat:           "411",
    InvalidData:             "412",
    InvalidOrder:            "420",
    OrderFailed:             "421",
    InvalidTransaction:      "430",
    FailedTransaction:       "431",
    DoubtingTransaction:     "432",
    TimeoutTransaction:      "433",
    CanceledTransaction:     "434",
    WaitingBuyerTransaction: "435",
    ProcessingTransaction:   "436",
    RefundedTransaction:     "437",
    InternalError:           "500",
    Maintenance:             "503",
    Unspecified:             "600",
}

var codeToRpcCode = map[Code]rpc.Code{
    Success:                 rpc.OK,
    Unauthorized:            rpc.InvalidArgument,
    InvalidFormat:           rpc.InvalidArgument,
    InvalidData:             rpc.InvalidArgument,
    InvalidOrder:            rpc.InvalidArgument,
    OrderFailed:             rpc.InvalidArgument,
    InvalidTransaction:      rpc.InvalidArgument,
    FailedTransaction:       rpc.InvalidArgument,
    InternalError:           rpc.Internal,
    Maintenance:             rpc.InvalidArgument,
    Unspecified:             rpc.InvalidArgument,
    DoubtingTransaction:     rpc.InvalidArgument,
    TimeoutTransaction:      rpc.InvalidArgument,
    CanceledTransaction:     rpc.InvalidArgument,
    WaitingBuyerTransaction: rpc.InvalidArgument,
    ProcessingTransaction:   rpc.InvalidArgument,
    RefundedTransaction:     rpc.InvalidArgument,
}

func (c Code) ToRpcCode() rpc.Code {
    return codeToRpcCode[c]
}

func (c Code) String() string {
    return codeToStr[c]
}

func (c Code) ToRPCCode() rpc.Code {
    return codeToRpcCode[c]
}
