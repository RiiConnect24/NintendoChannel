package dllist

import (
	"NintendoChannel/constants"
	"bufio"
	"database/sql"
	"fmt"
	"os"
)

type RecentRecommendationTable struct {
	TitleOffset uint32
	Medal       constants.Medal
	Unknown     uint8
}

const QueryRecommendations = `SELECT COUNT(game_id), game_id FROM recommendations GROUP BY game_id`

func (l *List) QueryRecommendations() {
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

	pool, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "rc24", password, "127.0.0.1", 3306, "rc24_nc"))
	if err != nil {
		panic(err)
	}

	rows, err := pool.Query(QueryRecommendations)
	checkError(err)

	for rows.Next() {
		var gameID string
		var count int
		err = rows.Scan(&count, &gameID)
		checkError(err)

		// First see if this game could exist in all regions
		isForRegion := false
		if gameID[3:] == "A" || gameID[3:] == "B" || gameID[3:] == "U" || gameID[3:] == "X" {
			isForRegion = true
		} else {
			// Now determine if the game exists for this region
			switch l.region {
			case constants.NTSC:
				if gameID[3:] == "E" || gameID[3:] == "N" {
					isForRegion = true
				}
				break
			case constants.Japan:
				if gameID[3:] == "J" {
					isForRegion = true
				}
				break
			case constants.PAL:
				if gameID[3:] == "P" || gameID[3:] == "L" || gameID[3:] == "M" {
					isForRegion = true
				}
				break
			}
		}

		if !isForRegion {
			continue
		}

		l.recommendations[gameID] = count
	}
}

func (l *List) MakeRecommendationTable() {
	l.Header.RecommendationTableOffset = l.GetCurrentSize()

	for gameID := range l.recommendations {
		// Now we find the title from our title table
		for i, title := range l.TitleTable {
			if string(title.TitleID[:]) == gameID {
				l.RecommendationTable = append(l.RecommendationTable, (236*uint32(i))+l.Header.TitleTableOffset)
				break
			}
		}
	}

	l.Header.NumberOfRecommendationTables = uint32(len(l.RecommendationTable))
}

func (l *List) MakeRecentRecommendationTable() {
	l.Header.RecentRecommendationTableOffset = l.GetCurrentSize()

	for gameID, num := range l.recommendations {
		for i, title := range l.TitleTable {
			if string(title.TitleID[:]) == gameID {
				l.RecentRecommendationTable = append(l.RecentRecommendationTable, RecentRecommendationTable{
					TitleOffset: (236 * uint32(i)) + l.Header.TitleTableOffset,
					Medal:       GetMedal(num),
					Unknown:     222,
				})
				break
			}
		}
	}

	l.Header.NumberOfRecentRecommendationTables = uint32(len(l.RecentRecommendationTable))
}
