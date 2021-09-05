// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type componentPool []Component

func (s *componentPool) Push(e Component) {
	*s = append(*s, e)
}

func (s *componentPool) Pop() (Component, bool) {
	lenght := len(*s)
	if lenght > 0 {
		last := lenght - 1
		entity := (*s)[last]
		*s = (*s)[:last]
		return entity, true
	}
	return nil, false
}
