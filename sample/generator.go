package sample

import (
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}

	return keyboard
}


func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	cpuName := randomCPUName(brand)
	numberOfCores := randomInt(2, 8)
	numberOfThreads := randomInt(numberOfCores, 12)
	minGhz := randomFloat64(1.0, 3.6)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &pb.CPU{
		Brand: brand,
		Name: cpuName,
		NumberOfCores: uint32(numberOfCores),
		NumberOfThreads: uint32(numberOfThreads),
		MinGhz: minGhz,
		MaxGhz: maxGhz,
	}

	return cpu
}

func NewGPU() *pb.GPU {
	brand := randomGPUBrand();
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.6)
	maxGhz := randomFloat64(minGhz, 3.0)

	memory := &pb.Memory {
		Value: uint64(randomInt(2, 6)),
		Unit: pb.Memory_GIGABYTE,
	}

	gpu := &pb.GPU{
		Brand: brand,
		Name: name,
		MinGhz: minGhz,
	 	MaxGhz: maxGhz,
	 	Memory: memory,
	}
	return gpu
}

func NewRAM() *pb.Memory {
	ram := &pb.Memory {
		Value: uint64(randomInt(4, 64)),
		Unit: pb.Memory_GIGABYTE,
	}

	return ram
}

func NewSSD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory {
			Value: uint64(randomInt(128, 1024)),
			Unit: pb.Memory_GIGABYTE,
		},
	}
	return ssd
}

func NewHDD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory {
			Value: uint64(randomInt(1, 6)),
			Unit: pb.Memory_TERABYTE,
		},
	}
	return ssd
}

func NewScreen() *pb.Screen {

	height := randomInt(1080, 4230)
	width := height * 16/9

	screen := &pb.Screen{
		SizeInch:   randomFloat32(13, 17),
		Panel: randomScreenPanel(),
		Resolution:      &pb.Screen_Resolution{
			Height: uint32(height),
			Width: uint32(width),
		},
		Multitouch: randomBool(),
	}

	return screen
}

func NewLaptop() *pb.Laptop {
	laptop := &pb.Laptop{
		Id:          uuid.New().String(),
		Brand:       "Apple",
		Name: 		"Macbook Pro",
		Cpu:         NewCPU(),
		Ram:         NewRAM(),
		Gpus:        []*pb.GPU{NewGPU()},
		Storages:    []*pb.Storage{NewSSD(), NewHDD()},
		Screen:       NewScreen(),
		Keyboard:    NewKeyboard(),
		Weight:      &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1,3),
		},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:  timestamppb.Now(),
	}

	return laptop
}