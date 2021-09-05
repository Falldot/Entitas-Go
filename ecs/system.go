// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type Initer interface {
	Initer(Contexts)
}

type Executer interface {
	Executer()
}

type Cleaner interface {
	Cleaner()
}

type Reactive interface {
	Trigger(contexts Contexts) Collector
	Filter(entity *Entity) bool
	Executer(entities []*Entity)
}

type Exiter interface {
	Exiter(Contexts)
}

type Systems interface {
	Init(contexts Contexts)
	Execute()
	Clean()
	Exit(contexts Contexts)

	Add(action interface{})
}

type systems struct {
	initers   []Initer
	executers []Executer
	cleaners  []Cleaner
	exiters   []Exiter

	reactives  []Reactive
	collectors []Collector
}

/////////////////////////////////////////////////////
// Constructor:
/////////////////////////////////////////////////////

func CreateSystemPool() Systems {
	return &systems{}
}

/////////////////////////////////////////////////////
// Public:
/////////////////////////////////////////////////////

func (s *systems) Init(contexts Contexts) {
	for _, v := range s.initers {
		v.Initer(contexts)
	}
	s.collectors = make([]Collector, len(s.reactives))
	for i, v := range s.reactives {
		s.collectors[i] = v.Trigger(contexts)
	}
}

func (s *systems) Execute() {
	for _, v := range s.executers {
		v.Executer()
	}
	for i, v := range s.reactives {
		if s.collectors[i] != nil {
			var buffer []*Entity
			entities := s.collectors[i].Entities()
			if len(entities) != 0 {
				for _, e := range entities {
					if v.Filter(e) {
						buffer = append(buffer, e)
					}
				}

				s.collectors[i].Clear()

				if len(buffer) != 0 {
					v.Executer(buffer)
				}
			}
		}
	}
}

func (s *systems) Clean() {
	for _, v := range s.cleaners {
		v.Cleaner()
	}
}

func (s *systems) Exit(contexts Contexts) {
	for _, v := range s.exiters {
		v.Exiter(contexts)
	}
}

func (s *systems) Add(action interface{}) {
	switch action := action.(type) {
	case Reactive:
		s.reactives = append(s.reactives, action)
	}

	switch action := action.(type) {
	case Initer:
		s.initers = append(s.initers, action)
	}

	switch action := action.(type) {
	case Executer:
		s.executers = append(s.executers, action)
	}

	switch action := action.(type) {
	case Cleaner:
		s.cleaners = append(s.cleaners, action)
	}

	switch action := action.(type) {
	case Exiter:
		s.exiters = append(s.exiters, action)
	}
}
