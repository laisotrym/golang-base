//go:generate stringer -type ActiveStatus -output active_status_string.go
package enum

// ActiveStatus is generated type for enum 'active_status'
type ActiveStatus int32

const (
    Inactive ActiveStatus = iota
    Active
)
