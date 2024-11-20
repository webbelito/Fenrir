package utils

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func GetColorFromString(cN string) raylib.Color {

	switch cN {
	case "white":
		return raylib.White
	case "black":
		return raylib.Black
	case "gray":
		return raylib.Gray
	case "lightGray":
		return raylib.LightGray
	case "darkGray":
		return raylib.DarkGray
	case "yellow":
		return raylib.Yellow
	case "gold":
		return raylib.Gold
	case "orange":
		return raylib.Orange
	case "pink":
		return raylib.Pink
	case "red":
		return raylib.Red
	case "maroon":
		return raylib.Maroon
	case "green":
		return raylib.Green
	case "lime":
		return raylib.Lime
	case "darkGreen":
		return raylib.DarkGreen
	case "skyBlue":
		return raylib.SkyBlue
	case "blue":
		return raylib.Blue
	case "darkBlue":
		return raylib.DarkBlue
	case "purple":
		return raylib.Purple
	case "violet":
		return raylib.Violet
	case "darkPurple":
		return raylib.DarkPurple
	case "beige":
		return raylib.Beige
	case "brown":
		return raylib.Brown
	case "darkBrown":
		return raylib.DarkBrown
	default:
		return raylib.White
	}
}
