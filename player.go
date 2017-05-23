package main

import (
	"time"

	sf "github.com/zyedidia/sfml/v2.3/sfml"
 	"strconv"
)

var (
	// Weapons that the player may use
	chainsaw *Weapon
	pistol *Weapon
	shotgun *Weapon
)

type Weapon struct {
	*sf.Sprite

	// Standby animations (see chainsaw)
	standbyTextures []*sf.Texture

	// Attacking anims
	activeTextures []*sf.Texture
	
	// At what frame does the weapon fire?
	shootFrame int
	// What range does the weapon have (decent for pistol, medium for shotgun, low for chainsaw)
	rng float32
	// Used to check if the animation is at the shoot frame
	activeIndices []int

	// is the weapon attacking
	attack bool

	animateSpeed int

	// Should the weapon continously animate when the mouse is held down (see chainsaw)
	continuousActive bool
}

type Player struct {
	currentWeapon *Weapon
}

func InitWeapons() {
	chainsaw = NewWeapon ([]int {0, 1}, []int {2, 3}, "chainsaw", 10, true, 2, 3)
	pistol = NewWeapon ([]int {13}, []int {12, 11, 10, 8, 10}, "pistol", 100, false, 10, 8)
	shotgun = NewWeapon ([]int {33}, []int {32, 28, 31, 28}, "shotgun", 100, false, 5, 32)
}

// Currently, the player's only purpose is to keep track of the currently equipped weapon, but later it may have health or other attributes
func NewPlayer (wpn *Weapon) *Player {
	ply := new (Player)
	ply.currentWeapon = wpn

	return ply
}

// Initialize the weapon with appropriate information

func NewWeapon (standbyIndices, activeIndices []int, name string, speed int, contA bool, rng float32, shootFrame int) *Weapon {
	wpn := new (Weapon)

	wpn.standbyTextures = make ([]*sf.Texture, 0)
	wpn.activeTextures = make ([]*sf.Texture, 0)

	wpn.rng = rng
	wpn.shootFrame = shootFrame
	wpn.activeIndices = activeIndices

	wpn.animateSpeed = speed
	wpn.Sprite = sf.NewSprite (res.images[name + "_" + strconv.Itoa (standbyIndices[0]) + ".png"])
	wpn.continuousActive = contA

	wpn.SetOrigin (sf.Vector2f {wpn.GetGlobalBounds ().Width/2, wpn.GetGlobalBounds ().Height})
	wpn.SetPosition (sf.Vector2f {screenWidth/2, screenHeight})

	for i := 0; i < len (standbyIndices); i ++ {
		wpn.standbyTextures = append (wpn.standbyTextures, res.images [name + "_" + strconv.Itoa (standbyIndices[i]) + ".png"])
	}

	for i := 0; i < len (activeIndices); i ++ {
		wpn.activeTextures = append (wpn.activeTextures, res.images [name + "_" + strconv.Itoa (activeIndices[i]) + ".png"])
	}


	// Animation goroutine

	go func () {
		for {
			wpn.Animate ()
		}
	} ()

	return wpn
}

// Checks if a weapon shot something, and returns the int value of the wall in the worldMap, as well as the location (x and y) of the tile in the map
// Only checks within a weapon's range
func CheckShoot (rng float32) (int, int, int) {
	pX := posX
	pY := posY

	for i := float32 (1); i < rng; i ++ {
		if (int (pX + dirX*i) >= mapWidth || int (pX + dirY*i) >= mapHeight) || (int (pX + dirX*i) < 0 || int (pX + dirY*i) < 0)  {
			return -1, -1, -1
		}

		if worldMap [int (pX + dirX*i)][int (pY + dirY*i)] != 0 {
			return int (pX + dirX*i), int (pY + dirY*i), worldMap [int (pX + dirX*i)][int (pY + dirY*i)]
		}
	}

	return -1, -1, -1
}

func (w *Weapon) Animate () {
	if w.attack {
		for i := 0; i < len (w.activeTextures); i ++ {
			
			// the attack frame (i.e when the pistol reaches its recoil frame)
			if w.activeIndices [i] == w.shootFrame {
				x, y, val := CheckShoot (w.rng)

				// A target has been shot, one down!
				if val == 9 {
					worldMap[x][y] = 10
					totalTargets --
				}
			}

			// Animating attack
			w.SetTexture (w.activeTextures[i], true)
			w.SetOrigin (sf.Vector2f {w.GetGlobalBounds ().Width/2, w.GetGlobalBounds ().Height})
			w.SetPosition (sf.Vector2f {screenWidth/2, screenHeight})

			time.Sleep (time.Duration (w.animateSpeed) * time.Millisecond)
		}

		// Only stop attack animating if it is a one shot type weapon (all but chainsaw)
		if !w.continuousActive {
			time.Sleep (time.Duration (w.animateSpeed) * time.Millisecond)

			w.attack = false
		}

	} else {
		// Animate standby
		for i := 0; i < len (w.standbyTextures); i ++ {
			w.SetTexture (w.standbyTextures[i], true)
			w.SetOrigin (sf.Vector2f {w.GetGlobalBounds ().Width/2, w.GetGlobalBounds ().Height})
			w.SetPosition (sf.Vector2f {screenWidth/2, screenHeight})

			time.Sleep (30 * time.Millisecond)
		}
	}
}