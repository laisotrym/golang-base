package safeweb_lib_error

type ClientError struct {
    ClientExisted                *Error `yaml:"clientExisted"`
    ClientNotExisted             *Error `yaml:"clientNotExisted"`
    ClientCodeDuplicated         *Error `yaml:"clientCodeDuplicated"`
    ClientMerchantBackendInvalid *Error `yaml:"clientMerchantBackendInvalid"`
}

type ClientStatus struct {
    ClientError ClientError `yaml:"client"`
}

var clientStatus *ClientStatus

func GetClientStatus() (ClientError, error) {
    if clientStatus == nil {
        if err := Load(&clientStatus); err != nil {
            return ClientError{}, err
        }
    }

    return clientStatus.ClientError, nil
}
