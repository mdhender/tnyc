// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package tnyc

type IslandID int
type CellID int
type CornerID int
type EdgeID int

const NoIslandID IslandID = 0

type Terrain string

type Point struct {
	X float64
	Y float64
}

type World struct {
	Islands []Island
	Cells   []Cell
	Corners []Corner
	Edges   []Edge
}

type Island struct {
	ID      IslandID
	CellIDs []CellID
}

type Cell struct {
	ID        CellID
	IslandID  IslandID
	CornerIDs []CornerID
	Terrain   Terrain
}

type Corner struct {
	ID    CornerID
	Point Point
}

type Edge struct {
	ID        EdgeID
	CornerIDs [2]CornerID
	CellIDs   []CellID
}
