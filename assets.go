package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	sf "github.com/zyedidia/sfml/v2.3/sfml"
)

type Resources struct {
	images map[string]*sf.Texture
}

func NewResources() *Resources {
	r := new(Resources)

	r.images = make(map[string]*sf.Texture)

	r.LoadAllImages("./assets/images")

	return r
}

func (r *Resources) LoadAllImages(dir string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			r.LoadAllImages(dir + "/" + f.Name())
		} else if filepath.Ext(f.Name()) == ".png" {
			texture := sf.NewTexture(dir + "/" + f.Name())
			r.images[f.Name()] = texture
		}
	}
}
