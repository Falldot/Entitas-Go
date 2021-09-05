// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type componentState []int

func (s *componentState) On(i int, ComponentTotal int) {
	(*s)[i/ComponentTotal] |= 1 << (i % ComponentTotal)
}

func (s *componentState) Off(i int, ComponentTotal int) {
	(*s)[i/ComponentTotal] &^= 1 << (i % ComponentTotal)
}

func (s *componentState) onOff(i int, ComponentTotal int) bool {
	return (*s)[i/ComponentTotal]&(1<<(i%ComponentTotal)) == 0
}
