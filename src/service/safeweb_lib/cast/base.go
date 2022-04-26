package safeweb_lib_cast

import (
    "reflect"
)

var stringType = reflect.TypeOf("")
var errUnexpectedType = "non-numeric type could not be converted to %s"
var errConvertType = "cannot convert %v to %s"
