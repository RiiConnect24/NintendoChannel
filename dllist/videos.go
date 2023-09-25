package dllist

import (
	"NintendoChannel/constants"
	"unicode/utf16"
	"database/sql"
	"fmt"
	"strings"
)

type VideoTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	VideoType   uint8
	Unknown     [14]byte
	Unknown2    uint8
	RatingID    uint8
	Unknown3    uint8
	IsNew       uint8
	// Starts at 1 and incremented by 1
	VideoIndex uint8
	Unknown4   [2]byte
	Title      [123]uint16
}

type NewVideoTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	Unknown     [15]byte
	Unknown2    uint8
	RatingID    uint8
	Unknown3    uint8
	Title       [102]uint16
}

type PopularVideosTable struct {
	ID          uint32
	VideoLength uint16
	TitleID     uint32
	BarColor    uint8
	_           [15]byte
	RatingID    uint8
	Unknown     uint8
	VideoRank   uint8
	Unknown2    uint8
	Title       [102]uint16
}

func (l *List) MakeVideoTable() {
	l.Header.VideoTableOffset = l.GetCurrentSize()
	
	pool, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "DellServerzz", "127.0.0.1", 3306, "nc"))
	if err != nil {
		panic(err)
	}

	var title [123]uint16
	tempTitle := utf16.Encode([]rune("Go to \"New Arrivals\" >\n\"New Videos\" to watch\nany video."))
	copy(title[:], tempTitle)

	l.VideoTable = append(l.VideoTable, VideoTable{
		ID:          uint32(2130871958),
		VideoLength: uint16(280),
		TitleID:     0,
		VideoType:   uint8(7),
		Unknown:     [14]byte{},
		Unknown2:    0,
		RatingID:    9,
		Unknown3:    1,
		IsNew:       0,
		VideoIndex:  uint8(0),
		Unknown4:    [2]byte{222, 222},
		Title:       title,
	})

	rows, err := pool.Query(constants.GetPopularVideoQueryString(l.language))
	checkError(err)

	index := 1
	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&id, &queriedTitle, &length, &videoType)
		checkError(err)

		var title [123]uint16
		tempTitle := utf16.Encode([]rune(strings.Replace(queriedTitle, "\\n", "\n", -1)))
		copy(title[:], tempTitle)

		l.VideoTable = append(l.VideoTable, VideoTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			VideoType:   uint8(videoType),
			Unknown:     [14]byte{},
			Unknown2:    0,
			RatingID:    9,
			Unknown3:    1,
			IsNew:       0,
			VideoIndex:  uint8(index),
			Unknown4:    [2]byte{222, 222},
			Title:       title,
		})
		index++
		if (index == 60) {
			break
		}
	}

	l.Header.NumberOfVideoTables = uint32(len(l.VideoTable))
}

func (l *List) MakeNewVideoTable() {
	l.Header.NewVideoTableOffset = l.GetCurrentSize()
	
	pool, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "DellServerzz", "127.0.0.1", 3306, "nc"))
	if err != nil {
		panic(err)
	}

	rows, err := pool.Query(constants.GetVideoQueryString(l.language))
	checkError(err)

	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&id, &queriedTitle, &length, &videoType)
		checkError(err)

		var title [102]uint16
		tempTitle := utf16.Encode([]rune(strings.Replace(queriedTitle, "\\n", "\n", -1)))
		copy(title[:], tempTitle)

		l.NewVideoTable = append(l.NewVideoTable, NewVideoTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			Unknown:     [15]byte{8, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Unknown2:    0,
			RatingID:    9,
			Unknown3:    1,
			Title:       title,
		})
	}

	l.Header.NumberOfNewVideoTables = uint32(len(l.NewVideoTable))
}

func (l *List) MakePopularVideoTable() {
	l.Header.PopularVideoTableOffset = l.GetCurrentSize()
	
	pool, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "DellServerzz", "127.0.0.1", 3306, "nc"))
	if err != nil {
		panic(err)
	}

	return

	rows, err := pool.Query(constants.GetPopularVideoQueryString(l.language))
	checkError(err)

	for rows.Next() {
		var id int
		var queriedTitle string
		var length int
		var videoType int

		err = rows.Scan(&id, &queriedTitle, &length, &videoType)
		checkError(err)

		var title [102]uint16
		tempTitle := utf16.Encode([]rune(strings.Replace(queriedTitle, "\\n", "\n", -1)))
		copy(title[:], tempTitle)

		l.PopularVideosTable = append(l.PopularVideosTable, PopularVideosTable{
			ID:          uint32(id),
			VideoLength: uint16(length),
			TitleID:     0,
			BarColor:    0,
			RatingID:    9,
			Unknown:     1,
			VideoRank:   1,
			Unknown2:    222,
			Title:       title,
		})
	}

	l.Header.NumberOfPopularVideoTables = uint32(len(l.PopularVideosTable))
}
