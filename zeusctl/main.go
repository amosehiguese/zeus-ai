package main

import (
	"github.com/amosehiguese/zeus-ai/zeusctl/ctlv1"
	"github.com/fatih/color"
)

func main() {
	color.NoColor = false
	ctlv1.MustStart()
}
