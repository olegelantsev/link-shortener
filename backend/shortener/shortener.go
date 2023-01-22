package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
)

type LinkShortener struct {
	shortLinkExists func(string) (bool, error)
	baseUrl         string
}

func NewLinkShortener(baseUrl string, shortLinkExists func(string) (bool, error)) *LinkShortener {
	linkShortener := LinkShortener{}
	linkShortener.baseUrl = baseUrl
	linkShortener.shortLinkExists = shortLinkExists
	return &linkShortener
}

func extendStringWithRandomCharacter(line string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomCharacter := letterBytes[rand.Intn(len(letterBytes))]
	return fmt.Sprintf("%s%c", line, randomCharacter)
}

func shortenUrl(link string, baseUrl string) string {
	result := md5.Sum([]byte(link))
	fiveBytes := result[0:5]
	encodedBytes := make([]byte, 10)
	base64.URLEncoding.Encode(encodedBytes, fiveBytes)
	return fmt.Sprintf("%s/%s", baseUrl, string(encodedBytes[0:5]))
}

func (o *LinkShortener) Shorten(link string) (string, error) {
	foundNonColliding := false
	var shortUrl string
	for !foundNonColliding {
		shortUrl = shortenUrl(link, o.baseUrl)
		var err error
		foundNonColliding, err = o.shortLinkExists(shortUrl)
		if err != nil {
			panic(err)
		}
		if !foundNonColliding {
			link = extendStringWithRandomCharacter(link)
		}
	}

	return shortUrl, nil
}
