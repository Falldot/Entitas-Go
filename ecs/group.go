// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type Group interface {
	/////////////////////////////////////////////////////
	//Public:

	Lenght() int
	GetEntities() []*Entity

	OnEntityAdded(func(Group, *Entity, int, Component))
	OnEntityRemoved(func(Group, *Entity, int, Component))
	OnEntityUpdate(func(Group, *Entity, int, Component))

	/////////////////////////////////////////////////////
	//Private:

	has(entity *Entity) bool

	updateEntity(entity *Entity, index int, component Component)

	handleEntity(entity *Entity, index int, component Component)
	handleEntitySilently(entity *Entity)

	addEntitySilently(entity *Entity)
	addEntity(entity *Entity, index int, component Component)

	removeEntitySilently(entity *Entity)
	removeEntity(entity *Entity, index int, component Component)
}

type group struct {
	matcher Matcher

	entities           map[int]*Entity
	entitiesCache      []*Entity
	entitieSingleCache *Entity

	onEntityAdded   onGroup
	onEntityRemoved onGroup
	onEntityUpdate  onGroup
}

/////////////////////////////////////////////////////
//Constructor:
/////////////////////////////////////////////////////

func newGroup(matcher Matcher) Group {
	return &group{
		matcher:  matcher,
		entities: make(map[int]*Entity),
	}
}

/////////////////////////////////////////////////////
//Public:
/////////////////////////////////////////////////////

func (g *group) Lenght() int {
	return len(g.entities)
}

func (g *group) GetEntities() []*Entity {
	if g.entitiesCache == nil {
		g.entitiesCache = make([]*Entity, 0, len(g.entities))
		for _, v := range g.entities {
			g.entitiesCache = append(g.entitiesCache, v)
		}
	}
	return g.entitiesCache
}

func (g *group) GetSingleEntity() *Entity {
	if g.entitieSingleCache == nil {
		count := len(g.entities)
		if count == 1 {
			for _, v := range g.entities {
				g.entitieSingleCache = v
			}
		} else {
			return nil
		}
	}
	return g.entitieSingleCache
}

func (g *group) OnEntityAdded(action func(Group, *Entity, int, Component)) {
	g.onEntityAdded.On(action)
}

func (g *group) OnEntityRemoved(action func(Group, *Entity, int, Component)) {
	g.onEntityRemoved.On(action)
}

func (g *group) OnEntityUpdate(action func(Group, *Entity, int, Component)) {
	g.onEntityUpdate.On(action)
}

/////////////////////////////////////////////////////
//Private:
/////////////////////////////////////////////////////

func (g *group) has(entity *Entity) bool {
	return g.entities[entity.ID()] != nil
}

func (g *group) updateEntity(entity *Entity, index int, component Component) {
	if g.onEntityUpdate != nil {
		g.onEntityUpdate.Execute(g, entity, index, component)
	}
}

func (g *group) handleEntity(entity *Entity, index int, component Component) {
	if g.matcher.matches(entity) {
		g.addEntity(entity, index, component)
	} else {
		g.removeEntity(entity, index, component)
	}
}

func (g *group) handleEntitySilently(entity *Entity) {
	if g.matcher.matches(entity) {
		g.addEntitySilently(entity)
	}
}

func (g *group) addEntitySilently(entity *Entity) {
	g.entities[entity.ID()] = entity
	if g.entitiesCache != nil {
		g.entitiesCache = append(g.entitiesCache, entity)
	}
}

func (g *group) addEntity(entity *Entity, index int, component Component) {
	g.addEntitySilently(entity)
	if g.onEntityAdded != nil {
		g.onEntityAdded.Execute(g, entity, index, component)
	}
}

func (g *group) removeEntitySilently(entity *Entity) {
	delete(g.entities, entity.ID())
	g.entitiesCache = nil
}

func (g *group) removeEntity(entity *Entity, index int, component Component) {
	g.removeEntitySilently(entity)
	if g.onEntityRemoved != nil {
		g.onEntityRemoved.Execute(g, entity, index, component)
	}
}
