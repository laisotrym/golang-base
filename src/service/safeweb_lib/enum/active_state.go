//go:generate stringer -type ActiveState -output active_state_string.go
package safeweb_lib_num

// ActiveState is generated type for enum 'active_state'
type ActiveState int32

const (
    Inactive ActiveState = iota
    Active
)
