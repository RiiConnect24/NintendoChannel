package info

import (
	"NintendoChannel/gametdb"
	"unicode/utf16"
)

type SupportedControllers struct {
	WiiRemote          uint8
	Nunchuk            uint8
	ClassicController  uint8
	GamecubeController uint8
}

func (i *Info) GetSupportedControllers(controllers *gametdb.Controllers) {
	wrotePeripheral := false
	// For some reason the peripheral text must be padded with 2 uint16 before any real text.
	temp := []uint16{0, 0}
	for _, s := range controllers.Controller {
		switch s.Type {
		case "wiimote":
			i.SupportedControllers.WiiRemote = 1
			break
		case "nunchuk":
			i.SupportedControllers.Nunchuk = 1
			break
		case "classiccontroller":
			i.SupportedControllers.ClassicController = 1
			break
		case "gamecube":
			i.SupportedControllers.GamecubeController = 1
			break
		case "mii":
			// Mii's aren't a controller, but they are considered one to GameTDB for some reason
			i.SupportedFeatures.Miis = 1
			break
		case "wheel":
			temp = append(temp, utf16.Encode([]rune("Wii Wheel, "))...)
			wrotePeripheral = true
		case "balanceboard":
			temp = append(temp, utf16.Encode([]rune("Wii Balance Board, "))...)
			wrotePeripheral = true
		case "wiispeak":
			temp = append(temp, utf16.Encode([]rune("Wii Speak, "))...)
			wrotePeripheral = true
		case "microphone":
			temp = append(temp, utf16.Encode([]rune("Microphone, "))...)
			wrotePeripheral = true
		case "guitar":
			temp = append(temp, utf16.Encode([]rune("Guitar, "))...)
			wrotePeripheral = true
		case "drums":
			temp = append(temp, utf16.Encode([]rune("Drums, "))...)
			wrotePeripheral = true
		case "dancepad":
			temp = append(temp, utf16.Encode([]rune("Dance Pad, "))...)
			wrotePeripheral = true
		case "keyboard":
			temp = append(temp, utf16.Encode([]rune("Keyboard, "))...)
			wrotePeripheral = true
		case "udraw":
			temp = append(temp, utf16.Encode([]rune("uDraw, "))...)
			wrotePeripheral = true
		case "amiibo":
			temp = append(temp, utf16.Encode([]rune("Amiibo, "))...)
			wrotePeripheral = true
		}
	}

	temp = temp[:len(temp)-2]

	if wrotePeripheral {
		copy(i.PeripheralsText[:], temp)
	}
}