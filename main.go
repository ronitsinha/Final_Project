package main

import (
	"runtime"
	"time"
	//"fmt"
	"math"
	"strconv"

	sf "github.com/zyedidia/sfml/v2.3/sfml"
	"github.com/kellpossible/math32"
)

const (
	screenWidth  = 512
	screenHeight = 384

	texWidth = 64
	texHeight = 64

	mapWidth = 24
	mapHeight = 24

)

var player *Player
var res *Resources


// How many targets left for the player to shoot
var totalTargets int
// How much time the player has
var timeLimit time.Duration = time.Minute - time.Second*10

// World map to be parsed to render textures
var worldMap [mapWidth][mapHeight]int = 
[mapWidth][mapHeight]int {
	{8,9,8,8,8,8,8,8,8,8,8,4,4,6,4,4,6,4,6,4,4,4,6,4},
	{8,0,0,0,0,0,0,0,0,0,8,9,0,0,0,0,0,0,0,0,0,0,0,4},
	{8,0,3,3,0,0,0,0,0,8,8,4,0,0,0,0,0,0,0,0,0,0,0,6},
	{8,0,0,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,6},
	{8,0,3,3,0,0,0,0,0,8,8,4,0,0,0,0,0,0,0,0,0,0,0,4},
	{8,0,0,0,0,0,0,0,0,0,8,4,0,0,0,0,0,6,6,6,0,6,4,6},
	{8,8,8,8,0,8,8,8,8,8,8,4,4,4,4,4,4,6,0,0,0,0,0,6},
	{7,7,9,7,0,7,7,7,7,0,8,0,8,0,8,0,8,4,0,4,0,6,0,6},
	{7,7,0,0,0,0,0,0,7,8,0,8,0,8,0,8,8,6,0,0,0,0,0,6},
	{7,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,8,6,0,0,0,0,0,4},
	{7,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,8,6,0,6,0,6,0,6},
	{7,7,0,0,0,0,0,0,7,8,0,8,0,8,0,8,8,5,4,6,0,6,6,6},
	{7,7,7,7,0,7,7,7,7,8,8,4,0,6,8,4,8,3,3,3,0,3,3,3},
	{2,2,2,2,0,2,2,2,2,4,6,4,0,0,6,0,6,3,0,0,0,0,0,3},
	{2,2,0,0,0,0,0,9,2,4,0,0,0,0,0,0,4,3,0,0,0,0,0,3},
	{2,0,0,0,0,0,0,0,2,4,0,0,0,0,0,0,4,3,0,0,0,0,0,3},
	{1,0,0,0,0,0,0,0,1,4,4,4,4,4,6,0,6,3,3,0,0,0,3,3},
	{2,0,0,0,0,0,0,0,2,2,2,1,2,2,2,6,6,0,0,5,0,5,0,5},
	{2,2,0,0,0,0,0,2,2,2,0,0,0,2,2,0,5,0,5,0,0,0,5,5},
	{2,0,0,0,0,0,0,0,2,0,0,0,0,0,2,5,0,5,0,5,0,5,0,5},
	{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,9},
	{2,0,0,0,0,0,0,0,2,0,0,0,0,0,2,5,0,5,0,5,0,5,0,5},
	{2,2,0,0,0,0,0,2,2,2,0,0,0,2,2,0,5,0,5,0,0,0,5,5},
	{2,2,2,2,1,2,2,2,2,2,2,1,2,2,2,5,5,5,5,5,5,5,5,5},
}

var tileWidth int
var tileHeight int

// Collection of textures to draw
var texture [10]*sf.Image

var posX float32 = 2
var posY float32 = 2

var dirX float32 = -1
var dirY float32 = 0

var gameOver bool

