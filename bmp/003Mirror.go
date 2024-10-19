package bmp

import "errors"

// Errors
var (
	ErrIncorrectMirrorValue = errors.New("Incorrect mirror value provided.")
)

// Mirrors the image
func (b *bmp) Mirror(flagValue string) error {
	newPixelArray := make([][]byte, 0, len(b.pixelArray))

	switch flagValue {
	case "h":
		// Mirror rows along the horizontal line
		for idx := len(b.pixelArray) - 1; idx >= 0; idx-- {
			newPixelArray = append(newPixelArray, b.pixelArray[idx])
		}
	case "v":
		// Initialize new pixel array
		rowSize := (uint32(b.dibHeader.BitsPerPixel)*(b.dibHeader.Width) + 31) / 32 * 4
		for idx := 0; idx < len(b.pixelArray); idx++ {
			newPixelArray = append(newPixelArray, make([]byte, rowSize))
		}

		// Mirror pixels along the vertical line
		for colIdx := 0; uint32(colIdx) < rowSize; colIdx += 3 {
			for rowIdx := 0; rowIdx < len(b.pixelArray); rowIdx++ {

				newPixelArray[rowIdx][colIdx] = b.pixelArray[rowIdx][rowSize-uint32(colIdx)-3]   // blue
				newPixelArray[rowIdx][colIdx+1] = b.pixelArray[rowIdx][rowSize-uint32(colIdx)-2] // green
				newPixelArray[rowIdx][colIdx+2] = b.pixelArray[rowIdx][rowSize-uint32(colIdx)-1] // red
			}
		}
	default:
		return ErrIncorrectMirrorValue
	}

	b.pixelArray = newPixelArray
	return nil
}
