package safeweb_lib_utils

import (
    "github.com/google/uuid"
)

// NewUUID return new unique Id
func NewUUID() string {
    return uuid.New().String()
}
