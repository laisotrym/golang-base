package safeweb_lib_types_set

const maxInt = int(^uint(0) >> 1)

var exists = struct{}{}

type StringSet struct {
    m map[string]struct{}
}

func (s *StringSet) GetRaw() []string {
    var result []string
    for key := range s.m {
        result = append(result, key)
    }
    return result
}

func NewStringSet(ts ...string) *StringSet {
    s := &StringSet{m: map[string]struct{}{}}
    s.Add(ts...)
    return s
}

func (s *StringSet) Add(items ...string) {
    for _, item := range items {
        s.m[item] = exists
    }
}

func (s *StringSet) Remove(items ...string) {
    for _, item := range items {
        delete(s.m, item)
    }
}

func (s *StringSet) Contains(value string) bool {
    _, c := s.m[value]
    return c
}

func (s *StringSet) Length() int {
    return len(s.m)
}

func (s *StringSet) IsEqual(t *StringSet) bool {
    // return false if they are no the same size
    if s.Length() != t.Length() {
        return false
    }
    
    equal := true
    t.Each(func(item string) bool {
        _, equal = s.m[item]
        return equal
    })
    
    return equal
}

func (s *StringSet) Each(f func(item string) bool) {
    for item := range s.m {
        if !f(item) {
            break
        }
    }
}

func Intersection(sets ...*StringSet) *StringSet {
    minPos := -1
    minSize := maxInt
    for i, set := range sets {
        if l := set.Size(); l < minSize {
            minSize = l
            minPos = i
        }
    }
    if minSize == maxInt || minSize == 0 {
        return NewStringSet()
    }
    
    t := sets[minPos].Copy()
    for i, set := range sets {
        if i == minPos {
            continue
        }
        for item := range t.m {
            if _, has := set.m[item]; !has {
                delete(t.m, item)
            }
        }
    }
    return t
}

// Copy returns a new Set with a copy of s.
func (s *StringSet) Copy() *StringSet {
    u := NewWithSize(s.Size())
    for item := range s.m {
        u.m[item] = exists
    }
    return u
}

// NewWithSize creates a new Set and gives make map a size hint.
func NewWithSize(size int) *StringSet {
    return &StringSet{make(map[string]struct{}, size)}
}

func (s *StringSet) Size() int {
    return len(s.m)
}
