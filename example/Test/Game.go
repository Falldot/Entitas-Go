package main

//go:generate go run github.com/Falldot/Entitas-Go

type Position struct {
	X, Y float64
}

type Direction struct {
	X, Y float64
}

type Speed int

type Health int

type Sprite struct {
	tag  string
	W, H int
}
