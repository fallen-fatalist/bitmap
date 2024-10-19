package flag

import (
	"fmt"

	"bitmap/utils"
)

// Public main variables
var (
	Command    string
	Arguments  []Argument
	SourceFile string
	OutputFile string
)

// Struct of apply command's arguments
type Argument struct {
	Name  string
	Value string
}

// Private variables
var (
	commands     = []string{"header", "apply"}
	helps        = []string{"-h", "--help", "help"}
	mirrorValues = []string{"h", "hor", "horizontal", "horizontally", "v", "ver", "vertical", "vertically"}
	filterValues = []string{"red", "green", "blue", "grayscale", "negative", "pixelate", "blur", "sepia"}
	rotateValues = []string{"right", "90", "180", "270", "left", "-90", "-180", "-270"}
)

// Errors
var (
	HelpCommand                   = fmt.Errorf("Help command called")
	ErrIncorrectNumberOfArguments = fmt.Errorf("Incorrect number of arguments")
	ErrTooManyArguments           = fmt.Errorf("Too many arguments entered")
	ErrTooLittleArguments         = fmt.Errorf("Too little arguments entered")
	ErrIncorrectCommandName       = fmt.Errorf("Incorrect command name")
	ErrIncorrectArgumentFormat    = fmt.Errorf("Incorrect argument(s) format, correct format: --<flag_name>=<value>")
	ErrIncorrectArgumentName      = fmt.Errorf("Incorrect argument(s) name")
	ErrIncorrectArgumentValue     = fmt.Errorf("Incorrect argument(s) value")
	ErrNotNumericArgumentValue    = fmt.Errorf("Argument(s) value is not numeric")
	ErrIncorrectOptionName        = fmt.Errorf("Incorrect option's name")
)

func Parse(args []string) error {
	// Get the command name
	if len(args) > 0 {
		Command = args[0]
		args = args[1:]
	}

	// Command validation
	switch {
	// help Case
	case utils.In(Command, helps) != -1 || Command == "":
		Command = ""
		return HelpCommand
	case Command == "header":
		if len(args) > 1 {
			return ErrTooManyArguments
		} else if len(args) == 0 || utils.In(args[0], helps) != -1 {
			return HelpCommand
		}

		SourceFile = args[0]
		return nil
	case Command == "apply":
		// Help command
		if utils.In(args[0], helps) != -1 {
			return HelpCommand
		} else if len(args) < 3 {
			return ErrIncorrectNumberOfArguments
		}

		// Arguments processing
		for _, arg := range args[:len(args)-2] {
			flagName, flagValue, err := getFlagNameAndValue("--", arg)
			if err != nil {
				return err
			}

			// Argument handling
			switch flagName {
			case "mirror":
				if utils.In(flagValue, mirrorValues) == -1 {
					return ErrIncorrectArgumentValue
				} else if utils.In(flagValue, []string{"h", "hor", "horizontal", "horizontally"}) != -1 {
					flagValue = "h"
				} else if utils.In(flagValue, []string{"v", "ver", "vertical", "vertically"}) != -1 {
					flagValue = "v"
				}
			case "filter":
				if utils.In(flagValue, filterValues) == -1 {
					return ErrIncorrectArgumentValue
				}
			case "rotate":
				if utils.In(flagValue, rotateValues) == -1 {
					return ErrIncorrectArgumentValue
				} else if flagValue == "right" {
					flagValue = "90"
				} else if flagValue == "left" {
					flagValue = "-90"
				}

			case "crop":
				// Size validation
				sizes := utils.Split(flagValue, "-")
				if len(sizes) != 2 && len(sizes) != 4 {
					return ErrIncorrectArgumentValue
				}
				// Numeric validation
				for _, size := range sizes {
					if !utils.IsNumeric(size) {
						return ErrNotNumericArgumentValue
					}
				}
			default:
				return ErrIncorrectOptionName
			}

			Arguments = append(Arguments, Argument{
				Name:  flagName,
				Value: flagValue,
			})
		}
		SourceFile = args[len(args)-2]
		OutputFile = args[len(args)-1]
		return nil
	default:
		Command = ""
		return ErrIncorrectCommandName
	}
}

