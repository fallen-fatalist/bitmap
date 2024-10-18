package bmp

import (
	"encoding/binary"
	"errors"
	"os"
)

// Errors
var (
	ErrIncorrectSignature        = errors.New("File is not in BMP format, or file is corrupted.")
	ErrFileIsCorrupted           = errors.New("BMP File is corrupted.")
	ErrIncorrectFileFormat       = errors.New("File's format does not match BMP format.")
	ErrNon24BitImageNotSupported = errors.New("Image with 24 bit color pallete is not supported.")
)

// File header of Device independent bitmap
// see (https://en.wikipedia.org/wiki/BMP_file_format#Bitmap_file_header)
type fileHeader struct {
	Signature uint16
	FileSize  uint32
	Reserved1 uint16
	Reserved2 uint16
	Offset    uint32
}

// Buffs where unused part of DIB header or ICC profile will be recorded (all metadata which differs from BITMAPINFOHEADER)
// see (https://en.wikipedia.org/wiki/BMP_file_format#DIB_header_(bitmap_information_header))
var (
	unusedBuf1 []byte
	unusedBuf2 []byte
)

// BMP file reading
func Load(fileName string) (*bmp, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bmp := &bmp{}

	// File header reading
	bmp.fileHeader = &fileHeader{}

	if err := binary.Read(file, binary.LittleEndian, bmp.fileHeader); err != nil {
		return nil, err
	}

	// validation for bmp signature
	if bmp.fileHeader.Signature != 19778 {
		return nil, ErrIncorrectFileFormat
	}

	// Reading DIB header
	bmp.dibHeader = &dibHeader{}

	if err := binary.Read(file, binary.LittleEndian, bmp.dibHeader); err != nil {
		return nil, err
	}

	// Reading unused bytes until pixel array
	unusedBuf1 = make([]byte, bmp.fileHeader.Offset-54)
	if err := binary.Read(file, binary.LittleEndian, &unusedBuf1); err != nil {
		return nil, err
	}

	// If not 24 bit, just skip the pixel array
	// it's done to make possible reading header of images which color pallete is not 24
	if bmp.dibHeader.BitsPerPixel != 24 {
		return bmp, ErrNon24BitImageNotSupported
	}

	// Reading  pixel array

	// row size is the number of pixels in one row by row
	rowSize := (uint32(bmp.dibHeader.BitsPerPixel)*(bmp.dibHeader.Width) + 31) / 32 * 4
	bmp.pixelArray = make([][]byte, bmp.dibHeader.Height)

	// reading pixel array row by row
	for idx, row := range bmp.pixelArray {
		row = make([]byte, rowSize)
		if err := binary.Read(file, binary.LittleEndian, &row); err != nil {
			return nil, err
		}
		bmp.pixelArray[idx] = row
	}

	// Reading unused bytes after pixel array
	unusedBuf2 = make([]byte, bmp.fileHeader.FileSize-bmp.dibHeader.ImageSize-uint32(len(unusedBuf1))-54)
	if err := binary.Read(file, binary.LittleEndian, &unusedBuf2); err != nil {
		return nil, err
	}

	return bmp, nil
}

func (b *bmp) Save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writing file header
	if err := binary.Write(file, binary.LittleEndian, b.fileHeader); err != nil {
		return err
	}

	// Writing DIB header
	if err := binary.Write(file, binary.LittleEndian, b.dibHeader); err != nil {
		return err
	}

	// Writing unused bytes before pixel array
	if err := binary.Write(file, binary.LittleEndian, unusedBuf1); err != nil {
		return err
	}

	// Writing pixel array row by row
	for _, row := range b.pixelArray {
		if err := binary.Write(file, binary.LittleEndian, row); err != nil {
			return err
		}
	}

	// Writing unsused bytes after pixel array
	if err := binary.Write(file, binary.LittleEndian, unusedBuf2); err != nil {
		return err
	}

	return nil
}
