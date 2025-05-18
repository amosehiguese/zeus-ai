package main

import (
	"github.com/fatih/color"

	"github.com/amosehiguese/zeus-ai/zeusctl/ctlv1"
)

func main() {
	color.NoColor = false
	ctlv1.MustStart()
}
