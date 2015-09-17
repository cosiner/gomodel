package store

type (
	Strings struct {
		Values []string
	}

	Ints struct {
		Values []int
	}

	Bools struct {
		Values []bool
	}

	KVs struct {
		Keys   []string
		Values []string
	}
)

func (s *Strings) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]string, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *Strings) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *Strings) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Values[index]
}

func (s *Strings) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]string, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of Strings")
}

func (s *Strings) Clear() {
	if s.Values != nil {
		s.Values = s.Values[:0]
	}
}

func (s *Ints) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]int, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *Ints) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *Ints) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Values[index]
}

func (s *Ints) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]int, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of Ints")
}

func (s *Ints) Clear() {
	if s.Values != nil {
		s.Values = s.Values[:0]
	}
}
func (s *Bools) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]bool, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *Bools) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *Bools) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Values[index]
}

func (s *Bools) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]bool, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of Bools")
}

func (s *Bools) Clear() {
	if s.Values != nil {
		s.Values = s.Values[:0]
	}
}

func (s *KVs) Init(size int) {
	if cap(s.Values) < size {
		s.Keys = make([]string, size)
		s.Values = make([]string, size)
	} else {
		s.Keys = s.Keys[:size]
		s.Values = s.Values[:size]
	}
}

func (s *KVs) Final(size int) {
	s.Keys = s.Keys[:size]
	s.Values = s.Values[:size]
}

func (s *KVs) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Keys[index]
	ptrs[1] = &s.Values[index]
}

func (s *KVs) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		keys := make([]string, 2*c)
		copy(keys, s.Keys)
		s.Keys = keys
		values := make([]string, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Keys = s.Keys[:c]
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of KVs")
}
