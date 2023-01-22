package shortener

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenUrl(t *testing.T) {
	invoked := false
	baseUrl := "https://example.com"
	mockedShortLinkExists := func(string) (bool, error) {
		invoked = true
		return false, nil
	}
	linkShortener := NewLinkShortener(baseUrl, mockedShortLinkExists)
	shortLink, err := linkShortener.Shorten("https://example.com/full-link/full-link/full-link")
	assert.Nil(t, err)
	assert.Equal(t, 25, len(shortLink))
	assert.True(t, invoked)
	assert.True(t, strings.HasPrefix(shortLink, baseUrl))
}

func TestShortenUrl_WithCollision(t *testing.T) {
	invokedTimes := 0
	baseUrl := "https://example.com"
	firstLink := ""
	mockedShortLinkExists := func(link string) (bool, error) {
		if invokedTimes == 0 {
			firstLink = link
		}
		invokedTimes += 1
		return (invokedTimes == 1), nil
	}
	linkShortener := NewLinkShortener(baseUrl, mockedShortLinkExists)
	shortLink, err := linkShortener.Shorten("https://example.com/full-link/full-link/full-link")
	assert.Nil(t, err)
	assert.True(t, len(firstLink) > 0)
	assert.NotEqual(t, firstLink, shortLink)
	assert.Equal(t, 25, len(shortLink))
	assert.Equal(t, 2, invokedTimes)
	assert.True(t, strings.HasPrefix(shortLink, baseUrl))
}