// Returns the Flag name and the value of the flags with format: --<flag_name>=<value>
func getFlagNameAndValue(prefix, argument string) (flagName string, flagValue string, err error) {
	// Escape case when prefix has more length than argument
	if len(prefix) >= len(argument) {
		return "", "", ErrIncorrectArgumentFormat
	}
	endIdx := -1
	// checks that prefix present in argument
	for idx := range prefix {
		if prefix[idx] != argument[idx] {
			return "", "", ErrIncorrectArgumentFormat
		}
	}
	// get the flag name
	for idx, char := range argument {
		if char == '=' {
			endIdx = idx
			break
		}
	}

	// no = sign met
	if endIdx == -1 || len(prefix) >= endIdx {
		return "", "", ErrIncorrectArgumentFormat
	}

	flagName = argument[len(prefix):endIdx]
	flagValue = argument[endIdx+1:]
	return
}

func GetFlags() {
	fmt.Println("Command :", Command)
	fmt.Println("Arguments :", Arguments)
	fmt.Println("Source file :", SourceFile)
	fmt.Println("Output file :", OutputFile)
}

func PrintHelp() {
	fmt.Println("Usage:")
	if Command == "" {
		fmt.Println("   bitmap <command> [arguments]")
		fmt.Println()
		fmt.Println("The commands are:")
		fmt.Println("   header    prints bitmap file header information; add --help flag to get detailed information")
		fmt.Println("   apply     applies processing to the image and saves it to the file, add --help flag to get detailed information")
	} else if Command == "header" {
		fmt.Println("	bitmap header <source_file>")
		fmt.Println()
		fmt.Println("Description:")
		fmt.Println("   Prints bitmap file header information:")
		fmt.Println("	- filetype")
		fmt.Println("	- file size in bytes")
		fmt.Println("	- header size")
		fmt.Println("	- DIB header size")
		fmt.Println("	- width in pixels")
		fmt.Println("	- height in pixels")
		fmt.Println("	- pixel size in bits")
		fmt.Println("	- image size in bytesf")
		fmt.Println("	<source_file> must go last in the arguments list")
	} else if Command == "apply" {
		fmt.Println("   bitmap apply [options] <source_file> <output_file>")
		fmt.Println("	several options may be applied in the same time")
		fmt.Println()
		fmt.Println("The options of apply are:")
		fmt.Println("	--mirror : mirrors a bitmap image either horizontally or vertically; several mirrors may be applied in the provided sequence")
		fmt.Println("		possible values of --mirror: horizontal, h, horizontally, hor; vertical, v, vertically, ver")
		fmt.Println(" 		usage example: ./bitmap apply --mirror=horizontal sample.bmp sample-mirrored-horizontal.bmp")
		fmt.Println()
		fmt.Println("	--filter : applies some image filter on image; several filters may be applied in the provided sequence")
		fmt.Println("		possible values of --filter:")
		fmt.Println("		- blue 		: filter retains only the blue channel")
		fmt.Println("		- red  		: filter retains only the red channel")
		fmt.Println("		- green		: filter retains only the green channel")
		fmt.Println("		- grayscale	: filter converts the image to grayscale")
		fmt.Println("		- negative 	: applies a negative filter")
		fmt.Println("		- sepia		: applies a reddish brown color effect")
		fmt.Println("		- pixelate 	: apply a pixelation effect, option pixelates the image with a block of 20 pixels by default")
		fmt.Println("		- blur 		: applies a blur effect")
		fmt.Println("		usage example: ./bitmap apply --filter=blur sample.bmp sample-filtered-blur.bmp")
		fmt.Println()
		fmt.Println("	--rotate : rotates a bitmap image by a specified angle; several rotates may be applied in the provided sequence")
		fmt.Println("		possible values of --rotate:")
		fmt.Println("		- right : rotates clockwise")
		fmt.Println("		- 90 	: rotates image 90 degrees")
		fmt.Println("		- 180 	: rotates image 180 degrees")
		fmt.Println("		- 270 	: rorates image 270 degrees")
		fmt.Println("		- left 	: rotates counter-clockwise")
		fmt.Println("		- -90 	: rotates image -90 degrees")
		fmt.Println("		- -180  : rotates image -180 degrees")
		fmt.Println("		- -270 	: rorates image -270 degrees")
		fmt.Println("		usage example:  ./bitmap apply --rotate=right --rotate=right sample.bmp sample-rotated-right-right.bmp")
		fmt.Println()
		fmt.Println("	--crop : crop trims a bitmap image according to specified parameters; several crops may be applied in the provided sequence")
		fmt.Println("		crop flag accepts either 2 or 4 values in pixels")
		fmt.Println("		flag format: --crop=OffsetX-OffsetY-Width-Height, Width and Height are optional")
		fmt.Println("		usage example: ./bitmap apply --crop=20-20-100-100 sample.bmp sample-cropped-20-20-80-80.bmp")
		fmt.Println("	<source_file> <output_file> must go last in the arguments list")
	}
}
