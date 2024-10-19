package bmp

import (
	"errors"
)

// Errors
var (
	ErrIncorrectFilterValue = errors.New("Incorrect value provided to Filter option")
	ErrIndexOutOfBound      = errors.New("Index out of bounds of pixel array")
)

const (
	eps = 0.01
)

func (b *bmp) Filter(flagValue string) error {
	switch flagValue {
	case "red":
		for rowIdx := range b.pixelArray {
			// Nullify blue and green colors
			for colIdx := uint32(0); colIdx < rowSize; colIdx += 3 {
				b.pixelArray[rowIdx][colIdx] = 0   // Blue
				b.pixelArray[rowIdx][colIdx+1] = 0 // Green
			}
		}
	case "green":
		for rowIdx := range b.pixelArray {
			// Nullify red and blue colors
			for colIdx := uint32(0); colIdx < rowSize; colIdx += 3 {
				b.pixelArray[rowIdx][colIdx] = 0   // Blue
				b.pixelArray[rowIdx][colIdx+2] = 0 // Red
			}
		}
	case "blue":
		for rowIdx := range b.pixelArray {
			// Nullify red and green colors
			for colIdx := uint32(0); colIdx < rowSize; colIdx += 3 {
				b.pixelArray[rowIdx][colIdx+1] = 0 // Green
				b.pixelArray[rowIdx][colIdx+2] = 0 // Red
			}
		}
	case "grayscale":
		// Weighted method
		// see (https://idmnyu.github.io/p5.js-image/Filters/index.html)
		for rowIdx := range b.pixelArray {
			for colIdx := uint32(0); colIdx+2 < rowSize; colIdx += 3 {
				// Scale color value according to weight
				RedVal := float32(b.pixelArray[rowIdx][colIdx]) * 0.11
				GreenVal := float32(b.pixelArray[rowIdx][colIdx+1]) * 0.59
				BlueVal := float32(b.pixelArray[rowIdx][colIdx+2]) * 0.3
				// Sum colors to get gray color
				LumaVal := byte(RedVal + GreenVal + BlueVal)
				// Assign color to pixel
				b.pixelArray[rowIdx][colIdx] = LumaVal   // Blue
				b.pixelArray[rowIdx][colIdx+1] = LumaVal // Green
				b.pixelArray[rowIdx][colIdx+2] = LumaVal // Red
			}
		}
	case "negative":
		for rowIdx := range b.pixelArray {
			for colIdx := uint32(0); colIdx+2 < rowSize; colIdx += 3 {
				// Invert color values
				b.pixelArray[rowIdx][colIdx] = 255 - b.pixelArray[rowIdx][colIdx]
				b.pixelArray[rowIdx][colIdx+1] = 255 - b.pixelArray[rowIdx][colIdx+1]
				b.pixelArray[rowIdx][colIdx+2] = 255 - b.pixelArray[rowIdx][colIdx+2]
			}
		}
	case "sepia":
		// Microsoft recommended values
		// see (https://idmnyu.github.io/p5.js-image/Filters/index.html)
		for rowIdx := range b.pixelArray {
			for colIdx := uint32(0); colIdx+2 < rowSize; colIdx += 3 {
				// Get color values
				blueColor, greenColor, redColor := b.pixelArray[rowIdx][colIdx], b.pixelArray[rowIdx][colIdx+1], b.pixelArray[rowIdx][colIdx+2]

				// Apply ratio to colors
				blueSepia := float32(redColor)*.272 + float32(greenColor)*.534 + float32(blueColor)*.131
				// condition of exceeding 255
				if blueSepia-255. > eps {
					blueSepia = 255.
				}
				greenSepia := float32(redColor)*.349 + float32(greenColor)*.686 + float32(blueColor)*.168
				// condition of exceeding 255
				if greenSepia-255 > eps {
					greenSepia = 255.
				}
				redSepia := float32(redColor)*.393 + float32(greenColor)*.769 + float32(blueColor)*.189
				// condition of exceeding 255
				if redSepia-255 > eps {
					redSepia = 255.
				}

				// Assign new colors
				b.pixelArray[rowIdx][colIdx] = byte(blueSepia)
				b.pixelArray[rowIdx][colIdx+1] = byte(greenSepia)
				b.pixelArray[rowIdx][colIdx+2] = byte(redSepia)
			}
		}
	case "pixelate":
		
	case "blur":
		// Box 5x5 blur method
		// kernel related weight of color sums
		// see (https://idmnyu.github.io/p5.js-image/Blur/index.html)

		// Initialize blurred pixel array
		blurredPixelArray := make([][]byte, len(b.pixelArray))
		for rowIdx := range blurredPixelArray {
			blurredPixelArray[rowIdx] = make([]byte, rowSize)
		}

		// Iterate over pixels and calculate color sums in box and assign sums to pixel
		for rowIdx := uint32(0); rowIdx < uint32(len(b.pixelArray)); rowIdx++ {
			for colIdx := uint32(0); colIdx+2 < rowSize; colIdx += 3 {
				// Fetch the sum in the box
				sum_B, sum_G, sum_R, pixelCount := b.boxPixelsSum(5, int(rowIdx), int(colIdx))

				// Assign center pixel average color value of all surrounding pixels including itself
				blurredPixelArray[rowIdx][colIdx] = byte(sum_B / uint16(pixelCount))   // Blue
				blurredPixelArray[rowIdx][colIdx+1] = byte(sum_G / uint16(pixelCount)) // Green
				blurredPixelArray[rowIdx][colIdx+2] = byte(sum_R / uint16(pixelCount)) // Red
			}
		}
		// Assign blurred array
		b.pixelArray = blurredPixelArray

	default:
		return ErrIncorrectFilterValue
	}

	return nil
}

// boxPixelsSum returns the sum of color values in the box
// sum_B - sum of blue color values in the box
// sum_G - sum of green color values in the box
// sum_R - sum of red color values in the box
func (b *bmp) boxPixelsSum(kernelSize, rowIdx, colIdx int) (sum_B, sum_G, sum_R uint16, neighborPixelCount int) {
	kernelOffset := kernelSize / 2

	for dx := -kernelOffset; dx < kernelOffset; dx++ {
		for dy := -kernelOffset; dy < kernelOffset; dy++ {
			// Fetching pixel
			pixel, err := b.fetchPixel(rowIdx+dx, colIdx+dy*3)
			if err != nil {
				continue
			}

			// Adding pixel color value to sum
			sum_B += uint16(pixel[0]) // Blue
			sum_G += uint16(pixel[1]) // Green
			sum_R += uint16(pixel[2]) // Red
			// Increment pixel count
			neighborPixelCount++
		}
	}

	return
}

// fetchPixel returns the slice of pixel's values
func (b *bmp) fetchPixel(rowIdx, colIdx int) ([]byte, error) {
	if rowIdx < 0 || rowIdx >= len(b.pixelArray) || colIdx < 0 || colIdx+3 >= int(rowSize) {
		return nil, ErrIndexOutOfBound
	}

	return b.pixelArray[rowIdx][colIdx : colIdx+3], nil
}
