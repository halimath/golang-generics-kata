package seq

import (
	"bufio"
	"errors"
)

// A sentinel error value used to describe the end of a Seq.
var Exhausted = errors.New("sequence exhausted")

// Seq defines the fundamental interface for sequences of data. T describes
// the element's type.
type Seq[T any] interface {
	// Next retrieves the next element from the Seq and returns it. If no
	// element is available, Exhausted will be returned as the error and T
	// is undefined (usually the default value). Any other non-nil error
	// denotes an error of the underlying datasource (i.e. i/o error).
	Next() (T, error)
}

// sliceSeq implements Seq for any type of go slice.
type sliceSeq[E any] struct {
	s   []E
	idx int
}

func (s *sliceSeq[E]) Next() (e E, err error) {
	if s.idx >= len(s.s) {
		return e, Exhausted
	}

	e = s.s[s.idx]

	s.idx++

	return
}

// FromSlice creates a Seq with s used as the underlying source. This operation
// does not copy s. If s changes the behaviour of the returned Seq is undefined.
func FromSlice[E any, S ~[]E](s S) Seq[E] {
	return &sliceSeq[E]{s: s}
}

type scannerSeq struct {
	scanner *bufio.Scanner
}

func (s scannerSeq) Next() ([]byte, error) {
	if !s.scanner.Scan() {
		return nil, Exhausted
	}

	if s.scanner.Err() != nil {
		return nil, s.scanner.Err()
	}

	return s.scanner.Bytes(), nil
}

func Scanner(s *bufio.Scanner) Seq[[]byte] {
	return &scannerSeq{scanner: s}
}

// Integer defines a constraint interface which contains all integral data types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// intRange implements an finite or infinite range of integral values of type
// E provided as a Seq.
type intRange[E Integer] struct {
	v, inc E
	end    *E
}

func (r *intRange[E]) Next() (e E, err error) {
	if r.end != nil && r.v+r.inc >= *r.end {
		return e, Exhausted
	}

	e = r.v
	r.v += r.inc

	return
}

// InfiniteIntRange creates a Seq that generates values of type E starting at
// start and incrementing by inc. The range generates values infinitely.
func InfiniteIntRange[E Integer](start, inc E) Seq[E] {
	return &intRange[E]{v: start, inc: inc}
}

// takeSeq implements a Seq that picks at most limit values from an underlying
// Seq.
type takeSeq[E any] struct {
	current, limit int
	from           Seq[E]
}

func (t *takeSeq[E]) Next() (e E, err error) {
	if t.current >= t.limit {
		return e, Exhausted
	}

	e, err = t.from.Next()
	t.current++

	return
}

// Take takes at most count elements from s.
func Take[E any](s Seq[E], count int) Seq[E] {
	return &takeSeq[E]{limit: count, from: s}
}

// ForEach applies f to each element from s. It returns nil on success or
// a non-nil error if s returns any non-nil value that is not Exhausted.
func ForEach[E any](s Seq[E], f func(E)) error {
	for {
		e, err := s.Next()
		if errors.Is(err, Exhausted) {
			return nil
		}
		if err != nil {
			return err
		}

		f(e)
	}
}

// mapSeq applies mapper to each element from s.
type mapSeq[I, O any] struct {
	s      Seq[I]
	mapper func(I) O
}

func (m mapSeq[I, O]) Next() (o O, err error) {
	i, err := m.s.Next()
	if err != nil {
		return
	}

	return m.mapper(i), nil
}

// Map applies m to each element from s.
func Map[I, O any](s Seq[I], m func(I) O) Seq[O] {
	return &mapSeq[I, O]{s: s, mapper: m}
}

// mapESeq applies mapper to each element from s.
type mapESeq[I, O any] struct {
	s      Seq[I]
	mapper func(I) (O, error)
}

func (m mapESeq[I, O]) Next() (o O, err error) {
	i, err := m.s.Next()
	if err != nil {
		return
	}

	return m.mapper(i)
}

func MapE[I, O any](s Seq[I], m func(I) (O, error)) Seq[O] {
	return &mapESeq[I, O]{s: s, mapper: m}
}

// Reduce applies reducer to each element of s in turn, passing in the previous
// result returned from the reducer. It eventually returns the last reducer's
// result. initial may either contain 0 or 1 element. If one is given it is
// used as the initial value passed to the first invocation of the reducer.
// Otherwise R's default value is used. If len(initial) > 1 Reduce panics.
func Reduce[I, R any](s Seq[I], reducer func(R, I) R, initial ...R) (r R, err error) {
	if len(initial) > 1 {
		panic("more then one initial value given")
	}

	if len(initial) == 1 {
		r = initial[0]
	}

	var i I
	for {
		i, err = s.Next()
		if errors.Is(err, Exhausted) {
			return r, nil
		}
		if err != nil {
			return
		}

		r = reducer(r, i)
	}
}

// Predicate defines a predicate function for values of type E.
type Predicate[E any] func(E) bool

// filterSeq only returns those elements from s where p(s) returns true.
type filterSeq[E any] struct {
	s Seq[E]
	p Predicate[E]
}

func (f *filterSeq[E]) Next() (e E, err error) {
	e, err = f.s.Next()
	if err != nil {
		return
	}

	if f.p(e) {
		return
	}

	return f.Next()
}

// Filter returns a Seq that pulls values from s and skips those where p(val) is
// false.
func Filter[E any](s Seq[E], p Predicate[E]) Seq[E] {
	return &filterSeq[E]{s: s, p: p}
}
