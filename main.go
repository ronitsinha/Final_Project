package main

import (
	"runtime"
	"math"

	sf "github.com/zyedidia/sfml/v2.3/sfml"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var res *Resources

func main() {
	runtime.LockOSThread ()

	var rect *sf.RectangleShape

	res = NewResources()

	InitWeapons ()

	player := NewPlayer (pistol)

	aux := 64.0

	// creates the level 
	// 0- nothing
	// 1- red wall
	// 2- green wall
	// 3- blue wall
	// 4- pink wall
	level1 := CreateLevel ([][]rune{
		{'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'},
		{'1', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '0', '3', '0', '0', '0', '0', '0', '0', '0', '2', '2', '0', '0', '0', '1'},
		{'1', '0', '3', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '2', '2', '1', '0', '0', '0', '0', '0', '0', '2', '2', '0', '0', '0', '1'},
		{'1', '0', '0', '1', '0', '0', '0', '0', '0', '0', '2', '2', '0', '0', '0', '1'},
		{'1', '0', '0', '1', '0', '0', '0', '0', '0', '0', '2', '2', '0', '0', '0', '1'},
		{'1', '0', '0', '0', '0', '0', '0', '4', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '0', '0', '0', '0', '0', '0', '4', '0', '0', '0', '0', '0', '0', '0', '1'},
		{'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'},
	})

	activeLevel := level1

	ActiveCamera := Camera{200, 200, 0, 1.0472}

	window := sf.NewRenderWindow(sf.VideoMode{screenWidth, screenHeight, 32}, "Final Project", sf.StyleDefault, nil)
	window.SetVerticalSyncEnabled(true)
	window.SetFramerateLimit(60)

	renderType := 2

	for window.IsOpen() {
		if event := window.PollEvent(); event != nil {
			switch event.Type {
			case sf.EventClosed:
				window.Close()
				break
			case sf.EventKeyPressed:
				// Renders 2d map when Tab is pressed
				if event.Key.Code == sf.KeyTab {
					renderType = 1
				} else if event.Key.Code == sf.KeyNum1 {
					player.currentWeapon = chainsaw
				} else if event.Key.Code == sf.KeyNum2 {
					player.currentWeapon = pistol
				}

				break
			case sf.EventKeyReleased:
				// Renders raycasted map when tab is not pressed
				if event.Key.Code == sf.KeyTab {
					renderType = 2
				}
				break
			case sf.EventMouseButtonPressed:
				if event.MouseButton.Button == sf.MouseLeft {
					player.currentWeapon.attack = true
				}
				break
			case sf.EventMouseButtonReleased:
				if event.MouseButton.Button == sf.MouseLeft {
					player.currentWeapon.attack = false
				}
				break
			}


		}

		// Movement/Rotation of camera (with collisions)

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
	    	if activeLevel.Tiles [int(ActiveCamera.Y + math.Sin(ActiveCamera.Angle) * 5)/activeLevel.TileHeight][int(ActiveCamera.X)/activeLevel.TileWidth] == '0' {
	    		ActiveCamera.Y += math.Sin(ActiveCamera.Angle) * 5
	    	}

	    	if activeLevel.Tiles [int(ActiveCamera.Y)/activeLevel.TileHeight][int(ActiveCamera.X + math.Cos(ActiveCamera.Angle) * 5)/activeLevel.TileWidth] == '0' {
	    		ActiveCamera.X += math.Cos(ActiveCamera.Angle) * 5
	    	}
	  	}

	  	if sf.KeyboardIsKeyPressed (sf.KeyS) {
	    	if activeLevel.Tiles [int(ActiveCamera.Y - math.Sin(ActiveCamera.Angle) * 5)/activeLevel.TileHeight][int(ActiveCamera.X)/activeLevel.TileWidth] == '0' {
	    		ActiveCamera.Y -= math.Sin(ActiveCamera.Angle) * 5
	    	}

	    	if activeLevel.Tiles [int(ActiveCamera.Y)/activeLevel.TileHeight][int(ActiveCamera.X - math.Cos(ActiveCamera.Angle) * 5)/activeLevel.TileWidth] == '0' {
	    		ActiveCamera.X -= math.Cos(ActiveCamera.Angle) * 5
	    	}
	  	}

	  	if sf.KeyboardIsKeyPressed (sf.KeyP) {
	    	ActiveCamera.Fov += 0.1
	  	}

	  	if sf.KeyboardIsKeyPressed (sf.KeyO) {
	    	ActiveCamera.Fov -= 0.1
	  	}

	  	if sf.KeyboardIsKeyPressed (sf.KeyM) {
	    	aux += 1.0
	  	}

	  	if sf.KeyboardIsKeyPressed (sf.KeyN) {
	    	aux -= 1.0
	  	}

		window.Clear(sf.ColorBlack)

		switch renderType {

		// Renders a 2D top-down map and the rays
		case 1:
			// Render tiles
			for i := 0; i < activeLevel.Width; i ++ {
				for j := 0; j < activeLevel.Height; j ++ {
					if activeLevel.Tiles [j][i] != '0' {
						rect = sf.NewRectangleShape (sf.Vector2f {float32 (activeLevel.TileWidth), float32 (activeLevel.TileHeight)})
						rect.SetPosition (sf.Vector2f {float32 (i * activeLevel.TileWidth), float32 (j * activeLevel.TileHeight)})
						switch activeLevel.Tiles [j][i] {
							case '1':
								rect.SetFillColor (sf.Color {179, 57, 57, 255})
								break
							case '2':
								rect.SetFillColor (sf.Color {45, 136, 45, 255})
								break
							case '3':
								rect.SetFillColor (sf.Color {34, 102, 102, 255})				
								break
							case '4':
								rect.SetFillColor (sf.Color {179, 57, 179, 255})
								break
						}
						window.Draw (rect)
					}
				}
			}

			// Render rays
			for i := -(ActiveCamera.Fov / 2); i < ActiveCamera.Fov/2; i += 0.01 {
				targetX, targetY := activeLevel.CastRay(ActiveCamera.X, ActiveCamera.Y, ActiveCamera.Angle+i)
				var c sf.Color

				switch activeLevel.Tiles[targetY/activeLevel.TileHeight][targetX/activeLevel.TileWidth] {
				case '1':
					c = sf.Color {179, 57, 57, 255}
					break
				case '2':
					c = sf.Color {45, 136, 45, 255}
					break
				case '3':
					c = sf.Color {34, 102, 102, 255}
					break
				case '4':
					c = sf.Color {179, 57, 179, 255}
					break
				}

				if activeLevel.GetTileSide (targetX, targetY) == "NS" {
					r := c.R - 20
					g := c.G - 20
					b := c.B - 20

					c = sf.Color {r, g, b, 255}
				}
				vertA := sf.Vertex {sf.Vector2f {float32 (ActiveCamera.X), float32(ActiveCamera.Y)}, c, sf.Vector2f {float32 (ActiveCamera.X), float32(ActiveCamera.Y)}}
				vertB := sf.Vertex {sf.Vector2f {float32 (targetX), float32(targetY)}, c, sf.Vector2f {float32 (targetX), float32(targetY)}}
			
				vertArray := sf.NewVertexArray ([]sf.Vertex {vertA, vertB}, sf.PrimitiveLines)
				window.DrawPrimitives (vertArray)
			}

			// Render camera as a white rectangle
			rect = sf.NewRectangleShape (sf.Vector2f {20, 20})
			rect.SetPosition (sf.Vector2f {float32 (ActiveCamera.X - 10), float32 (ActiveCamera.Y - 10)})
			rect.SetFillColor (sf.ColorWhite)
			window.Draw (rect)
			break
		case 2:
			// Renders the raycasted map
			// TODO: render textures (see assets/images/textures)


			rect = sf.NewRectangleShape (sf.Vector2f {float32 (screenWidth), float32 (screenHeight / 2)})
			rect.SetPosition (sf.Vector2f {0, float32 (screenHeight/2)})
			rect.SetFillColor (sf.Color {179, 108, 57, 255})
			
			j := 0
			for i := -(ActiveCamera.Fov/2); i < ActiveCamera.Fov/2; i += (ActiveCamera.Fov / float64 (screenWidth)) * math.Cos (i) {
				targetX, targetY := activeLevel.CastRay(ActiveCamera.X, ActiveCamera.Y, ActiveCamera.Angle+i)
			
				var clr sf.Color

				switch activeLevel.Tiles [targetY/activeLevel.TileHeight][targetX/activeLevel.TileWidth] {
				case '1':
					clr = sf.Color {179, 57, 57, 255}
					break
				case '2':
					clr = sf.Color {45, 136, 45, 255}
					break
				case '3':
					clr = sf.Color {34, 102, 102, 255}
					break
				case '4':
					clr = sf.Color {179, 57, 179, 255}
					break
				}

				// Darken color if on X-axis
				if activeLevel.GetTileSide (targetX, targetY) == "NS" {
					r := clr.R - 20
					g := clr.G - 20
					b := clr.B - 20

					clr = sf.Color {r, g, b, 255}
				}

				distance := math.Sqrt (math.Pow (ActiveCamera.X-float64 (targetX), 2) + math.Pow (ActiveCamera.Y - float64 (targetY), 2))
				z := distance * math.Cos (i)
				lineHeight := float64 (screenHeight) / z * 64

				vertA := sf.Vertex {sf.Vector2f {float32 (j), float32(screenHeight/2 - int (lineHeight))}, clr, sf.Vector2f {float32 (j), float32(screenHeight/2 - int (lineHeight))}}
				vertB := sf.Vertex {sf.Vector2f {float32 (j), float32(screenHeight/2 + int (lineHeight))}, clr, sf.Vector2f {float32 (j), float32(screenHeight/2 + int (lineHeight))}}
				vertArray := sf.NewVertexArray ([]sf.Vertex {vertA, vertB}, sf.PrimitiveLines)
				window.DrawPrimitives (vertArray)
				j ++
			}

			player.currentWeapon.Update ()

			window.Draw (player.currentWeapon)

			break
		}

		window.Display()
	}
}