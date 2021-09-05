// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type onComponent []func(*Entity, int, Component)

func (e *onComponent) On(action func(*Entity, int, Component)) {
	*e = append(*e, action)
}

func (e *onComponent) Execute(en *Entity, id int, c Component) {
	for _, fn := range *e {
		fn(en, id, c)
	}
}

type onEntity []func(*Entity)

func (e *onEntity) On(action func(*Entity)) {
	*e = append(*e, action)
}

func (e *onEntity) Execute(en *Entity) {
	for _, fn := range *e {
		fn(en)
	}
}

type onGroup []func(Group, *Entity, int, Component)

func (e *onGroup) On(action func(Group, *Entity, int, Component)) {
	*e = append(*e, action)
}

func (e *onGroup) Execute(g Group, en *Entity, id int, c Component) {
	for _, fn := range *e {
		fn(g, en, id, c)
	}
}

type onGroupCreate []func(Group)

func (e *onGroupCreate) On(action func(Group)) {
	*e = append(*e, action)
}

func (e *onGroupCreate) Execute(g Group) {
	for _, fn := range *e {
		fn(g)
	}
}
