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
        coverImg, _, err := image.Decode(resp.Body)
        checkError(err)

        // Check if the image is a PNG with a transparent background.
        _, isPNG := coverImg.(*image.NRGBA)
        if isPNG {
            // Handle transparent PNGs here.
            coverImgResized := resizeImageWithAspectRatio(coverImg, 384, 384)
            err = jpeg.Encode(buffer, coverImgResized, nil)
            checkError(err)
        } else {
            // For non-PNG images, create a new RGBA image with a white background.
            newImage := image.NewRGBA(image.Rect(0, 0, 384, 384))
            draw.Draw(newImage, newImage.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

            // Resize the image with the new dimensions.
            coverImgResized := resizeImageWithAspectRatio(coverImg, 384, 384)

            // Calculate the offset to center the resized image on the white background.
            offsetX := (384 - coverImgResized.Bounds().Dx()) / 2
            offsetY := (384 - coverImgResized.Bounds().Dy()) / 2
            offset := image.Pt(offsetX, offsetY)

            // Draw the resized image onto the newImage with transparency.
            draw.Draw(newImage, newImage.Bounds().Add(offset), coverImgResized, image.Point{}, draw.Over)

            err = jpeg.Encode(buffer, newImage, nil)
            checkError(err)
        }
    }

    i.Header.PictureSize = uint32(buffer.Len())
}

func resizeImageWithAspectRatio(img image.Image, width, height int) image.Image {
    imgBounds := img.Bounds()
    imgWidth, imgHeight := imgBounds.Dx(), imgBounds.Dy()
    
    // Calculate the aspect ratio of the original image.
    aspectRatio := float64(imgWidth) / float64(imgHeight)

    // Calculate the new dimensions while preserving the aspect ratio.
    newWidth := width
    newHeight := int(float64(newWidth) / aspectRatio)

    // Check if the new height exceeds the specified height.
    if newHeight > height {
        newHeight = height
        newWidth = int(float64(newHeight) * aspectRatio)
    }

    // Resize the image with the new dimensions.
    resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

    // Create a new image with the specified dimensions and draw the resized image onto it.
    resultImg := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(resultImg, resultImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

    offsetX := (width - newWidth) / 2
    offsetY := (height - newHeight) / 2
    offset := image.Pt(offsetX, offsetY)

    draw.Draw(resultImg, resultImg.Bounds().Add(offset), resizedImg, image.Point{}, draw.Over)

    return resultImg
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

	buffer.Write(constants.ImagesSmall[regionToRatingGroup[region]][i.RatingID-8])
	i.Header.RatingPictureSize = uint32(len(constants.ImagesSmall[regionToRatingGroup[region]][i.RatingID-8]))
}
