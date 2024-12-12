package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	eqOpenCooldown int16 = eqFramesCooldown
)

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}

	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}

	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}

	if rl.IsKeyPressed(rl.KeyP) {
		musicPasued = !musicPasued
	}

	// open / close eq
	if rl.IsKeyDown(rl.KeyE) && eqOpenCooldown < 1 {
		eqOpenCooldown = 15
		eqOpen = !eqOpen
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		playerLeftClick = true
	} else {
		playerLeftClick = false
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		playerRightClick = true
	} else {
		playerRightClick = false
	}

	if eqOpenCooldown < -30768 {
		eqOpenCooldown = 0
	}

	eqOpenCooldown--
}
