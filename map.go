package main

import (
  "fmt"
  "math"

  sf "github.com/zyedidia/sfml/v2.3/sfml"
)


type Camera struct {
  X, Y, Angle, Fov float64
}

type Level struct {
  Height, Width, TileWidth, TileHeight int
  Tiles                                [][]rune
}

func (targetLevel *Level) CastRay(x, y, Angle float64) (int, int) {
  for targetLevel.Tiles[int(y)/targetLevel.TileHeight][int(x)/targetLevel.TileWidth] == '0' {
    x += math.Cos(Angle)
    y += math.Sin(Angle)
  }
  return int(x), int(y)
}

func (targetLevel *Level) GetTileSide(x, y int) (side string) {
  dx := float64(targetLevel.TileWidth/2-x%targetLevel.TileWidth) / float64(targetLevel.TileWidth)
  dy := float64(targetLevel.TileHeight/2-y%targetLevel.TileHeight) / float64(targetLevel.TileHeight)
  if math.Abs(float64(dx)) >= math.Abs(float64(dy)) {
    return "LW"
  } else {
    return "NS"
  }
}

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

func Move (ActiveCamera Camera, aux float64, renderType int) {
  if sf.KeyboardIsKeyPressed (sf.KeyA) {
    ActiveCamera.Angle -= 0.1
    if ActiveCamera.Angle < 0 {
      ActiveCamera.Angle = 2 * math.Pi
    }
  }

  if sf.KeyboardIsKeyPressed (sf.KeyD) {
    ActiveCamera.Angle += 0.1
    if ActiveCamera.Angle > 2*math.Pi {
      ActiveCamera.Angle = 0
    }
  }

  if sf.KeyboardIsKeyPressed (sf.KeyW) {
    ActiveCamera.X += math.Cos(ActiveCamera.Angle) * 5
    ActiveCamera.Y += math.Sin(ActiveCamera.Angle) * 5
  }

  if sf.KeyboardIsKeyPressed (sf.KeyS) {
    ActiveCamera.X -= math.Cos(ActiveCamera.Angle) * 5
    ActiveCamera.Y -= math.Sin(ActiveCamera.Angle) * 5
  }

  if sf.KeyboardIsKeyPressed (sf.KeyP) {
    ActiveCamera.Fov += 0.1
  }

  if sf.KeyboardIsKeyPressed (sf.KeyO) {
    ActiveCamera.Fov -= 0.1
  }

  if sf.KeyboardIsKeyPressed (sf.KeyR) {
    if renderType == 1 {
      renderType = 2
    } else {
      renderType = 1
    }
  }

  if sf.KeyboardIsKeyPressed (sf.KeyM) {
    aux += 1.0
  }

  if sf.KeyboardIsKeyPressed (sf.KeyN) {
    aux -= 1.0
  }
}