// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
package ecs

type EntityBase interface {
	/////////////////////////////////////////////////////
	//Public:

	CreateEntity() *Entity
	Group(Matcher) Group
	Collector(Matcher) Collector

	OnGroupCreated(func(Group))

	OnEntityCreated(func(*Entity))
	OnEntityWillBeDestroyed(func(*Entity))
	OnEntityDestroyed(func(*Entity))

	/////////////////////////////////////////////////////
	//Private:
	destroyEntity(e *Entity)
	removeComponent(e *Entity, tp int, c Component)
}

type entityBase struct {
	componentTotal int

	entities []*Entity
	entityPool
	componentPool []componentPool

	groups      map[int]Group
	groupsIndex [][]Group

	onGroupCreated onGroupCreate

	onEntityCreated         onEntity
	onEntityWillBeDestroyed onEntity
	onEntityDestroyed       onEntity
}

/////////////////////////////////////////////////////
// Constructor:
/////////////////////////////////////////////////////

func CreateEntityBase(componentTotal int) EntityBase {
	return &entityBase{
		componentTotal: componentTotal,
		groupsIndex:    make([][]Group, componentTotal),
		componentPool:  make([]componentPool, componentTotal),
		groups:         make(map[int]Group),
	}
}

/////////////////////////////////////////////////////
// Public:
/////////////////////////////////////////////////////

func (eb *entityBase) CreateEntity() *Entity {
	var entity *Entity

	if id, ok := eb.entityPool.Pop(); ok {
		entity = eb.entities[id]
	} else {
		entity = createEntity(len(eb.entities), &eb.componentPool, eb.componentTotal)
		eb.entities = append(eb.entities, entity)
		entity.OnEntityDestroy(eb.destroyEntity)
		entity.OnComponentRemoved(eb.removeComponent)
		entity.OnComponentRemoved(eb.updateGroupsComponentRemoved)
		entity.OnComponentOff(eb.updateGroupsComponentRemoved)
		entity.OnComponentAdd(eb.updateGroupsComponentAdded)
		entity.OnComponentReplaced(eb.updateGroupsComponentReplaced)
	}

	if eb.onEntityCreated != nil {
		eb.onEntityCreated.Execute(entity)
	}

	return entity
}

func (eb *entityBase) OnGroupCreated(action func(Group)) {
	eb.onGroupCreated.On(action)
}

func (eb *entityBase) OnEntityCreated(action func(*Entity)) {
	eb.onEntityCreated.On(action)
}

func (eb *entityBase) OnEntityWillBeDestroyed(action func(*Entity)) {
	eb.onEntityWillBeDestroyed.On(action)
}

func (eb *entityBase) OnEntityDestroyed(action func(*Entity)) {
	eb.onEntityDestroyed.On(action)
}

func (eb *entityBase) Group(matcher Matcher) Group {
	hash := matcher.getHashCode()

	if g, ok := eb.groups[hash]; ok {
		return g
	}

	g := newGroup(matcher)

	for _, v := range eb.entities {
		g.handleEntitySilently(v)
	}

	eb.groups[hash] = g

	for _, v := range matcher.getIndices() {
		eb.groupsIndex[v] = append(eb.groupsIndex[v], g)
	}

	if eb.onGroupCreated != nil {
		eb.onGroupCreated.Execute(g)
	}

	return g
}

func (eb *entityBase) Collector(matcher Matcher) Collector {
	g := eb.Group(matcher)
	return newCollector(g)
}

/////////////////////////////////////////////////////
// Private:
/////////////////////////////////////////////////////

func (eb *entityBase) updateGroupsComponentRemoved(entity *Entity, index int, component Component) {
	entity.componentState.Off(index, eb.componentTotal)
	for _, v := range eb.groupsIndex[index] {
		v.handleEntity(entity, index, component)
	}
}

func (eb *entityBase) updateGroupsComponentAdded(entity *Entity, index int, component Component) {
	entity.componentState.On(index, eb.componentTotal)
	for _, v := range eb.groupsIndex[index] {
		v.handleEntity(entity, index, component)
	}
}

func (eb *entityBase) updateGroupsComponentReplaced(entity *Entity, index int, component Component) {
	for _, v := range eb.groupsIndex[index] {
		v.updateEntity(entity, index, component)
	}
}

func (eb *entityBase) destroyEntity(e *Entity) {
	e.onEntityDestroy = e.onEntityDestroy[:1]
	if eb.onEntityWillBeDestroyed != nil {
		eb.onEntityWillBeDestroyed.Execute(e)
	}
	eb.entityPool.Push(e.id)
	if eb.onEntityDestroyed != nil {
		eb.onEntityDestroyed.Execute(e)
	}
}

func (eb *entityBase) removeComponent(e *Entity, tp int, c Component) {
	eb.componentPool[tp].Push(c)
	e.components[tp] = nil
	e.componentState.Off(tp, eb.componentTotal)
}
