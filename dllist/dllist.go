package dllist

import (
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"NintendoChannel/info"
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type List struct {
	Header          Header
	RatingsTable    []RatingTable
	TitleTypesTable []TitleTypeTable
	CompaniesTable  []CompanyTable
	TitleTable      []TitleTable
	// NewTitleTable is an array of pointers to titles in TitleTable
	NewTitleTable             []uint32
	VideoTable                []VideoTable
	NewVideoTable             []NewVideoTable
	DemoTable                 []DemoTable
	RecommendationTable       []uint32
	RecentRecommendationTable []RecentRecommendationTable
	PopularVideosTable        []PopularVideosTable
	DetailedRatingTable       []DetailedRatingTable

	// Below are variables that help us keep state
	region      constants.Region
	ratingGroup constants.RatingGroup
	language    constants.Language
	// map[game_id]amount_voted
	recommendations map[string]int
	imageBuffer     *bytes.Buffer
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Nintendo Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

var pool *sql.DB
var ctx = context.Background()

func MakeDownloadList(overwrite bool) {
	file, err := os.Open("sql.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the password from the file
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	password := scanner.Text()

	// Check for errors while scanning
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Initialize database
	pool, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "rc24", password, "127.0.0.1", 3306, "rc24_nc"))
	if err != nil {
		panic(err)
	}

	// Ensure this Postgresql connection is valid.
	defer pool.Close()
	gametdb.PrepareGameTDB()
	info.GetTimePlayed(ctx, pool)

	wg := sync.WaitGroup{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	semaphore := make(chan struct{}, 3)

	wg.Add(10)
	for _, region := range constants.Regions {
		for _, language := range region.Languages {
			go func(_region constants.RegionMeta, _language constants.Language) {
				defer wg.Done()
				semaphore <- struct{}{}

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

				fmt.Printf("Starting worker - Region: %s, Language: %s\n", reg[_region.Region], lang[_language])
				list := List{
					region:          _region.Region,
					ratingGroup:     _region.RatingGroup,
					language:        _language,
					imageBuffer:     new(bytes.Buffer),
					recommendations: map[string]int{},
				}

				list.QueryRecommendations()

				list.MakeHeader()
				list.MakeRatingsTable()
				list.MakeTitleTypeTable()
				list.MakeCompaniesTable()
				list.MakeTitleTable(overwrite)
				list.MakeNewTitleTable()
				list.MakeVideoTable()
				list.MakeNewVideoTable()
				list.MakeDemoTable()
				list.MakeRecommendationTable()
				list.MakeRecentRecommendationTable()
				list.MakePopularVideoTable()
				list.MakeDetailedRatingTable()
				list.WriteRatingImages()

				temp := bytes.NewBuffer(nil)
				list.WriteAll(temp)
				list.Header.Filesize = uint32(temp.Len())
				temp.Reset()
				list.WriteAll(temp)

				crcTable := crc32.MakeTable(crc32.IEEE)
				checksum := crc32.Checksum(temp.Bytes(), crcTable)
				list.Header.CRC32 = checksum

				temp.Reset()
				list.WriteAll(temp)

				// Compress then write
				compressed, err := lz10.Compress(temp.Bytes())
				checkError(err)

				err = os.MkdirAll(fmt.Sprintf("lists/%d/%d/", _region.Region, _language), 0755)
				checkError(err)
				err = os.WriteFile(fmt.Sprintf("lists/%d/%d/dllist.bin", _region.Region, _language), compressed, 0666)
				checkError(err)
				fmt.Printf("Finished worker - Region: %s, Language: %s\n", reg[_region.Region], lang[_language])
				<-semaphore
			}(region, language)
		}
	}

	wg.Wait()
}

// Write writes the current values in Votes to an io.Writer method.
// This is required as Go cannot write structs with non-fixed slice sizes,
// but can write them individually.
func (l *List) Write(writer io.Writer, data any) {
	err := binary.Write(writer, binary.BigEndian, data)
	checkError(err)
}

func (l *List) WriteAll(writer io.Writer) {
	l.Write(writer, l.Header)
	l.Write(writer, l.RatingsTable)
	l.Write(writer, l.TitleTypesTable)
	l.Write(writer, l.CompaniesTable)
	l.Write(writer, l.TitleTable)
	l.Write(writer, l.NewTitleTable)
	l.Write(writer, l.VideoTable)
	l.Write(writer, l.NewVideoTable)
	l.Write(writer, l.DemoTable)
	l.Write(writer, l.RecommendationTable)
	l.Write(writer, l.RecentRecommendationTable)
	l.Write(writer, l.PopularVideosTable)
	l.Write(writer, l.DetailedRatingTable)
}

// GetCurrentSize returns the current size of our List struct.
// This is useful for calculating the current offset of List.
func (l *List) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer(nil)
	l.WriteAll(buffer)
	buffer.Write(l.imageBuffer.Bytes())

	return uint32(buffer.Len())
}
