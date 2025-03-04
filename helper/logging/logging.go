package logging

import (
	"gcw/helper/cl"
	"log"
	"os"
)

func getDirectory() string {
	path, _ := os.Getwd()
	return path + "."
}

func Low(location, message, errorstatus string) {
	log.Print(cl.Blue("Error:"), cl.Magenta(getDirectory()+location), ":", cl.Cyan(message), ": ", errorstatus)
}

func Warn(location, message, errorstatus string) {
	log.Print(cl.Yellow("WARNING:"), cl.Magenta(getDirectory()+location), ":", cl.Cyan(message), ": ", errorstatus)
}

func High(location, message, errorstatus string) {
	log.Print(cl.Red("RISK:"), cl.Magenta(getDirectory()+location), ":", cl.Cyan(message), ": ", errorstatus)
}
