package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConvertParams(t *testing.T) {
	assert := assert.New(t)
	req, _ := http.NewRequest("POST", `/files?converts={"pic":"120x90"}`, nil)

	convert, err := GetConvertParams(req)

	assert.Nil(err)
	assert.Equal("120x90", convert["pic"])
}
