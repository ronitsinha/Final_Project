package main

import (
  "fmt"
  "math"
)


type Camera struct {
  X, Y, Angle, Fov float64
}

type Level struct {
  Height, Width, TileWidth, TileHeight int
  Tiles                                [][]rune
}

// Casts as ray from position x, y in direction Angle and returns the coordinates of where the ray hits a wall/non-walkable
func (targetLevel *Level) CastRay(x, y, Angle float64) (int, int) {
  for targetLevel.Tiles[int(y)/targetLevel.TileHeight][int(x)/targetLevel.TileWidth] == '0' {
    x += math.Cos(Angle)
    y += math.Sin(Angle)
  }
  return int(x), int(y)
}

// Determines where the side hide is North/South or Lest/West
func (targetLevel *Level) GetTileSide(x, y int) (side string) {
  dx := float64(targetLevel.TileWidth/2-x%targetLevel.TileWidth) / float64(targetLevel.TileWidth)
  dy := float64(targetLevel.TileHeight/2-y%targetLevel.TileHeight) / float64(targetLevel.TileHeight)
  if math.Abs(float64(dx)) >= math.Abs(float64(dy)) {
    return "LW"
  } else {
    return "NS"
  }
}

// Creates a level from an array of rune arrays
func CreateLevel(tiles [][]rune) (newLevel *Level) {
  newLevel = new(Level)
  newLevel.Tiles = tiles
  newLevel.Height = len(tiles)
  newLevel.Width = len(tiles[0])
  newLevel.TileWidth = screenWidth / newLevel.Width
  newLevel.TileHeight = screenHeight / newLevel.Height
  fmt.Println("Creating level...", "height:", newLevel.Height, "width:", newLevel.Width, "tilewidth:", newLevel.TileWidth, "tileheight:", newLevel.TileHeight)
  return newLevel
}