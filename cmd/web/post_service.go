package main

import (
	"errors"
	"strings"
)

func IsValidData(title, content string, categories []string) error {
	
	if len(title) < 5 || len(title) > 50 || len(strings.TrimSpace(title)) == 0 {
		return errors.New("the title must be between 5-50 chars")
	}

	if content == "" || len(strings.TrimSpace(content)) == 0 {
		return errors.New("the content can't be empty")
	}

	if len(categories) < 1 || len(categories) > 3 {
		return errors.New("choose one of the available categories (max 3)")
	}

	return nil
}

func IsValidComment(content string) error {

	if content == "" || len(strings.TrimSpace(content)) == 0 {
		return errors.New("the comment can't be empty")
	}

	if len(content) < 5 {
		return errors.New("the content must be at least 5 chars")
	}

	return nil
}
