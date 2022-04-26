package safeweb_lib_utils

import (
	"github.com/chilts/sid"
)

// NewShortID return new short Id as token
func NewShortID() string {
	return sid.Id()
}
