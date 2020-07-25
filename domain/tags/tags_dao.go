package tags

import (
	"errors"
	"log"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/datasources/postres/conn"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/logger"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

const (
	errorNoRows = "no rows in result set"

	queryCreateTagInd = `INSERT INTO tagscollection (username, targetpost, tagcontent) VALUES ($1, $2, $3) ON CONFLICT (username, targetpost, tagcontent) DO NOTHING returning id`
	queryUpsertTag    = `INSERT INTO tags (postid, tagcontent) VALUES ($1, $2) ON CONFLICT (postid, tagcontent) DO NOTHING returning id`
	queryReadTagsPost = `SELECT id, targetpost, tagcontent FROM tagscollection WHERE targetpost=$1`
)

//Read

func (t *TagInd) Create() rest_errors.RestErr {
	stmt, err := conn.DB.Prepare(queryCreateTagInd)
	if err != nil {
		logger.Error("error when trying to prepare save item statement", err)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}
	defer stmt.Close()

	saveErr := stmt.QueryRow(t.Username, t.TargetPost, t.TagContent).Scan(&t.ID)
	if saveErr != nil {
		logger.Error("error saving item", saveErr)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}

	// stmtTwo, err := conn.DB.Prepare(queryUpsertTag)
	// if err != nil {
	// 	logger.Error("error when trying to prepare save item statement", err)
	// 	return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	// }
	// defer stmtTwo.Close()

	// var ID int64
	// saveErrTwo := stmtTwo.QueryRow(t.TargetPost, t.TagContent).Scan(&ID)
	// if err != nil {
	// 	logger.Error("error saving item", saveErrTwo)
	// 	return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	// }
	// log.Println(ID)

	return nil
}

func (t *Tag) ReadTagsPost(ID int64, username string) ([]Tag, rest_errors.RestErr) {
	stmt, err := conn.DB.Prepare(queryReadTagsPost)
	if err != nil {
		logger.Error("error preparing studio statement for item", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer stmt.Close()

	log.Println(username)

	rows, err := stmt.Query(ID)
	if err != nil {
		logger.Error("error selecting studio items", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer rows.Close()

	res := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		// username := ""
		if err := rows.Scan(&tag.ID, &tag.PostID, &tag.TagContent); err != nil {
			logger.Error("error scanning item row into struct", err)
			return nil, rest_errors.NewInternalServerError("Error parsing DB response", errors.New("DB error"))
		}

		// if username == "" {
		// 	tag.HasVoted = false
		// } else {
		// 	tag.HasVoted = true
		// }

		res = append(res, tag)
	}

	if len(res) == 0 {
		logger.Info("No results from search")
		return nil, rest_errors.NewNotFoundError("No results from search")
	}

	return res, nil
}
