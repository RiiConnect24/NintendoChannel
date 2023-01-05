package dllist

import (
	"NintendoChannel/constants"
	"NintendoChannel/gametdb"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"io"
	"log"
	"os"
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
	imageBuffer *bytes.Buffer
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Nintendo Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

var pool *pgxpool.Pool
var ctx = context.Background()

type WorkerCtx struct {
	region   constants.RegionMeta
	language constants.Language
}

func Worker(worker <-chan WorkerCtx, b chan<- int) {
	for w := range worker {
		fmt.Printf("Starting worker - Region: %d, Language: %d\n", w.region.Region, w.language)
		list := List{
			region:      w.region.Region,
			ratingGroup: w.region.RatingGroup,
			language:    w.language,
			imageBuffer: new(bytes.Buffer),
		}

		list.MakeHeader()
		list.MakeRatingsTable()
		list.MakeTitleTypeTable()
		list.MakeCompaniesTable()
		list.MakeTitleTable()
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

		err = os.WriteFile(fmt.Sprintf("lists/dllist_%d_%d.bin", w.region.Region, w.language), compressed, 0666)
		checkError(err)
		fmt.Printf("Finished worker - Region: %d, Language: %d\n", w.region.Region, w.language)
		b <- 1
	}
}

func SpawnWorker(number int, contexts chan WorkerCtx, results chan int) {
	for i := 0; i < number; i++ {
		go Worker(contexts, results)
	}
}

func ResolveWorkers(number int, contexts chan WorkerCtx, results chan int) {
	for i := 0; i < number; i++ {
		<-results
		SpawnWorker(1, contexts, results)
	}
}

func MakeDownloadList() {
	// Initialize database
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", "noahpistilli", "2006", "127.0.0.1", "nc")
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	// Ensure this Postgresql connection is valid.
	defer pool.Close()
	gametdb.PrepareGameTDB()

	contexts := make(chan WorkerCtx, 10)
	results := make(chan int, 10)

	SpawnWorker(3, contexts, results)

	for _, region := range constants.Regions {
		for _, language := range region.Languages {
			contexts <- WorkerCtx{region: region, language: language}
		}
	}

	ResolveWorkers(10, contexts, results)
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
