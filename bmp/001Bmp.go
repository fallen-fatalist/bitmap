package bmp

// Device Independent Bitmap
// for more detail see: https://en.wikipedia.org/wiki/BMP_file_format
type bmp struct {
	fileHeader *fileHeader
	dibHeader  *dibHeader
	pixelArray [][]byte
}

// Device independent bitmap header
// header corresponds to BITMAPINFOHEADER version for 24 bits
// see: https://upload.wikimedia.org/wikipedia/commons/7/75/BMPfileFormat.svg
type dibHeader struct {
	Size                  uint32
	Width                 uint32
	Height                uint32
	ColorPlane            uint16
	BitsPerPixel          uint16
	CompressionMethod     uint32
	ImageSize             uint32
	HorizontalResolution  uint32
	VerticalResolution    uint32
	ColorsNumber          uint32
	ImportantColorsNumber uint32
}

// Fields will be added in future to correspond BITMAPV5HEADER
// redChannel            uint32
// greenChannel          uint32
// blueChannel           uint32
// alphaChannel          uint32
// colorSpaceType        uint32
// colorSpaceEndpoints   uint32
// gammaRedChannel       uint32
// gammaGreenChannel     uint32
// gammaBlueChannel      uint32
// intent                uint32
// iccProfileData        uint32
// iccProfileSize        uint32
// reserved              uint32

// Compression methods
// see (https://en.wikipedia.org/wiki/BMP_file_format#DIB_header_(bitmap_information_header))
// Will be added in future
// const (
// 	BI_RGB = iota
// 	BI_RLE8
// 	BI_RLE4
// 	BI_BITFIELDS
// 	BI_JPEG
// 	BI_PNG
// 	BI_ALPHABITFIELDS
// 	BI_CMYK     = 11
// 	BI_CMYKRLE8 = 12
// 	BI_CMYKRLE4 = 13
// )

// Color table for color pallete <= 8 bits (will be added in future)
// type colorTable []byte

func (b *bmp) GetPixelNumber() uint16 {
	return b.dibHeader.BitsPerPixel
}
