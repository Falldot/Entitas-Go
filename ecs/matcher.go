// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type Matcher interface {

	/////////////////////////////////////////////////////
	//Public:
	AllOf(indices ...int) Matcher
	AnyOf(indices ...int) Matcher
	NoneOf(indices ...int) Matcher

	/////////////////////////////////////////////////////
	//Private:

	getIndices() []int
	matches(entity *Entity) bool
	getHashCode() int
}

type matcher struct {
	allOfIndices  []int
	anyOfIndices  []int
	noneOfIndices []int

	indices []int

	hash        int
	isHashCache bool
}

/////////////////////////////////////////////////////
// Constructor:
/////////////////////////////////////////////////////

func NewMatcher() Matcher {
	return &matcher{}
}

/////////////////////////////////////////////////////
// Public:
/////////////////////////////////////////////////////

func (m *matcher) AllOf(indices ...int) Matcher {
	m.allOfIndices = append(m.allOfIndices, indices...)
	m.indices = append(m.indices, indices...)
	return m
}

func (m *matcher) AnyOf(indices ...int) Matcher {
	m.anyOfIndices = append(m.anyOfIndices, indices...)
	m.indices = append(m.indices, indices...)
	return m
}

func (m *matcher) NoneOf(indices ...int) Matcher {
	m.noneOfIndices = append(m.noneOfIndices, indices...)
	m.indices = append(m.indices, indices...)
	return m
}

/////////////////////////////////////////////////////
// Private:
/////////////////////////////////////////////////////

func (m *matcher) getIndices() []int {
	return m.indices
}

func (m *matcher) matches(entity *Entity) bool {
	return (m.allOfIndices == nil || entity.Has(m.allOfIndices...)) &&
		(m.anyOfIndices == nil || entity.Any(m.anyOfIndices...)) &&
		(m.noneOfIndices == nil || !entity.Has(m.noneOfIndices...))
}

func (m *matcher) getHashCode() int {
	if !m.isHashCache {
		var hash int
		hash = applyHash(hash, m.allOfIndices, 3, 53)
		hash = applyHash(hash, m.anyOfIndices, 307, 367)
		hash = applyHash(hash, m.noneOfIndices, 647, 683)
		m.hash = hash
		m.isHashCache = true
	}
	return m.hash
}

func applyHash(hash int, indices []int, i1 int, i2 int) int {
	if indices != nil {
		for _, v := range indices {
			hash ^= v * i1
		}
		hash ^= len(indices) * i2
	}
	return hash
}
