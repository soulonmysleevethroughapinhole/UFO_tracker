package tags

import (
	"errors"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/datasources/postres/conn"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/logger"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

const (
	errorNoRows = "no rows in result set"

	queryCreateTagInd = `INSERT INTO tagcollection (username, `
	queryUpsertTag    = ``
	queryReadTagsPost = ``
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

	stmtTwo, err := conn.DB.Prepare(queryUpsertTag)
	if err != nil {
		logger.Error("error when trying to prepare save item statement", err)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}
	defer stmtTwo.Close()

	saveErrTwo := stmtTwo.QueryRow(t.TargetPost, t.TagContent).Scan()
	if err != nil {
		logger.Error("error saving item", saveErrTwo)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}

	return nil
}

func (t *Tag) ReadTagsPost(ID int64, username string) ([]Tag, rest_errors.RestErr) {
	stmt, err := conn.DB.Prepare(queryReadTagsPost)
	if err != nil {
		logger.Error("error preparing studio statement for item", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(ID, username)
	if err != nil {
		logger.Error("error selecting studio items", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer rows.Close()

	res := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		username := ""
		if err := rows.Scan(&tag.ID, &tag.PostID, &tag.TagContent, &tag.VoteAmt, username); err != nil {
			logger.Error("error scanning item row into struct", err)
			return nil, rest_errors.NewInternalServerError("Error parsing DB response", errors.New("DB error"))
		}

		if username == "" {
			tag.HasVoted = false
		} else {
			tag.HasVoted = true
		}

		res = append(res, tag)
	}

	if len(res) == 0 {
		logger.Info("No results from search")
		return nil, rest_errors.NewNotFoundError("No results from search")
	}

	return res, nil
}
