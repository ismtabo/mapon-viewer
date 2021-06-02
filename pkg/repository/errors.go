package repository

import (
	"fmt"
	"regexp"

	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var conflictRegexp = regexp.MustCompile(`dup key: { (?P<key>[a-zA-Z0-9_.-]+): "(?P<value>.*)" }`)

func mapMongoError(err error) error {
	if err == mongo.ErrNoDocuments {
		return errors.NewNotFoundError()
	}
	if mErr, ok := err.(mongo.CommandError); ok && mErr.Code == 11000 {
		return mapMongoConflictError(mErr)
	}
	return errors.NewInternalServerError(err)
}

func mapMongoConflictError(err mongo.CommandError) error {
	match := conflictRegexp.FindStringSubmatch(err.Message)
	if len(match) == 0 {
		return errors.NewConflictError(err.Message)
	}
	var conflictKey string
	var conflictValue string
	for i, name := range conflictRegexp.SubexpNames() {
		switch name {
		case "key":
			conflictKey = match[i]
		case "value":
			conflictValue = match[i]
		}
	}
	return errors.NewConflictError(fmt.Sprintf("$.%s: %s already exists", conflictKey, conflictValue))
}
