package main

import (
	"time"

	sf "github.com/zyedidia/sfml/v2.3/sfml"
 	"strconv"
)

var (
	// There will be one weapon for each set of weapon sprites see assets/images/weapons
	chainsaw *Weapon
	pistol *Weapon
)

type Weapon struct {
	*sf.Sprite

	standbyTextures []*sf.Texture
	activeTextures []*sf.Texture

	attack bool
}

type Player struct {
	currentWeapon *Weapon
}

func InitWeapons() {
	chainsaw = NewWeapon ([]int {0, 1}, []int {2, 3}, "chainsaw")
	pistol = NewWeapon ([]int {13}, []int {12, 11, 10}, "pistol")
}

// Currently, the player's only purpose is to keep track of the currently equipped weapon, but later it may have health or other attributes
func NewPlayer (wpn *Weapon) *Player {
	ply := new (Player)
	ply.currentWeapon = wpn

	return ply
}

func NewWeapon (standbyIndices, activeIndices []int, name string) *Weapon {
	wpn := new (Weapon)

	wpn.standbyTextures = make ([]*sf.Texture, 0)
	wpn.activeTextures = make ([]*sf.Texture, 0)

	wpn.Sprite = sf.NewSprite (res.images[name + "_" + strconv.Itoa (standbyIndices[0]) + ".png"])

	wpn.SetOrigin (sf.Vector2f {wpn.GetGlobalBounds ().Width/2, wpn.GetGlobalBounds ().Height})
	wpn.SetPosition (sf.Vector2f {screenWidth/2, screenHeight})

	for i := 0; i < len (standbyIndices); i ++ {
		wpn.standbyTextures = append (wpn.standbyTextures, res.images [name + "_" + strconv.Itoa (standbyIndices[i]) + ".png"])
	}

	for i := 0; i < len (activeIndices); i ++ {
		wpn.activeTextures = append (wpn.activeTextures, res.images [name + "_" + strconv.Itoa (activeIndices[i]) + ".png"])
	}

	return wpn
}

// TODO: Fix weird animation stuff
func (w *Weapon) Update () {
	go w.Animate ()
}

func (w *Weapon) Animate () {
	if w.attack {
		for i := 0; i < len (w.activeTextures); i ++ {
			w.SetTexture (w.activeTextures[i], true)
			time.Sleep (30 * time.Millisecond)

			w.SetOrigin (sf.Vector2f {w.GetGlobalBounds ().Width/2, w.GetGlobalBounds ().Height})
			w.SetPosition (sf.Vector2f {screenWidth/2, screenHeight})
		}	
	} else {
		for i := 0; i < len (w.standbyTextures); i ++ {
			w.SetTexture (w.standbyTextures[i], true)
			time.Sleep (30 * time.Millisecond)

			w.SetOrigin (sf.Vector2f {w.GetGlobalBounds ().Width/2, w.GetGlobalBounds ().Height})
			w.SetPosition (sf.Vector2f {screenWidth/2, screenHeight})
		}
	}

	w.attack = false
}