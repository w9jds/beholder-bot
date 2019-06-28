package main

import (
	"database/sql"
	"fmt"
	"log"
)

func storeNewGame(guildID, categoryID, textID, voiceID, dmID string) bool {
	statement := `INSERT INTO games (guild_id, category_id, text_id, voice_id, dm_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := postgres.Exec(statement, guildID, categoryID, textID, voiceID, dmID)
	if err != nil {
		log.Printf("Error storing new game into postgres guild: %s category: %s", guildID, categoryID)
		log.Println(err)
		return false
	}

	return true
}

func storeNewMap(guildID, channelID, name, messageID string) bool {
	statement := `INSERT INTO maps (guild_id, text_id, message_id, name) VALUES ($1, $2, $3, $4)`
	_, err := postgres.Exec(statement, guildID, channelID, messageID, name)
	if err != nil {
		log.Printf("Error storing new map into postgres name: %s category: %s", name, channelID)
		log.Println(err)
		return false
	}

	return true
}

func getStoredMap(guildID, channelID, name string) (string, error) {
	var messageID string
	statement := `SELECT message_id FROM maps WHERE guild_id='%s' and text_id='%s' and name='%s'`

	row := postgres.QueryRow(fmt.Sprintf(statement, guildID, channelID, name))
	switch error := row.Scan(&messageID); error {
	case sql.ErrNoRows:
		return "", sql.ErrNoRows
	case nil:
		return messageID, nil
	default:
		return "", error
	}
}

// func storeNextSession(categoryID string, nextSession time) {
// 	statement := `INSERT INTO sessions (guild_id, category_id, date) VALUES ($1, $2, $3, $4)`
// 	_, err := postgres.Exec(statement, guildID, categoryID, nextSession)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func isGameChannel(guildID, textID string) bool {
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

func isGameDM(guildID, textID, authorID string) bool {
	_, dmID, err := getGameFromTextChannel(guildID, textID)
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

func getGameFromTextChannel(guildID, textID string) (string, string, error) {
	var dmID, categoryID string
	statement := `SELECT dm_id, category_id from games where guild_id='%s' and text_id='%s'`

	row := postgres.QueryRow(fmt.Sprintf(statement, guildID, textID))
	switch error := row.Scan(&dmID, &categoryID); error {
	case sql.ErrNoRows:
		return "", "", sql.ErrNoRows
	case nil:
		return categoryID, dmID, nil
	default:
		return "", "", error
	}
}

func removeChannel(guildID, channelID string) {
	var categoryID, textID string
	selectGame := `SELECT category_id, text_id WHERE guild_id='%[1]' and (text_id='%[2]' or category_id='%[2]')`
	gameStatement := `DELETE FROM games WHERE guild_id='%[1]' and (text_id='%[2]' or category_id='%[2]')`
	mapsStatement := `DELETE FROM maps WHERE guild_id='%[1]' and text_id='%[2]'`

	row := postgres.QueryRow(fmt.Sprintf(selectGame, guildID, channelID))
	switch error := row.Scan(&categoryID, &textID); error {
	case sql.ErrNoRows:
		return
	default:
		log.Println(error)
		return
	}

	_, err := postgres.Exec(fmt.Sprintf(gameStatement, guildID, channelID))
	if err != nil {
		log.Println(err)
	}

	_, err = postgres.Exec(fmt.Sprintf(mapsStatement, guildID, textID))
	if err != nil {
		log.Println(err)
	}
}

func removeMaps(guildID, channelID, messageID string) {
	mapsStatement := `DELETE FROM maps WHERE guild_id='%s' and text_id='%s' and message_id='%s'`

	_, err := postgres.Exec(fmt.Sprintf(mapsStatement, guildID, channelID, messageID))
	if err != nil {
		log.Println(err)
	}
}