func main () {
	runtime.LockOSThread()
	res = NewResources ()

	for x := 0; x < len (worldMap); x ++ {
		for y := 0; y < len (worldMap[0]); y ++ {
			if worldMap [x][y] == 9 {
				totalTargets ++
			}
		}
	}

	startTime := time.Now ()

	// Load fonts

	fntAbove := sf.NewFont ("./assets/fonts/Digitalt.ttf")
	
	scoreTxt := sf.NewText ("Targets Remaining: " + strconv.Itoa (totalTargets), fntAbove, 16)
	scoreTxt.SetColor (sf.ColorGreen)

	timeTxt := sf.NewText ("Time Remaining: " + strconv.Itoa (int ((timeLimit - (time.Since (startTime)))/time.Second)), fntAbove, uint (16))
	timeTxt.SetOrigin (sf.Vector2f {timeTxt.GetGlobalBounds ().Width, 0})
	timeTxt.SetPosition (sf.Vector2f{screenWidth, 0})

	gameOverTxt := sf.NewText ("Game Over!", fntAbove, uint (32))
	gameOverTxt.SetOrigin (sf.Vector2f {timeTxt.GetGlobalBounds ().Width/2, timeTxt.GetGlobalBounds ().Height/2})
	gameOverTxt.SetPosition (sf.Vector2f{screenWidth/2, screenHeight/2})
	gameOverTxt.SetColor (sf.ColorRed)

	var planeX float32 = 0
	var planeY float32 = 0.66

	tileWidth = screenWidth / mapWidth
	tileHeight = screenHeight / mapHeight

	window := sf.NewRenderWindow(sf.VideoMode{screenWidth, screenHeight, 32}, "Final Project", sf.StyleDefault, nil)
	window.SetVerticalSyncEnabled(true)
	window.SetFramerateLimit(60)

	var dt float32

	renderType := 2

	// Load Textures

	texture = [10]*sf.Image {
		sf.NewImage ("./assets/images/walls/eagle.png"),
		sf.NewImage ("./assets/images/walls/redbrick.png"),
		sf.NewImage ("./assets/images/walls/purplestone.png"),
		sf.NewImage ("./assets/images/walls/greystone.png"),
		sf.NewImage ("./assets/images/walls/bluestone.png"),
		sf.NewImage ("./assets/images/walls/mossy.png"),
		sf.NewImage ("./assets/images/walls/wood.png"),
		sf.NewImage ("./assets/images/walls/colorstone.png"),
		sf.NewImage ("./assets/images/walls/targetBrick.png"),
		sf.NewImage ("./assets/images/walls/targetBrickHit.png"),
	}

	// Create global weapon vars (see player.go) and instantiate the player with a pistol

	InitWeapons ()
	player = NewPlayer (pistol)

	for window.IsOpen() {
		start := time.Now ()

		if event := window.PollEvent(); event != nil {
			switch event.Type {
			case sf.EventClosed:
				window.Close()
			case sf.EventKeyPressed:
				// Renders 2d map when Tab is pressed
				if event.Key.Code == sf.KeyTab {
					renderType = 1
				} else if event.Key.Code == sf.KeyNum1 {
					player.currentWeapon = chainsaw
				} else if event.Key.Code == sf.KeyNum2 {
					player.currentWeapon = pistol
				} else if event.Key.Code == sf.KeyNum3 {
					player.currentWeapon = shotgun
				}

				break
			case sf.EventKeyReleased:
				// Renders raycasted map when tab is not pressed
				if event.Key.Code == sf.KeyTab {
					renderType = 2
				}
				break
			case sf.EventMouseButtonPressed:
				// Fire/discharge weapon

				if event.MouseButton.Button == sf.MouseLeft {
					player.currentWeapon.attack = true
				}
				break
			case sf.EventMouseButtonReleased:
				// Stop firing the weapon (for chainsaw only)
				if event.MouseButton.Button == sf.MouseLeft {
					player.currentWeapon.attack = false
				}
				break
			}
		}

		window.Clear(sf.ColorBlack)

		DrawRaycast (window, posX, posY, renderType, dirX, dirY, planeX, planeY)		
		
		if renderType == 2 {
			window.Draw (player.currentWeapon)
		}

		timeTxt.SetOrigin (sf.Vector2f {timeTxt.GetGlobalBounds ().Width, 0})
		timeTxt.SetPosition (sf.Vector2f{screenWidth, 0})

		// Only draw game over text if game is actually over
		if gameOver {
			gameOverTxt.SetOrigin (sf.Vector2f {gameOverTxt.GetGlobalBounds ().Width/2, gameOverTxt.GetGlobalBounds ().Height/2})
			gameOverTxt.SetPosition (sf.Vector2f{screenWidth/2, screenHeight/2})
			window.Draw (gameOverTxt)
		} else {
			timeTxt.SetString ("Time Remaining: " + strconv.Itoa (int ((timeLimit - (time.Since (startTime)))/time.Second)))
			scoreTxt.SetString ("Targets Remaining: " + strconv.Itoa (totalTargets))
		}

		window.Draw (scoreTxt)
		window.Draw (timeTxt)

		var moveSpeed float32 = dt * 5
		var rotSpeed float32 = dt * 3

		// Movement w/ collision detection against worldMap grid
		if !gameOver {
			if sf.KeyboardIsKeyPressed (sf.KeyW) {
				if worldMap [int (posX + dirX * moveSpeed)][int (posY)] == 0 {
					posX += dirX * moveSpeed
				}

				if worldMap [int (posX)][int (posY + dirY * moveSpeed)] == 0 {
					posY += dirY * moveSpeed
				}
			}

			if sf.KeyboardIsKeyPressed (sf.KeyS) {
				if worldMap [int (posX - dirX * moveSpeed)][int (posY)] == 0 {
					posX -= dirX * moveSpeed
				}

				if worldMap [int (posX)][int (posY - dirY * moveSpeed)] == 0 {
					posY -= dirY * moveSpeed
				}
			}

			// Camera Rotation
			if sf.KeyboardIsKeyPressed (sf.KeyD) {
				oldDirX := dirX
				dirX = dirX * math32.Cos (-rotSpeed) - dirY * math32.Sin (-rotSpeed)
				dirY = oldDirX * math32.Sin (-rotSpeed) + dirY * math32.Cos (-rotSpeed)

				oldPlaneX := planeX
				planeX = planeX * math32.Cos (-rotSpeed) - planeY * math32.Sin (-rotSpeed)
				planeY = oldPlaneX * math32.Sin (-rotSpeed) + planeY * math32.Cos (-rotSpeed)
			}

			if sf.KeyboardIsKeyPressed (sf.KeyA) {
				oldDirX := dirX
				dirX = dirX * math32.Cos (rotSpeed) - dirY * math32.Sin (rotSpeed)
				dirY = oldDirX * math32.Sin (rotSpeed) + dirY * math32.Cos (rotSpeed)

				oldPlaneX := planeX
				planeX = planeX * math32.Cos (rotSpeed) - planeY * math32.Sin (rotSpeed)
				planeY = oldPlaneX * math32.Sin (rotSpeed) + planeY * math32.Cos (rotSpeed)
			}
		}

		window.Display()
		dt = float32 (time.Since (start)) / float32 (time.Second)
	
		// If the player has shot all the targets OR time has run out, then the game is over
		if totalTargets == 0 || time.Since (startTime) >= timeLimit {
			gameOver = true
		}
	}
}

