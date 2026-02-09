package config

import "time"

const FileName = "walkings.json"
const FolderName = "WALKINGS"
const BackupFolder = "/media/veikko/VK DATA/DATABASES/"

var DDriveSave = false

var Date = time.Now().Format("02.01.2006")

// var Questions = []string{"Name:", "Distance:", "Duration:", "Pace:", "Steps:", "Calories:", "Date:"}
// var AddSuggestions = []string{"", "", "", "", "", "", ""}

var LocalPath = "./" + FolderName + "/" + FileName
var BackupPath = BackupFolder + FolderName + "/" + FileName
var BackupPathWithDate = BackupFolder + FolderName + "/" + "(" + Date + ") " + FileName
