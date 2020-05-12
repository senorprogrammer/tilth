package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTagMap(t *testing.T) {
	tests := []struct {
		name        string
		pages       []*Page
		expectedLen int
	}{
		{
			name:        "with no pages",
			pages:       []*Page{},
			expectedLen: 0,
		},
		{
			name: "with pages",
			pages: []*Page{
				{TagsStr: "go, ada"},
			},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewTagMap(tt.pages).Tags

			assert.Equal(t, tt.expectedLen, len(actual))
		})
	}
}

func Test_TagMap_Add(t *testing.T) {
	tests := []struct {
		name        string
		tag         *Tag
		expectedLen int
	}{
		{
			name:        "with an invalid tag",
			tag:         &Tag{},
			expectedLen: 0,
		},
		{
			name:        "with a new tag",
			tag:         &Tag{Name: "go"},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tMap := NewTagMap([]*Page{})
			tMap.Add(tt.tag)

			actual := tMap.Tags

			assert.Equal(t, tt.expectedLen, len(actual))
		})
	}
}

func Test_TagMap_BuildFromPages(t *testing.T) {
	tests := []struct {
		name        string
		pages       []*Page
		expectedLen int
	}{
		{
			name:        "with no pages",
			pages:       []*Page{},
			expectedLen: 0,
		},
		{
			name: "with pages",
			pages: []*Page{
				{TagsStr: "go"},
				{TagsStr: "ada"},
			},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tMap := NewTagMap([]*Page{})
			tMap.BuildFromPages(tt.pages)

			actual := tMap.Tags

			assert.Equal(t, tt.expectedLen, len(actual))
		})
	}
}

func Test_TagMap_Get(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedLen int
	}{
		{
			name:        "with missing tag",
			input:       "ada",
			expectedLen: 0,
		},
		{
			name:        "with valid tag",
			input:       "go",
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		pages := []*Page{&Page{TagsStr: "go"}}
		tMap := NewTagMap(pages)

		actual := tMap.Get(tt.input)

		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedLen, len(actual))
		})
	}
}

func Test_TagMap_Len(t *testing.T) {
	tests := []struct {
		name        string
		page        *Page
		expectedLen int
	}{
		{
			name:        "with missing tag",
			page:        &Page{},
			expectedLen: 0,
		},
		{
			name:        "with valid tag",
			page:        &Page{TagsStr: "go"},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		pages := []*Page{tt.page}
		tMap := NewTagMap(pages)

		actual := tMap.Len()

		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedLen, actual)
		})
	}
}

func Test_TagMap_SortedTagNames(t *testing.T) {
	pages := []*Page{&Page{TagsStr: "go, ada, lua"}}
	tMap := NewTagMap(pages)

	expected := []string{"ada", "go", "lua"}
	actual := tMap.SortedTagNames()

	assert.Equal(t, expected, actual)
}