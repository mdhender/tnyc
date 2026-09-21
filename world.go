// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package tnyc

type IslandID int
type ProvinceID int
type CornerID int
type EdgeID int

const NoIslandID IslandID = -1

type Terrain string
type ElevationBand string
type HeatBand string
type MoistureBand string

type Point struct {
	X float64
	Y float64
}

type Bounds struct {
	Minimum Point
	Maximum Point
}

type World struct {
	Generation Generation
	Bounds     Bounds
	Islands    []Island
	Provinces  []Province
	Corners    []Corner
	Edges      []Edge
}

type Island struct {
	ID          IslandID
	ProvinceIDs []ProvinceID
}

type Province struct {
	ID            ProvinceID
	IslandID      IslandID
	Center        Point
	CornerIDs     []CornerID
	Terrain       Terrain
	Elevation     float64
	ElevationBand ElevationBand
	Relief        float64
	Heat          float64
	HeatBand      HeatBand
	Moisture      float64
	MoistureBand  MoistureBand
}

type Corner struct {
	ID    CornerID
	Point Point
}

type Edge struct {
	ID          EdgeID
	CornerIDs   [2]CornerID
	ProvinceIDs []ProvinceID
	Elevation   float64
}

type Generation struct {
	Config GenerationConfig
	Result GenerationResult
}

type GenerationConfig struct {
	Seed               string
	ProvinceCount      int
	IslandCount        int
	AspectRatio        string
	OceanFraction      float64
	EdgeBarrierWidth   float64
	EdgeRamp           []float64
	AttractantCount    int
	AttractantRamp     []float64
	AttractantJitter   float64
	SoftmaxTemperature float64
	ControlPenalty     float64
	MaxRounds          int
	Relaxations        int
	PolarIceFraction   float64
	PeakChillFraction  float64
}

type GenerationResult struct {
	ProvinceCount      int
	LandProvinceCount  int
	InitialIslandCount int
	IslandCount        int
	MergeCount         int
	RoundsAttempted    int
	OceanFraction      float64
}
