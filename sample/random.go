package sample

import (
	"github.com/Adetunjii/go-grpc/pb"
	"math/rand"
)

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(a ...string) string {
	n := len(a);
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomInt(min, max int32) int32 {
	return min + rand.Int31n(max - min + 1)
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_OPS
	}

	return pb.Screen_OLED
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64() * (max - min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32() * (max - min)
}

func randomCPUName(brand string) string {
	if brand == "intel" {
		return randomStringFromSet("Pentium","Core i3", "Core i5", "Core i7", "Core i9")
	}

	return randomStringFromSet("Ryzen 7 PRO 2700U", "Ryzen 5 PRO 3500U")
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet("RTX 2060", "RTX 2070", "GTX 1660-Ti")
	}

	return randomStringFromSet("RX 580", "RX 590")
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomBool() bool {
	return rand.Intn(2) == 1
}
