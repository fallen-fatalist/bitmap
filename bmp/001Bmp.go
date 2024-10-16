package bmp

// Device Independent Bitmap
// for more detail see: https://en.wikipedia.org/wiki/BMP_file_format
type bmp struct {
	fileHeader *fileHeader
	dibHeader  *dibHeader
	colorTable colorTable
	pixelArray [][]byte
}

// Device independent bitmap header
// this header correspond to BITMAPV5HEADER
// see: https://upload.wikimedia.org/wikipedia/commons/7/75/BMPfileFormat.svg
type dibHeader struct {
	size                  uint32
	width                 uint32
	height                uint32
	colorPlane            uint16
	bitsPerPixel          uint16
	compressionMethod     uint32
	imageSize             uint32
	horizontalResolution  uint32
	verticalResolution    uint32
	colorsNumber          uint32
	importantColorsNumber uint32
	redChannel            uint32
	greenChannel          uint32
	blueChannel           uint32
	alphaChannel          uint32
	colorSpaceType        uint32
	colorSpaceEndpoints   uint32
	gammaRedChannel       uint32
	gammaGreenChannel     uint32
	gammaBlueChannel      uint32
	intent                uint32
	iccProfileData        uint32
	iccProfileSize        uint32
	reserved              uint32
}

type colorTable []byte
