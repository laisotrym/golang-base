package safeweb_lib_types_common

type Pair struct {
    values [2]interface{}
}

func MakePair(k, v interface{}) Pair {
    return Pair{values: [2]interface{}{k, v}}
}

func (p Pair) Get(i int) interface{} {
    return p.values[i]
}
