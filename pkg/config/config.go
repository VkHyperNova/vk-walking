package config

import "path/filepath"

var DefaultContent = `{"walkings": []}`

var	file = "walkings.json"
var baseLocal = "DATABASES"
var	baseBackup = "/media/veikko/VK DATA/"

var LocalFile = filepath.Join(baseLocal, "WALKINGS", file)
var BackupFile = filepath.Join(baseBackup, "DATABASES", "WALKINGS", file)