func DrawRaycast (window *sf.RenderWindow, posX, posY float32, renderType int, dirX, dirY float32, planeX, planeY float32) {
	
		if renderType == 1 {
			// Render top-down view
			for x := 0; x < mapWidth; x ++ {
				for y := 0; y < mapHeight; y ++ {
					rect := sf.NewRectangleShape (sf.Vector2f {float32 (tileWidth), float32 (tileHeight)})
					rect.SetPosition (sf.Vector2f {float32 (x * tileWidth), float32(y * tileHeight)})

					switch worldMap[x][y] {
						case 0:
							rect.SetFillColor (sf.ColorBlack)
							break
						default:
							rect.SetFillColor (sf.ColorRed)
							break
					}
					window.Draw (rect)
				}
			}

			rect := sf.NewRectangleShape (sf.Vector2f {float32 (tileWidth/2), float32 (tileHeight/2)})
			rect.SetOrigin (sf.Vector2f {rect.GetGlobalBounds ().Width/2, rect.GetGlobalBounds ().Height/2})
			rect.SetPosition (sf.Vector2f {posX * float32 (tileWidth), posY * float32(tileHeight)})

			window.Draw (rect)

			vertA := sf.Vertex {sf.Vector2f {posX * float32 (tileWidth), posY * float32(tileHeight)}, sf.ColorWhite, sf.Vector2f {posX * float32 (tileWidth), posY * float32(tileHeight)}}
			vertB := sf.Vertex {sf.Vector2f {(posX + dirX) * float32 (tileWidth), (posY + dirY) * float32(tileHeight)}, sf.ColorWhite, sf.Vector2f {posX * float32 (tileWidth), posY * float32(tileHeight)}}
		
			line := sf.NewVertexArray ([]sf.Vertex {vertA, vertB}, sf.PrimitiveLines)
			window.DrawPrimitives (line)
		} else if renderType == 2 {
			for x := 0; x < screenWidth; x ++ {
				// Camera position used for distance of walls from player
				var cameraX float32 = 2 * float32 (x) / float32 (screenWidth) - 1
				var rayPosX float32 = posX
				var rayPosY float32 = posY
				
				var rayDirX float32 = dirX + planeX * cameraX
				var rayDirY float32 = dirY + planeY * cameraX

				var mapX int = int (rayPosX)
				var mapY int = int (rayPosY)

				var sideDistX float32
				var sideDistY float32

				// The x dist of the ray
				deltaDistX := math32.Sqrt (1 + (rayDirY * rayDirY) / (rayDirX * rayDirX))
				
				// the y distance of the ray
				deltaDistY := math32.Sqrt (1 + (rayDirX * rayDirX) / (rayDirY * rayDirY))
				
				// the perpendicular distance to a wall
				var perpWallDist float32

				var stepX int
				var stepY int

				// Raycasting
				hit := 0
				var side int

				if rayDirX < 0 {
					stepX = -1
					sideDistX = (rayPosX - float32 (mapX)) * deltaDistX
				} else {
					stepX = 1
					sideDistX = (float32(mapX) + 1.0 - rayPosX) * deltaDistX
				}
				if rayDirY < 0 {
					stepY = -1
					sideDistY = (rayPosY - float32(mapY)) * deltaDistY
				} else {
					stepY = 1
					sideDistY = (float32(mapY) + 1.0 - rayPosY) * deltaDistY
				}

				for hit == 0 {
					if sideDistX < sideDistY {
						sideDistX += deltaDistX
						mapX += stepX
						side = 0
					} else {
						sideDistY += deltaDistY
						mapY += stepY
						side = 1
					}

					if worldMap[mapX][mapY] > 0 {
						hit = 1
					}
				}

				if side == 0 {
					perpWallDist = (float32 (mapX) - rayPosX + (1 - float32 (stepX)) / 2) / rayDirX
				} else {
					perpWallDist = (float32 (mapY) - rayPosY + (1 - float32 (stepY)) / 2) / rayDirY
				}

				// Rendering

				lineHeight := int (screenHeight / perpWallDist)

				drawStart := -lineHeight / 2 + screenHeight / 2
				
				if drawStart < 0 {
					drawStart = 0
				}

				drawEnd := lineHeight / 2 + screenHeight / 2

				if drawEnd >= screenHeight {
					drawEnd = screenHeight - 1
				}

				texNum := worldMap [mapX][mapY] - 1

				var wallX float32

				if side == 0 {
					wallX = rayPosY + perpWallDist * rayDirY
				} else {
					wallX = rayPosX + perpWallDist * rayDirX
				}

				// What part of a wall texture should be drawn?

				wallX -= float32 (math.Floor (float64 (wallX)))

				texX := int (wallX * float32 (texWidth))

				if side == 0 && rayDirX > 0 {
					texX = texWidth - texX - 1
				}

				if side == 1 && rayDirY < 0 {
					texX = texWidth - texX - 1
				}

				vertices := make ([]sf.Vertex, 0)

				for y := drawStart; y < drawEnd; y ++ {
					d := y * 256 - screenHeight * 128 + lineHeight * 128
					texY := ((d * texHeight) / lineHeight) / 256
					// Finds the correct pixels
					color := texture[texNum].GetPixel (uint(texX), uint (texY))
				
					if side == 1 {
						r := color.R / 2
						g := color.G / 2
						b := color.B / 2

						color = sf.Color {r, g, b, 255}
					}

					// Add the pixel as a vertex in a lines
					vertices = append (vertices, sf.Vertex {sf.Vector2f {float32(x), float32 (y)}, color, sf.Vector2f {float32(x), float32 (y)}})
				}


				// Draw the line
				verLine := sf.NewVertexArray (vertices, sf.PrimitiveLines)
				window.DrawPrimitives (verLine)

				// Floor/Ceiling casting
				floorVertices := make ([]sf.Vertex, 0)
				ceilVertices := make ([]sf.Vertex, 0)

				var floorXWall float32
				var floorYWall float32

				if side == 0 && rayDirX > 0 {
					floorXWall = float32 (mapX)
					floorYWall = float32 (mapY) + wallX
				} else if side == 0 && rayDirX < 0 {
					floorXWall = float32 (mapX) + 1
					floorYWall = float32 (mapY) + wallX
				} else if side == 1 && rayDirY > 0 {
					floorXWall = float32 (mapX) + wallX
					floorYWall = float32 (mapY)
				} else {
					floorXWall = float32 (mapX) + wallX
					floorYWall = float32 (mapY) + 1
				}

				// Similar rendering method from before

				var distWall float32
				var distPlayer float32
				var currentDist float32

				distWall = perpWallDist
				distPlayer = 0

				if drawEnd < 0 {drawEnd = screenHeight}

				for y := drawEnd + 1; y < screenHeight; y ++ {
					currentDist = screenHeight / (2 * float32(y) - float32 (screenHeight))

					weight := (currentDist - distPlayer) / (distWall - distPlayer)

					currentFloorX := weight * floorXWall + (1 - weight) * posX
					currentFloorY := weight * floorYWall + (1 - weight) * posY

					var floorTexX int
					var floorTexY int


					// Getting tex coords that can work for both floor and ceiling
					floorTexX = int (currentFloorX * texWidth) % texWidth
					floorTexY = int (currentFloorY * texHeight) % texHeight

					color := texture[3].GetPixel (uint(floorTexX), uint (floorTexY))
					ceilColor := texture[6].GetPixel (uint(floorTexX), uint (floorTexY))

					// Making floor and ceiling lines with pixel vertices
					floorVertices = append (floorVertices, sf.Vertex {sf.Vector2f {float32(x), float32 (y)}, color, sf.Vector2f {float32(x), float32 (y)}})
					ceilVertices = append (ceilVertices, sf.Vertex {sf.Vector2f {float32(x), float32 (screenHeight - y)}, ceilColor, sf.Vector2f {float32(x), float32 (screenHeight - y)}})	
				}

				// Drawing floor/ceiling lines.
				floorVerLine := sf.NewVertexArray (floorVertices, sf.PrimitiveLines)
				ceilVerLine := sf.NewVertexArray (ceilVertices, sf.PrimitiveLines)
				window.DrawPrimitives (floorVerLine)
				window.DrawPrimitives (ceilVerLine)

				// THis is the solid color rendering mode, commented out just as a backup

				// var color sf.Color

				// switch worldMap[mapX][mapY] {
				// case 1:
				// 	color = sf.ColorRed
				// case 2:
				// 	color = sf.ColorGreen
				// case 3:
				// 	color = sf.ColorBlue
				// case 4:
				// 	color = sf.ColorWhite
				// case 5:
				// 	color = sf.ColorYellow
				// }

				// if side == 1 {
				// 	r := color.R / 2
				// 	g := color.G / 2
				// 	b := color.B / 2

				// 	color = sf.Color {r, g, b, 255}
				// }

				// vertA := sf.Vertex {sf.Vector2f{float32 (x), float32 (drawStart)}, color, sf.Vector2f{float32 (x), float32 (drawStart)}}
				// vertB := sf.Vertex {sf.Vector2f{float32 (x), float32 (drawEnd)}, color, sf.Vector2f{float32 (x), float32 (drawEnd)}}


				// verLine := sf.NewVertexArray ([]sf.Vertex {vertA, vertB}, sf.PrimitiveLines)
				// window.DrawPrimitives (verLine)



				// This is for printing out memory usage (for debugging)

				// var memstats runtime.MemStats
				// runtime.ReadMemStats(&memstats)
				// fmt.Println(memstats.Alloc, memstats.Sys)
			}
		}
}