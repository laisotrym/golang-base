// Code generated by "stringer -type ActiveStatus -output active_status_string.go"; DO NOT EDIT.

package enum

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Inactive-0]
	_ = x[Active-1]
}

const _ActiveStatus_name = "InactiveActive"

var _ActiveStatus_index = [...]uint8{0, 8, 14}

func (i ActiveStatus) String() string {
	if i < 0 || i >= ActiveStatus(len(_ActiveStatus_index)-1) {
		return "ActiveStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActiveStatus_name[_ActiveStatus_index[i]:_ActiveStatus_index[i+1]]
}
