package info

import (
	"NintendoChannel/constants"
	"bytes"
	_ "embed"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"github.com/disintegration/imaging"
)

var regionToStr = map[constants.Region]string{
	constants.Japan: "JA",
	constants.PAL:   "EN",
	constants.NTSC:  "US",
}
var titleTypeToStr = map[constants.TitleType]string{
	constants.Wii:             "wii",
	constants.NintendoDS:      "ds",
	constants.NintendoThreeDS: "3ds",
}

var consoleToImageType = map[constants.TitleType]string{
	constants.Wii:             "cover",
	constants.NintendoDS:      "box",
	constants.NintendoThreeDS: "box",
}

var consoleToTempImageType = map[constants.TitleType][]byte{
	constants.Wii:             Placeholder3DS,
	constants.NintendoDS:      PlaceholderDS,
	constants.NintendoThreeDS: Placeholder3DS,
}

//go:embed 3ds.jpg
var Placeholder3DS []byte

//go:embed ds.jpg
var PlaceholderDS []byte

//go:embed wii.jpg
var PlaceholderWii []byte

func (i *Info) WriteCoverArt(buffer *bytes.Buffer, titleType constants.TitleType, region constants.Region, gameID string) {
	url := fmt.Sprintf("https://art.gametdb.com/%s/%s/%s/%s.png", titleTypeToStr[titleType], consoleToImageType[titleType], regionToStr[region], gameID)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		buffer.Write(consoleToTempImageType[titleType])
	} else {
		coverImg, err := png.Decode(resp.Body)
		checkError(err)

		coverImgBounds := coverImg.Bounds()
		coverImgWidth, coverImgHeight := coverImgBounds.Dx(), coverImgBounds.Dy()

		var coverImgResized image.Image
		if titleTypeToStr[titleType] != "3DS" && titleTypeToStr[titleType] != "NDS" {
			coverImgResized = resizeImage(coverImg, int(float64(coverImgWidth)*(384.0/float64(coverImgHeight))), 384)
		} else {
			coverImgResized = resizeImage(coverImg, 384, int(float64(coverImgHeight)*(384.0/float64(coverImgWidth))))
		}

		coverImgResizedBounds := coverImgResized.Bounds()
		coverImgResizedWidth, coverImgResizedHeight := coverImgResizedBounds.Dx(), coverImgResizedBounds.Dy()

		offsetX := (384 - coverImgResizedWidth) / 2
		offsetY := (384 - coverImgResizedHeight) / 2
		offset := image.Pt(offsetX, offsetY)

		newImage := image.NewRGBA(image.Rect(0, 0, 384, 384))
		drawImage(newImage, coverImgResized, offset)

		err = jpeg.Encode(buffer, newImage, nil)
		checkError(err)
	}

	i.Header.PictureSize = uint32(buffer.Len())
}

func resizeImage(img image.Image, width, height int) image.Image {
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	return resizedImg
}

func drawImage(dst draw.Image, src image.Image, offset image.Point) {
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)
	draw.Draw(dst, src.Bounds().Add(offset), src, image.Point{}, draw.Src)
}

func (i *Info) WriteRatingImage(buffer *bytes.Buffer, region constants.Region) {
	i.Header.RatingPictureOffset = i.GetCurrentSize(buffer)

	regionToRatingGroup := map[constants.Region]constants.RatingGroup{
		constants.Japan: constants.CERO,
		constants.NTSC:  constants.ESRB,
		constants.PAL:   constants.PEGI,
	}

	buffer.Write(constants.Images[regionToRatingGroup[region]][i.RatingID-8])
	i.Header.RatingPictureSize = uint32(len(constants.Images[regionToRatingGroup[region]][i.RatingID-8]))
}
