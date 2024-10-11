package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	_ "image/webp"

	"github.com/fatih/color"
	"github.com/PontusTideman/meme/cli"
	"github.com/PontusTideman/meme/image"
	"github.com/PontusTideman/meme/output"
)

func main() {
	opt := cli.ParseOptions()

	if opt.Help {
		opt.PrintUsage()

	} else if opt.ListTemplates {
		for _, id := range cli.ImageIds {
			fmt.Fprintln(output.Stdout, color.CyanString("%s", id))
		}

	} else if opt.Valid() {
		st := image.Load(opt)
		st = image.RenderImage(opt, st)

		if opt.ClientID != "" {
			url := image.Upload(opt, st)
			output.Info(url)
		} else {
			file := image.Save(opt, st)
			output.Info(file)
		}
	}
}
