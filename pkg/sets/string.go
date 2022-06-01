package sets

// Set implements a set of strings
type Set struct {
	set map[string]struct{}
}

// NewSet returns a set of given elements
func NewSet(elems ...string) *Set {
	s := &Set{
		set: map[string]struct{}{},
	}
	s.Add(elems...)
	return s
}

// Add adds an element to the set
func (s *Set) Add(elems ...string) {
	for _, e := range elems {
		s.set[e] = struct{}{}
	}
}

// Iter returns an iteratable map
func (s *Set) Iter() map[string]struct{} {
	return s.set
}

// Union executes a set union operatio with given set
func (s *Set) Union(t *Set) {
	for e := range t.Iter() {
		s.Add(e)
	}
}

// IsElem says if the given string is an element of set or not
func (s *Set) IsElem(e string) bool {
	_, ok := s.set[e]
	return ok
}

// Elems returns the array of elements
func (s *Set) Elems() []*string {
	elems := []*string{}
	for e := range s.Iter() {
		tmp := e
		elems = append(elems, &tmp)
	}
	return elems
}
