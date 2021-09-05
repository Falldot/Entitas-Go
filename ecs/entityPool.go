// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type entityPool []int

func (s *entityPool) Push(e int) {
	*s = append(*s, e)
}

func (s *entityPool) Pop() (int, bool) {
	lenght := len(*s)
	if lenght > 0 {
		last := lenght - 1
		entity := (*s)[last]
		*s = (*s)[:last]
		return entity, true
	}
	return 0, false
}
