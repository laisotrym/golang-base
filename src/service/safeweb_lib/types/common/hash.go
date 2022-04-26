package safeweb_lib_types_common

type HashKey struct {
    Key   string
    Field string
}

func NewHashKey(key string, field string) *HashKey {
    return &HashKey{Key: key, Field: field}
}
