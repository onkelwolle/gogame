package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/onkelwolle/gogame/entities"
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {

	//react to key presses

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X += 2
	}

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer == false {
			continue
		}
		if sprite.X < g.player.X {
			sprite.X += 1
		}
		if sprite.X > g.player.X {
			sprite.X -= 1
		}
		if sprite.Y < g.player.Y {
			sprite.Y += 1
		}
		if sprite.Y > g.player.Y {
			sprite.Y -= 1
		}
	}

	for _, potion := range g.potions {
		if potion.X == g.player.X && potion.Y == g.player.Y {
			g.player.Health += potion.AmtHeal
			fmt.Println("Player health is now: ", g.player.Health)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.tilemapJSON.Layers {
		for i, tileID := range layer.Data {
			x := i % layer.Width
			y := i / layer.Width

			x *= 16
			y *= 16

			srcX := (tileID - 1) % 22
			srcY := (tileID - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(
					image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)

	// Draw the player image
	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

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

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   100,
				Y:   100,
			},
			Health: 100,
		},
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
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
