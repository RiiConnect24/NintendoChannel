package info

import (
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"hash/crc32"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf16"
)

type Info struct {
	Header               Header
	SupportedControllers SupportedControllers
	SupportedFeatures    SupportedFeatures
	SupportedLanguages   SupportedLanguages
	_                    [10]byte
	Title                [31]uint16
	Subtitle             [31]uint16
	ShortTitle           [31]uint16
	DescriptionText      [3][41]uint16
	GenreText            [29]uint16
	PlayersText          [41]uint16
	PeripheralsText      [44]uint16
	_                    [80]byte
	DisclaimerText       [2400]uint16
	RatingID             uint8
	DistributionDateText [41]uint16
	WiiPointsText        [41]uint16
	CustomText           [10][41]uint16
	// Writing a blank one if it doesn't exist and not point to it
	// is way more efficient than writing everything individually
	TimePlayed TimePlayed
}

var timePlayed = map[string]TimePlayed{}

func (i *Info) MakeInfo(fileID uint32, game *gametdb.Game, title, synopsis string, region constants.Region, language constants.Language, titleType constants.TitleType, ratingDescriptors [7]string) {
	// Make other fields
	i.GetSupportedControllers(&game.Controllers)
	i.GetSupportedFeatures(&game.Features)
	i.GetSupportedLanguages(game.Languages)

	// Make title clean
	if strings.Contains(title, ": ") {
		splitTitle := strings.Split(title, ": ")
		copy(i.Title[:], utf16.Encode([]rune(splitTitle[0])))
		copy(i.Subtitle[:], utf16.Encode([]rune(splitTitle[1])))
	} else if strings.Contains(title, " - ") {
		splitTitle := strings.Split(title, " - ")
		copy(i.Title[:], utf16.Encode([]rune(splitTitle[0])))
		copy(i.Subtitle[:], utf16.Encode([]rune(splitTitle[1])))
	} else if len(title) > 31 {
		wrappedTitle := wordwrap.WrapString(title, 31)
		for p, s := range strings.Split(wrappedTitle, "\n") {
			switch p {
			case 0:
				copy(i.Title[:], utf16.Encode([]rune(s)))
				break
			case 1:
				copy(i.Subtitle[:], utf16.Encode([]rune(s)))
				break
			default:
				break
			}
		}
	} else {
		copy(i.Title[:], utf16.Encode([]rune(title)))
	}

	// Make synopsis
	wrappedSynopsis := strings.Split(wordwrap.WrapString(strings.Replace(strings.Replace(synopsis, "\n", "", -1), "  ", " ", -1), 40), "\n")
	if len(wrappedSynopsis) <= 3 {
		for i2, s := range wrappedSynopsis {
			copy(i.DescriptionText[i2][:], utf16.Encode([]rune(s)))
		}
	} else {
		for i2, s := range wrappedSynopsis {
			if i2 == 10 {
				break
			} else if i2 == 9 {
				s = strings.Split(s, ".")[0] + "."
			}

			copy(i.CustomText[i2][:], utf16.Encode([]rune(s)))
		}
	}

	// Write the online players text if any
	if game.Features.OnlinePlayers != 0 {
		temp := []uint16{0, 0}
		copy(i.PlayersText[:], append(temp, utf16.Encode([]rune(fmt.Sprintf("%d Players (Online)", game.Features.OnlinePlayers)))...))
	}

	var temp_ []uint16 // Declare a slice to store UTF-16 encoded values

	for _, s := range strings.Split(game.Genre, ",") {
		convertedString := s // Since s is already a string, no need to convert
		capitalized := capitalizeString(convertedString)

		// Append the utf16 encoded value of capitalized to temp_
		encoded := utf16.Encode([]rune(capitalized))
		temp_ = append(temp_, encoded...)

		// Append utf16 encoded value of ", " to temp_ if it's not the last entry
		commaSpace := utf16.Encode([]rune(", "))
		temp_ = append(temp_, commaSpace...)
	}

	temp_ = temp_[:len(temp_)-2]

	var genreText []uint16
	genreText = append(genreText, temp_...)

	copy(i.GenreText[:], genreText)

	copy(i.DisclaimerText[:], utf16.Encode([]rune("Game information is provided by GameTDB.")))

	if v, ok := timePlayed[game.ID[:4]]; ok {
		i.Header.TimesPlayedTableOffset = 6744
		i.TimePlayed = v
	}

	temp := new(bytes.Buffer)
	imageBuffer := new(bytes.Buffer)
	i.WriteAll(temp, imageBuffer)

	i.Header.PictureOffset = i.GetCurrentSize(imageBuffer)
	i.WriteCoverArt(imageBuffer, titleType, region, game.ID)
	i.WriteDetailedRatingImage(imageBuffer, region, ratingDescriptors, fileID)
	i.WriteRatingImage(imageBuffer, region)
	i.Header.Filesize = i.GetCurrentSize(imageBuffer)
	temp.Reset()

	i.WriteAll(temp, imageBuffer)
	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(temp.Bytes(), crcTable)
	i.Header.CRC32 = checksum
	temp.Reset()

	var reg [3]string
	reg[0] = "JP"
	reg[1] = "GB"
	reg[2] = "US"

	var lang [7]string
	lang[0] = "ja"
	lang[1] = "en"
	lang[2] = "de"
	lang[3] = "fr"
	lang[4] = "es"
	lang[5] = "it"
	lang[6] = "nl"

	i.WriteAll(temp, imageBuffer)

	err := os.MkdirAll(fmt.Sprintf("./infos/%d/%d/", region, language), 0755)
	checkError(err)
	err = os.WriteFile(fmt.Sprintf("./infos/%d/%d/%d.info", region, language, fileID), temp.Bytes(), 0666)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Nintendo Channel info file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func (i *Info) WriteAll(buffer, imageBuffer *bytes.Buffer) {
	err := binary.Write(buffer, binary.BigEndian, *i)
	checkError(err)

	buffer.Write(imageBuffer.Bytes())
}

func (i *Info) GetCurrentSize(_buffer *bytes.Buffer) uint32 {
	buffer := bytes.NewBuffer(nil)
	i.WriteAll(buffer, _buffer)
	return uint32(buffer.Len())
}

func capitalizeString(input string) string {
	words := strings.Fields(input) // Split the input into words
	var capitalizedWords []string

	for _, word := range words {
		if len(word) > 0 {
			// Check if the word is "and" or "of"
			if word == "of" || word == "and" {
				// Add the word as is (not capitalized)
				capitalizedWords = append(capitalizedWords, word)
			} else {
				// Capitalize the first letter of the word
				capitalizedWord := string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:])
				capitalizedWords = append(capitalizedWords, capitalizedWord)
			}
		}
	}

	return strings.Join(capitalizedWords, " ")
}

func GetTimePlayed(ctx context.Context, pool *sql.DB) {
	rows, err := pool.Query(`SELECT game_id, COUNT(game_id), SUM(times_played), SUM(time_played) FROM time_played GROUP BY game_id`)
	checkError(err)

	for rows.Next() {
		var gameID string
		var numberOfPlayers int
		var totalTimesPlayed int
		var totalTimePlayed int

		err = rows.Scan(&gameID, &numberOfPlayers, &totalTimesPlayed, &totalTimePlayed)
		checkError(err)

		timePlayed[gameID] = TimePlayed{
			TotalTimePlayed:           uint32(totalTimePlayed / 60),
			TimeSpentPlayingPerPerson: uint32(totalTimePlayed / numberOfPlayers),
			TotalTimesPlayed:          uint32(totalTimesPlayed),
			TimesPlayedPerPerson:      uint32((float64(totalTimesPlayed / numberOfPlayers)) / 0.01),
		}
	}
}
