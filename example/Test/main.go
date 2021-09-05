package main

import (
	ecs "test/Entitas"
)

func main() {
	contexts := ecs.CreateContexts()
	game := contexts.Game()

	// System registration
	systems := ecs.CreateSystemPool()
	systems.Add(&Translate{})
	systems.Add(&ReactiveTranslate{})

	// Create entity
	player := game.CreateEntity()

	// Add component
	player.AddPosition(10, 30)
	player.AddDirection(0, 0)
	player.AddSpeed(5)

	// Remove component
	player.RemoveSpeed()

	// Replace component
	player.ReplacePosition(30, 10)

	// On or Off component
	player.OffDirection()
	player.OnDirection()

	// Destroy entity
	//player.Destroy()

	// GameLoop
	systems.Init(contexts)
	systems.Execute()
	systems.Clean()
	systems.Exit(contexts)
}

type Translate struct {
	group ecs.Group
}

func (s *Translate) Initer(contexts ecs.Contexts) {
	game := contexts.Game()
	s.group = game.Group(ecs.NewMatcher().AllOf(ecs.Position))
}

func (s *Translate) Executer() {
	for _, e := range s.group.GetEntities() {
		pos := e.GetPosition()
		e.ReplacePosition(pos.X+10, pos.X+10)
	}
}

type ReactiveTranslate struct {
}

func (s *ReactiveTranslate) Trigger(contexts ecs.Contexts) ecs.Collector {
	game := contexts.Game()
	return game.Collector(ecs.NewMatcher().AllOf(ecs.Position)).OnUpdate().OnAdd()
}

func (s *ReactiveTranslate) Filter(entity *ecs.Entity) bool {
	return entity.Has(ecs.Position)
}

func (s *ReactiveTranslate) Executer(entities []*ecs.Entity) {
	for _, e := range entities {
		pos := e.GetPosition()
		e.ReplacePosition(pos.X+10, pos.X+10)
	}
}
