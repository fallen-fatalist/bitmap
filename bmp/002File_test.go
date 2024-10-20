package bmp

import (
	"encoding/binary"
	"os"
	"testing"
)

func TestLoadSave(t *testing.T) {
	type testData struct {
		name       string
		sourceFile string
		outputFile string
	}

	tests := []testData{
		{
			name:       "54 bytes header size BMP file",
			sourceFile: "../samples/sample.bmp",
			outputFile: "../samples/sample-saved.bmp",
		},
		{
			name:       "138 bytes header size BMP file",
			sourceFile: "../samples/sample_640x426.bmp",
			outputFile: "../samples/sample_640x426-saved.bmp",
		},
		{
			name:       "138 bytes header HD resolution size BMP file",
			sourceFile: "../samples/sample_1280x853.bmp",
			outputFile: "../samples/sample_1280x853-saved.bmp",
		},
		{
			name:       "138 bytes header FullHD resolution size BMP file",
			sourceFile: "../samples/sample_1920x1280.bmp",
			outputFile: "../samples/sample_1920x1280-saved.bmp",
		},
		{
			name:       "138 bytes header 5K resolution size BMP file",
			sourceFile: "../samples/sample_5184x3456.bmp",
			outputFile: "../samples/sample_5184x3456-saved.bmp",
		},
		{
			name:       "Square bmp image",
			sourceFile: "../samples/marilyn.bmp",
			outputFile: "../samples/marilyn-saved.bmp",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testBmp, err := Load(test.sourceFile)
			if err != nil {
				t.Fatalf("Error while loading %s: %s\n", test.sourceFile, err)
			}

			err = testBmp.Save(test.outputFile)
			if err != nil {
				t.Fatalf("Error while saving %s: %s\n", test.outputFile, err)
			}

			// Compare saved and original file:
			sampleFile, err := os.Open(test.sourceFile)
			if err != nil {
				t.Fatalf("Error while opening %s: %s\n", test.sourceFile, err)
			}
			defer sampleFile.Close()
			savedFile, err := os.Open(test.outputFile)
			if err != nil {
				t.Fatalf("Error while opening %s: %s\n", test.outputFile, err)
			}
			defer os.Remove(test.outputFile)
			defer savedFile.Close()

			// Allocate buffers for file
			OriginalBuf := make([]byte, testBmp.fileHeader.FileSize)
			SavedBuf := make([]byte, testBmp.fileHeader.FileSize)

			// Free memory
			testBmp = nil

			err = binary.Read(sampleFile, binary.LittleEndian, &OriginalBuf)
			if err != nil {
				t.Fatalf("Error while reading %s: %s\n", test.sourceFile, err)
			}
			err = binary.Read(savedFile, binary.LittleEndian, &SavedBuf)
			if err != nil {
				t.Fatalf("Error while reading %s: %s\n", test.outputFile, err)
			}

			if len(OriginalBuf) != len(SavedBuf) {
				t.Fatalf("Error while comparing %s and %s\n, sizes differ.", test.sourceFile, test.outputFile)
			}

			for idx, val := range OriginalBuf {
				if val != SavedBuf[idx] {
					t.Fatalf("Error while comparing %s and %s\n, images differ.", test.sourceFile, test.outputFile)
				}
			}
		})
	}
}
