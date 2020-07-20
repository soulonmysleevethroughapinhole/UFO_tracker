package comments

import (
	"errors"
	"fmt"
	"strings"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/datasources/postres/conn"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/logger"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

const (
	errorNoRows = "no rows in result set"

	queryCreateComment         = `INSERT INTO comments (username, threadid, parentid, content) VALUES ($1, $2, $3, $4) RETURNING id`
	queryReadComment           = `SELECT username, threadid, parentid, content, postdate FROM comments WHERE id = $1`
	queryReadAllCommentsThread = `SELECT id, username, threadid, parentid, content, postdate `
	queryUpdateComment         = ``
	queryDeleteComment         = ``
)

func (c *Comment) Read(commentID int64) rest_errors.RestErr {
	stmt, err := conn.DB.Prepare(queryReadComment)
	if err != nil {
		logger.Error("Error reading comment", err)
		return rest_errors.NewInternalServerError(fmt.Sprintf("Error reading comment %d", c.ID), errors.New("DB Error"))
	}
	defer stmt.Close()

	res := stmt.QueryRow(commentID)
	getErr := res.Scan(&c.Username, &c.ThreadID, &c.ParentID, &c.Content, &c.PostDate)
	if getErr != nil {
		if strings.Contains(getErr.Error(), "404") {
			logger.Error("error getting item - not found", getErr)
			return rest_errors.NewNotFoundError(fmt.Sprintf("%d not found", c.ID))
		}
		logger.Error("error getting item", getErr)
		return rest_errors.NewInternalServerError(fmt.Sprintf("Error getting audio %d", c.ID), errors.New("DB error"))
	}
	return nil
}

func (c *Comment) Create() rest_errors.RestErr {
	stmt, err := conn.DB.Prepare(queryCreateComment)
	if err != nil {
		logger.Error("error when trying to prepare save item statement", err)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}
	defer stmt.Close()

	saveErr := stmt.QueryRow(c.Username, c.ThreadID, c.ParentID, c.Content).Scan(c.ID)
	if saveErr != nil {
		logger.Error("error saving item", saveErr)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}
	return nil
}

func (c *Comment) Update() rest_errors.RestErr {
	return nil
}

func (c *Comment) ReadAll(threadID int64) (Comments, rest_errors.RestErr) {
	//case case case this shit
	stmt, err := conn.DB.Prepare(queryReadAllCommentsThread)
	if err != nil {
		logger.Error("error preparing studio statement for item", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(threadID)
	if err != nil {
		logger.Error("error selecting studio items", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer rows.Close()

	res := make(Comments, 0)
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&c.ID, &c.Username, &c.ThreadID, &c.ParentID, &c.Content, &c.PostDate); err != nil {
			logger.Error("error scanning item row into struct", err)
			return nil, rest_errors.NewInternalServerError("Error parsing DB response", errors.New("DB error"))
		}

		res = append(res, comment)
	}

	if len(res) == 0 {
		logger.Info("No results from search")
		return nil, rest_errors.NewNotFoundError("No results from search")
	}

	return res, nil
}
