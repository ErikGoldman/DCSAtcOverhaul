package commands

import "math/rand"

// Generator interface for random numbers
type Random interface {
	IntN(n int) int
}

// Real implementation using math/rand
type RealGenerator struct{}

func (r *RealGenerator) IntN(n int) int {
	return rand.Intn(n)
}

// Mock implementation for testing
type MockGenerator struct {
	Values []int
	index  int
}

func (m *MockGenerator) IntN(n int) int {
	if m.index >= len(m.Values) {
		panic("Ran out of random values")
	}
	value := m.Values[m.index]
	m.index++
	return value
}

func (m *MockGenerator) Reset() {
	m.Values = []int{}
	m.index = 0
}

func (m *MockGenerator) ResetTo(intValues []int) {
	m.Values = intValues
	m.index = 0
}

func (m *MockGenerator) PushInt(n int) {
	m.Values = append(m.Values, n)
}
