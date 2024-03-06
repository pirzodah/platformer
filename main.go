package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	screenWidth  = 2000
	screenHeight = 700
)

var (
	playerImage     *ebiten.Image
	backgroundImage *ebiten.Image
	playerX         = 0
	playerY         = 0
	jumping         = false
	jumpVelocity    = 0
	gravity         = 3
	groundLevel     = 630

	obstacleX      = 250 // Сдвиг ящика вправо
	obstacleY      = 350 // Установка новой координаты Y ящика
	obstacleWidth  = 40  // Уменьшаем ширину ящика в два раза
	obstacleHeight = 40  // Уменьшаем высоту ящика в два раза

	obstacle2X      = 460 // Координата X второго препятствия
	obstacle2Y      = 350 // Координата Y второго препятствия
	obstacle2Width  = 40  // Ширина второго препятствия
	obstacle2Height = 40  // Высота второго препятствия

)

type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return nil // Закрываем игру, возвращая nil
	}

	// Управление игроком
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		playerX -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		playerX += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if !jumping {
			jumping = true
			jumpVelocity = 50
		}
	}

	// Прыжок игрока
	if jumping {
		playerY -= jumpVelocity
		jumpVelocity -= gravity
		if playerY > groundLevel {
			playerY = groundLevel
			jumping = false
		}
	} else if playerY < groundLevel {
		// Установка playerY на уровень земли только если персонаж не в состоянии прыжка и находится выше уровня земли
		playerY = groundLevel
	}

	// Ограничение перемещения игрока в пределах экрана
	if playerX < 0 {
		playerX = 0
	}
	if playerX > screenWidth-16 {
		playerX = screenWidth - 16
	}
	if playerY < 0 {
		playerY = 0
	}
	if playerY > screenHeight-16 {
		playerY = screenHeight - 16
	}

	// Проверка коллизий с препятствием
	if checkCollision(playerX, playerY, 16, 16, obstacleX, obstacleY, obstacleWidth, obstacleHeight) {
		// Обработка столкновения
		playerX = obstacleX - 16 // Устанавливаем персонажа перед препятствием
		playerY = obstacleY - 16 // Устанавливаем персонажа на верхний край препятствия
	}

	if checkCollision(playerX, playerY, 16, 16, obstacle2X, obstacle2Y, obstacle2Width, obstacle2Height) {
		// Обработка столкновения
		playerX = obstacle2X - 16 // Устанавливаем персонажа перед вторым препятствием
		playerY = obstacle2Y - 16 // Устанавливаем персонажа на верхний край второго препятствия
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Отображение фонового изображения с уменьшением размера в два раза
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.8, 0.7) // Уменьшение размера в два раза
	screen.DrawImage(backgroundImage, op)

	// Рисование игрока с уменьшением размера в два раза
	geoM := ebiten.GeoM{}
	geoM.Translate(float64(playerX), float64(playerY))
	geoM.Scale(0.4, 0.4) // Уменьшение размера в два раза
	screen.DrawImage(playerImage, &ebiten.DrawImageOptions{GeoM: geoM})

	// Рисование препятствия 1
	vector.DrawFilledRect(screen, float32(obstacleX), float32(obstacleY), float32(obstacleWidth), float32(obstacleHeight), color.Gray{}, false)

	// Рисование препятствия 2
	vector.DrawFilledRect(screen, float32(obstacle2X), float32(obstacle2Y), float32(obstacle2Width), float32(obstacle2Height), color.Gray{}, false)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Загрузка изображения игрока
	img, _, err := ebitenutil.NewImageFromFile("1594832828.png")
	if err != nil {
		panic(err)
	}
	playerImage = img

	// Загрузка фонового изображения
	backgroundImage, _, err = ebitenutil.NewImageFromFile("Game_Background_190.png")
	if err != nil {
		panic(err)
	}

	playerX = 50 // Персонаж начинает игру немного правее

	// Запуск игры
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

// Вспомогательная функция для проверки коллизий
func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}
