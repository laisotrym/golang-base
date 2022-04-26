package safeweb_lib_error

type GeneralError struct {
    Success                   *Error `yaml:"success"`
    InternalServerError       *Error `yaml:"internalServerError"`
    DataFormatInvalid         *Error `yaml:"dataFormatInvalid"`
    DataInvalid               *Error `yaml:"dataInvalid"`
    MissingRequiredParam      *Error `yaml:"missingRequiredParam"`
    FromTimeExceedToTime      *Error `yaml:"fromTimeExceedToTime"`
    TimeRangeExceedLimitation *Error `yaml:"timeRangeExceedLimitation"`
}

type GeneralStatus struct {
    GeneralError GeneralError `yaml:"general"`
}

var generalStatus *GeneralStatus

func GetGeneralStatus() (GeneralError, error) {
    if generalStatus == nil {
        if err := Load(&generalStatus); err != nil {
            return GeneralError{}, err
        }
    }

    return generalStatus.GeneralError, nil
}
