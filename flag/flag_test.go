package flag

import "testing"

func TestParse(t *testing.T) {
	type testData struct {
		name       string
		args       []string
		err        error
		command    string
		outputArgs []Argument
		sourceFile string
		outputFile string
	}

	tests := []testData{
		{
			name: "Help command",
			args: []string{"--help"},
			err:  HelpCommand,
		},
		{
			name: "Help command with other args",
			args: []string{"--help", "other"},
			err:  ErrIncorrectNumberOfArguments,
		},
		{
			name:    "Incorrect command name",
			args:    []string{"--wrong"},
			err:     ErrIncorrectCommandName,
			command: "",
		},
		{
			name:    "Apply command without arguments",
			args:    []string{"apply", "source_file", "output_file"},
			err:     ErrIncorrectNumberOfArguments,
			command: "apply",
		},
		{
			name:    "Apply command with incorrect number of args",
			args:    []string{"apply", "source_file"},
			err:     ErrIncorrectNumberOfArguments,
			command: "apply",
		},
		{
			name:       "Apply command with mirror flag",
			args:       []string{"apply", "--mirror=horizontal", "source_file", "output_file"},
			outputArgs: []Argument{{Name: "mirror", Value: "h"}},
			sourceFile: "source_file",
			outputFile: "output_file",
			command:    "apply",
		},
		{
			name:       "Apply command with filter flag",
			args:       []string{"apply", "--filter=blur", "source_file", "output_file"},
			outputArgs: []Argument{{Name: "filter", Value: "blur"}},
			sourceFile: "source_file",
			outputFile: "output_file",
			command:    "apply",
		},
		{
			name:       "Apply command with several filter arguments",
			args:       []string{"apply", "--filter=blur", "--filter=blue", "source_file", "output_file"},
			outputArgs: []Argument{{Name: "filter", Value: "blur"}, {Name: "filter", Value: "blue"}},
			sourceFile: "source_file",
			outputFile: "output_file",
			command:    "apply",
		},
		{
			name:    "Apply without output file",
			args:    []string{"apply", "--rotate=90", "source_file"},
			err:     ErrIncorrectNumberOfArguments,
			command: "apply",
		},
		{
			name:       "Header command with source file",
			args:       []string{"header", "source_file"},
			sourceFile: "source_file",
			command:    "header",
		},
		{
			name:    "Header command with source file and help",
			args:    []string{"header", "source_file", "--help"},
			err:     ErrTooManyArguments,
			command: "header",
		},
		{
			name:    "Header command with incorrect number of args",
			args:    []string{"header"},
			err:     HelpCommand,
			command: "header",
		},
		{
			name:       "Several options called",
			args:       []string{"apply", "--filter=blur", "--rotate=90", "source_file", "output_file"},
			outputArgs: []Argument{{Name: "filter", Value: "blur"}, {Name: "rotate", Value: "90"}},
			err:        nil,
			command:    "apply",
			sourceFile: "source_file",
			outputFile: "output_file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Parse(test.args)
			if err != test.err {
				t.Errorf("Parse() error = %v, wantErr %v", err, test.err)
			} else if Command != test.command {
				t.Errorf("Parse() command = %v, want %v", Command, test.command)
			} else if SourceFile != test.sourceFile {
				t.Errorf("Parse() sourceFile = %v, want %v", SourceFile, test.sourceFile)
			} else if OutputFile != test.outputFile {
				t.Errorf("Parse() outputFile = %v, want %v", OutputFile, test.outputFile)
			} else if len(test.outputArgs) != len(Arguments) {
				t.Errorf("Parse() Arguments = %v, want %v", Arguments, test.outputArgs)
			}

			for idx, arg := range test.outputArgs {
				if idx < len(Arguments) && arg.Name != Arguments[idx].Name {
					t.Errorf("Parse() Arguments = %v, want %v", Arguments, test.outputArgs)
				} else if idx < len(Arguments) && arg.Value != Arguments[idx].Value {
					t.Errorf("Parse() Arguments = %v, want %v", Arguments, test.outputArgs)
				}
			}
			Arguments = nil
			Command = ""
			SourceFile = ""
			OutputFile = ""
		})
	}
}
