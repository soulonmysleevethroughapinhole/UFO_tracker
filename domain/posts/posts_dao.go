package posts

import (
	"errors"
	"fmt"	"strings"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/datasources/postres/conn"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/logger"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)


const (
	errorNoRows = "no rows in result set"

	queryCreatePost = ``
	queryReadPost   = ``
	queryReadAllPost   = ``
	queryUpdatePost = ``
	queryDeletePost = ``
)

func (p *Post) Read() rest_errors.RestErr {
	stmt, err := conn.DB.Prepare(queryReadPost)
	if err != nil {
		logger.Error("Error reading post", err)
		return rest_errors.NewInternalServerError(fmt.Sprintf("Error reading post %d", p.ID), errors.New("DB Error"))
	}
	defer stmt.Close()

	res := stmt.QueryRow(p.ID)
	getErr := res.Scan(&p.Username, &p.Title, &p.ContentType, &p.Content, &p.PostDate)
	if getErr != nil {
		if strings.Contains(getErr.Error(), "404") {
			logger.Error("error getting item - not found", getErr)
			return rest_errors.NewNotFoundError(fmt.Sprintf("%d not found", p.ID))
		}
		logger.Error("error getting item", getErr)
		return rest_errors.NewInternalServerError(fmt.Sprintf("Error getting audio %d", p.ID), errors.New("DB error"))
	}
	return nil
}

func (p *Post) Create() rest_errors.RestErr {
	stmt, err := conn.DB.Prepare(queryCreatePost)
	if err != nil {
		logger.Error("error when trying to prepare save item statement", err)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}
	defer stmt.Close()

	saveErr := stmt.QueryRow(p.Username, p.Title, p.ContentType, p.Content)
	if saveErr != nil {
		logger.Error("error saving item", saveErr)
		return rest_errors.NewInternalServerError("Error saving item", errors.New("DB error"))
	}
	return nil
}

func (p *Post) Update() rest_errors.RestErr {
	return nil	
}

func (p *Post) ReadAll() (Posts, rest_errors.RestErr) {
	stmt, err := conn.DB.Prepare(queryReadAllPost)
	if err != nil {
		logger.Error("error preparing studio statement for item", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer stmt.Close()

	rors, err := stmt.Query()
	if err != nil {
		logger.Error("error selecting studio items", err)
		return nil, rest_errors.NewInternalServerError("Error searching documents", errors.New("DB error"))
	}
	defer rows.Close()

	res := make(Posts, 0)
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &p.ContentType, &p.Content, &p.PostDate); err != nil {
			logger.Error("error scanning item row into struct", err)
			return nil, rest_errors.NewInternalServerError("Error parsing DB response", errors.New("DB error"))
		}
		res = append(res, post)
	}

	if len(res) == 0 {
		logger.Info("No results from search")
		return nil, rest_errors.NewNotFoundError("No results from search")
	}

	return res, nil
}

//func (p *Post) ReadRising() (Posts, rest_errors.RestErr) {}

//func (p *Post) ReadHot() (Posts, rest_errors.RestErr) {}

//func (p *Post) ReadTop() (Posts, rest_errors.RestErr) {}
