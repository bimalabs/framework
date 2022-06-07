package interfaces

import (
	"fmt"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/fatih/color"
)

type Database struct {
}

func (d *Database) Run(servers []configs.Server) {
	if len(servers) > 0 {
		color.New(color.FgCyan, color.Bold).Printf("✓ ")
		fmt.Println("Serving DB Auto Migration Juices...")
	}

	for _, server := range servers {
		go server.RegisterAutoMigrate()
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
