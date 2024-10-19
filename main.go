package main

import (
	"bitmap/bmp"
	"bitmap/flag"
	"fmt"
	"os"
)

func main() {
	// Flag proccessing
	err := flag.Parse(os.Args[1:])
	if err != nil {
		if err == flag.HelpCommand {
			flag.PrintHelp()
			return
		}

		fmt.Fprintf(os.Stderr, "Error parsing flags: %s; to get help add --help flag.\n", err)
		os.Exit(1)
	}

	bmpFile, err := bmp.Load(flag.SourceFile)
	if err != nil {
		if err == bmp.ErrNon24BitImageNotSupported {

		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	switch flag.Command {
	case "header":
		bmpFile.PrintHeader()
		return
	case "apply":
		// Processing validation
		if bmpFile.GetPixelNumber() != 24 {
			fmt.Fprintf(os.Stderr, "File: %s is not 24 bit color pallete", flag.SourceFile)
			os.Exit(1)
		}
		// Arguments proccessing
		for _, arg := range flag.Arguments {
			switch arg.Name {
			case "mirror":
				err := bmpFile.Mirror(arg.Value)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while Mirroring the BMP image: %s.\n", err)
					os.Exit(1)
				}
			case "filter":
				err := bmpFile.Filter(arg.Value)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error while Filtering the BMP image: %s.\n", err)
					os.Exit(1)
				}
			case "crop":
				return
			case "rotate":
				return
			}
		}
	}

	bmpFile.Save(flag.OutputFile)
}
