package safeweb_lib_types_set

var intExists = struct{}{}

type IntSet struct {
    m map[int]struct{}
}

func NewIntSet(ts ...int) *IntSet {
    s := &IntSet{m: map[int]struct{}{}}
    s.Add(ts...)
    return s
}

func (s *IntSet) Add(items ...int) {
    for _, item := range items {
        s.m[item] = intExists
    }
}

func (s *IntSet) Values() []int {
    values := make([]int, len(s.m))
    idx := 0
    for key := range s.m {
        values[idx] = key
        idx++
    }
    return values
}

func (s *IntSet) Remove(items ...int) {
    for _, item := range items {
        delete(s.m, item)
    }
}

func (s *IntSet) Contains(value int) bool {
    _, c := s.m[value]
    return c
}

func (s *IntSet) Length() int {
    return len(s.m)
}

func (s *IntSet) IsEqual(t *IntSet) bool {
    // return false if they are no the same size
    if s.Length() != t.Length() {
        return false
    }
    
    equal := true
    t.Each(func(item int) bool {
        _, equal = s.m[item]
        return equal
    })
    
    return equal
}

func (s *IntSet) Each(f func(item int) bool) {
    for item := range s.m {
        if !f(item) {
            break
        }
    }
}
