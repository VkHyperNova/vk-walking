package main

import (
	"vk-walking/pkg/cmd"
	"vk-walking/pkg/util"
)

func main() {
	
	util.CreateDDriveFiles()
	util.CreateLocalFiles()

	// Start
	cmd.CommandLine()
}
