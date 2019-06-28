package main

import (
	"database/sql"
	"log"
)

func storeNewGame(guildID string, categoryID string, textID string, voiceID string, dmID string) {
	statement := `INSERT INTO games (guild_id, category_id, text_id, voice_id, dm_id) VALUES ($1, $2, $3, $4)`
	_, err := postgres.Exec(statement, guildID, categoryID, textID, voiceID, dmID)
	if err != nil {
		log.Printf("Error storing new game into postgres guild: %s category: %s", guildID, categoryID)
		log.Fatal(err)
	}
}

// func storeNextSession(categoryID string, nextSession time) {
// 	statement := `INSERT INTO sessions (guild_id, category_id, date) VALUES ($1, $2, $3, $4)`
// 	_, err := postgres.Exec(statement, guildID, categoryID, nextSession)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func isGameChannel(guildID string, textID string) bool {
	_, _, err := getGameFromTextChannel(guildID, textID)
	switch err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		log.Fatal(err)
		return false
	}
}

func isGameDM(guildID string, textID string, authorID string) bool {
	dmID, _, err := getGameFromTextChannel(guildID, textID)
	switch err {
	case sql.ErrNoRows:
		return false
	case nil:
		return authorID == dmID
	default:
		log.Fatal(err)
		return false
	}
}

func getGameFromTextChannel(guildID string, textID string) (string, string, error) {
	var dmID, categoryID string
	statement := `SELECT dm_id, category_id from games where guild_id='$1' and text_id='$2'`

	row := postgres.QueryRow(statement, guildID, textID)
	switch error := row.Scan(&dmID, &categoryID); error {
	case sql.ErrNoRows:
		return "", "", sql.ErrNoRows
	case nil:
		return categoryID, dmID, nil
	default:
		return "", "", error
	}
}
