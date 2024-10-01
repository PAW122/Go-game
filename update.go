package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func update() {
	runing = !rl.WindowShouldClose()

	var newX, newY float32
	newX = playerDest.X
	newY = playerDest.Y

	// Określ nową potencjalną pozycję gracza
	if playerUp {
		newY -= playerSpeed
	}
	if playerDown {
		newY += playerSpeed
	}
	if playerLeft {
		newX -= playerSpeed
	}
	if playerRight {
		newX += playerSpeed
	}

	// Obliczenie indeksu nowego położenia gracza w siatce
	gridX := int(newX / tileDest.Width)
	gridY := int(newY / tileDest.Height)

	// Upewnij się, że nie przekraczamy granic mapy
	if gridX >= 0 && gridX < mapW && gridY >= 0 && gridY < mapH {
		// Sprawdzenie kolizji
		if !isCollision(gridX, gridY) {
			playerDest.X = newX
			playerDest.Y = newY
		}
	}

	frameCount++
	// zmień klatkę animacji co x klatek
	if playerMoving {
		if frameCount >= animationSpeed {
			playerFrame++
			if playerFrame >= animationFrames { // Zakładamy 4 klatki animacji
				playerFrame = 0
			}
			frameCount = 0
		}

		frameCountIdle = 0 // Resetujemy licznik klatek idle
	} else {
		frameCount++
		if frameCount%45 == 1 {
			playerFrame++
			if playerFrame > 1 {
				playerFrame = 0
			}
		}
	}

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPasued {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(
		float32(playerDest.X-(playerDest.Width/2)),
		float32(playerDest.Y-(playerDest.Height/2)),
	)

	playerMoving = false
	playerUp, playerDown, playerRight, playerLeft = false, false, false, false
}
