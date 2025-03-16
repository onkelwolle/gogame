package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/onkelwolle/gogame/animations"
	"github.com/onkelwolle/gogame/constants"
	"github.com/onkelwolle/gogame/entities"
	"github.com/onkelwolle/gogame/spritesheet"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *TilemapJSON
	tilesets          []Tileset
	tilemapImg        *ebiten.Image
	cam               *Camera
	colliders         []image.Rectangle
}

func NewGame() *Game {
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   100,
				Y:   100,
			},
			Health: 100,
			Animations: map[entities.PlayerState]*animations.Animation{
				entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
				entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
				entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
				entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
			},
		},
		playerSpriteSheet: playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   50,
					Y:   50,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150,
					Y:   150,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   75,
					Y:   75,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   200,
					Y:   200,
				},
				AmtHeal: 10,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		tilesets:    tilesets,
		cam:         NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	return &game
}

func (g *Game) Update() error {

	//react to key presses

	g.player.Dx = 0
	g.player.Dy = 0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Dx = 2
	}

	g.player.X += g.player.Dx

	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy

	CheckCollisionVertical(g.player.Sprite, g.colliders)

	activeAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnimation != nil {
		activeAnimation.Update()
	}

	for _, sprite := range g.enemies {

		sprite.Dx = 0
		sprite.Dy = 0

		if sprite.FollowsPlayer == false {
			continue
		}
		if sprite.X < g.player.X {
			sprite.Dx += 1
		}
		if sprite.X > g.player.X {
			sprite.Dx -= 1
		}
		if sprite.Y < g.player.Y {
			sprite.Dy += 1
		}
		if sprite.Y > g.player.Y {
			sprite.Dy -= 1
		}

		sprite.X += sprite.Dx
		CheckCollisionHorizontal(sprite.Sprite, g.colliders)

		sprite.Y += sprite.Dy
		CheckCollisionVertical(sprite.Sprite, g.colliders)
	}

	for _, potion := range g.potions {
		if potion.X == g.player.X && potion.Y == g.player.Y {
			g.player.Health += potion.AmtHeal
			fmt.Println("Player health is now: ", g.player.Health)
		}
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*constants.Tilesize,
		float64(g.tilemapJSON.Layers[0].Height)*constants.Tilesize,
		320,
		240,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	opts := ebiten.DrawImageOptions{}

	for layerIndex, layer := range g.tilemapJSON.Layers {
		for i, tileID := range layer.Data {

			if tileID == 0 {
				continue
			}

			x := i % layer.Width
			y := i / layer.Width

			x *= constants.Tilesize
			y *= constants.Tilesize

			img := g.tilesets[layerIndex].Img(tileID)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + constants.Tilesize))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(
				img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	// Draw the player image

	playerFrame := 0
	activeAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnimation != nil {
		playerFrame = activeAnimation.Frame()
	}

	screen.DrawImage(
		g.player.Img.SubImage(
			g.playerSpriteSheet.Rect(playerFrame),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, constants.Tilesize, constants.Tilesize)).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, constants.Tilesize, constants.Tilesize)).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.cam.X),
			float32(collider.Min.Y)+float32(g.cam.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{0xff, 0x00, 0x00, 0xff},
			true,
		)
	}
}
