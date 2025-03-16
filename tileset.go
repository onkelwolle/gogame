package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/onkelwolle/gogame/constants"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid

	// get the x and y position of the tile in the tileset
	srcX := id % 22
	srcY := id / 22

	// convert the x and y position to pixel coordinates
	srcX *= constants.Tilesize
	srcY *= constants.Tilesize

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize,
		),
	).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynamicTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynamicTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynamicTileset) Img(id int) *ebiten.Image {
	id -= d.gid

	return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		// return a new dynamic tileset

		var dynamicTilesetJSON DynamicTilesetJSON
		err = json.Unmarshal(contents, &dynamicTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynamicTileset := DynamicTileset{}
		dynamicTileset.gid = gid
		dynamicTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynamicTilesetJSON.Tiles {

			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.Trim(tileJSONPath, "../")
			tileJSONPath = strings.Trim(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dynamicTileset.imgs = append(dynamicTileset.imgs, img)
		}

		return &dynamicTileset, nil
	}
	// return a new uniform tileset

	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	tileJSONPath := uniformTilesetJSON.Path
	tileJSONPath = filepath.Clean(tileJSONPath)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.Trim(tileJSONPath, "../")
	tileJSONPath = strings.Trim(tileJSONPath, "../")
	tileJSONPath = filepath.Join("assets/", tileJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		return nil, err
	}

	uniformTileset.img = img
	uniformTileset.gid = gid

	return &uniformTileset, nil
}
