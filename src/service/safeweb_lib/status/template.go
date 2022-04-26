package safeweb_lib_status

import "fmt"

type Template string

type TemplateCode int

const (
    // 200
    DefaultSuccessTemplate TemplateCode = 2000
    // 410
    DefaultUnauthorizedTemplate TemplateCode = 4100
    IPNotAllowedAccessTemplate  TemplateCode = 4101
    WrongChecksumTemplate       TemplateCode = 4102
    // 411
    DefaultInvalidFormatTemplate     TemplateCode = 4110
    OtherInvalidFormatTemplate       TemplateCode = 4111
    OneArgumentInvalidFormatTemplate TemplateCode = 4112
    // 412
    DefaultInvalidDataTemplate             TemplateCode = 4120
    NotExistInvalidDataTemplate            TemplateCode = 4121
    ExistInvalidDataTemplate               TemplateCode = 4122
    NonActiveInvalidDataTemplate           TemplateCode = 4123
    NotExistOrNonActiveInvalidDataTemplate TemplateCode = 4124
    InvalidDataTemplate                    TemplateCode = 4125
    MisMatchDataTemplate                   TemplateCode = 4126
    OtherInvalidDataTemplate               TemplateCode = 4127
    // 420
    DefaultInvalidOrderTemplate TemplateCode = 4200
    // 421
    DefaultOrderFailedTemplate TemplateCode = 4210
    // 430
    DefaultInvalidTransactionTemplate TemplateCode = 4300
    OtherInvalidTransactionTemplate   TemplateCode = 4301
    // 431
    DefaultFailedTransactionTemplate TemplateCode = 4310
    // 500
    DefaultInternalErrorTemplate TemplateCode = 5000
    CustomInternalErrorTemplate  TemplateCode = 5001
    // 503
    DefaultMaintenanceTemplate TemplateCode = 5030
    // 600
    DefaultUnspecifiedTemplate TemplateCode = 6000
)

var templateCodeToTemplate = map[TemplateCode]Template{
    DefaultSuccessTemplate: "%s thành công",
    // 410
    DefaultUnauthorizedTemplate: "Không thể xác thực chữ ký hoặc thông tin",
    IPNotAllowedAccessTemplate:  "IP không được phép truy cập",
    WrongChecksumTemplate:       "Sai checksum",
    // 411
    DefaultInvalidFormatTemplate:     "%s sai định dạng, định dạng đúng là %s",
    OneArgumentInvalidFormatTemplate: "%s sai định dạng",
    OtherInvalidFormatTemplate:       "%s",
    // 412
    DefaultInvalidDataTemplate:             "%s không hợp lệ, dữ liệu hợp lệ là %s",
    NotExistInvalidDataTemplate:            "%s không tồn tại",
    ExistInvalidDataTemplate:               "%s đã tồn tại",
    NonActiveInvalidDataTemplate:           "%s chưa kích hoạt",
    NotExistOrNonActiveInvalidDataTemplate: "%s không tồn tại hoặc chưa kích hoạt",
    InvalidDataTemplate:                    "%s không hợp lệ",
    MisMatchDataTemplate:                   "%s, %s không cùng %s",
    OtherInvalidDataTemplate:               "%s",
    // 420
    DefaultInvalidOrderTemplate: "Đơn hàng %s không hợp lệ do %s",
    // 421
    DefaultOrderFailedTemplate: "Xử lý đơn hàng %s không thành công do %s",
    // 430
    DefaultInvalidTransactionTemplate: "Giao dịch không hợp lệ do %s",
    OtherInvalidTransactionTemplate:   "%s",
    // 431
    DefaultFailedTransactionTemplate: "Giao dịch %s thất bại do %s",
    // 500
    DefaultInternalErrorTemplate: "Lỗi hệ thống",
    CustomInternalErrorTemplate:  "%s",
    // 503
    DefaultMaintenanceTemplate: "Hệ thống %s đang bảo trì",
    // 600
    DefaultUnspecifiedTemplate: "Lỗi ngoài danh mục mô tả (Các lỗi không có trong danh sách mã lỗi đã liệt kê)",
}

var templateToCode = map[TemplateCode]Code{
    DefaultSuccessTemplate: Success,
    
    DefaultUnauthorizedTemplate: Unauthorized,
    WrongChecksumTemplate:       Unauthorized,
    IPNotAllowedAccessTemplate:  Unauthorized,
    
    DefaultInvalidFormatTemplate:           InvalidFormat,
    OneArgumentInvalidFormatTemplate:       InvalidFormat,
    OtherInvalidFormatTemplate:             InvalidFormat,
    DefaultInvalidDataTemplate:             InvalidData,
    NotExistInvalidDataTemplate:            InvalidData,
    ExistInvalidDataTemplate:               InvalidData,
    NonActiveInvalidDataTemplate:           InvalidData,
    OtherInvalidDataTemplate:               InvalidData,
    NotExistOrNonActiveInvalidDataTemplate: InvalidData,
    InvalidDataTemplate:                    InvalidData,
    MisMatchDataTemplate:                   InvalidData,
    
    DefaultInvalidOrderTemplate: InvalidOrder,
    
    DefaultOrderFailedTemplate: OrderFailed,
    
    DefaultInvalidTransactionTemplate: InvalidTransaction,
    OtherInvalidTransactionTemplate:   InvalidTransaction,
    
    DefaultFailedTransactionTemplate: FailedTransaction,
    
    DefaultInternalErrorTemplate: InternalError,
    CustomInternalErrorTemplate:  InternalError,
    
    DefaultMaintenanceTemplate: Maintenance,
    
    DefaultUnspecifiedTemplate: Unspecified,
}

func (t TemplateCode) ToTemplate() Template {
    return templateCodeToTemplate[t]
}

func (t TemplateCode) With(params ...interface{}) *Trace {
    c, ok := templateToCode[t]
    if !ok {
        return NewErrorTrace(Unspecified, DefaultUnspecifiedTemplate.ToTemplate().String())
    }
    
    if len(params) == 0 {
        return NewErrorTrace(c, t.ToTemplate().String())
    }
    
    return NewErrorTrace(c, fmt.Sprintf(t.ToTemplate().String(), params...))
}

func (t Template) String() string {
    return string(t)
}

func (t TemplateCode) FormatString(params ...interface{}) string {
    if len(params) == 0 {
        return t.ToTemplate().String()
    }
    return fmt.Sprintf(t.ToTemplate().String(), params...)
}
