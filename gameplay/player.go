package gameplay

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Player struct {
	Name            string `json:"name"`
	Uid             int64  `json:"uid"`
	ReviewGiven     bool   `json:"review_given"`
	SongSubmitted   bool   `json:"song_ubmited"`
	SongPresented   bool   `json:"song_presented"`
	Song            *Song
	ReceivedReviews []Review
}

type Song struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Review struct {
	Rating     int    `json:"rating"`
	FromPlayer string `json:"from_player"`
	UserReview string `json:"user_review"`
}

func CreatePlayer(name string, uid int64) *Player {
	p := Player{Name: name, Uid: uid}
	p.ReviewGiven = false
	p.SongSubmitted = false
	p.SongPresented = false
	return &p
}

func (p *Player) AddSong(description, url string) {
	s := &Song{Description: description, Url: url}
	p.Song = s
}

func parseRating(review string) (int, error) {
	ratingPat, _ := regexp.Compile("[0-9]+/[0-9]+$")
	match := ratingPat.FindString(review)
	if match == "" {
		return -1, errors.New("Couldn't parse points")
	}
	points := strings.Split(match, "/")
	numPoints, _ := strconv.Atoi(points[0])
	return numPoints, nil
}

func (p *Player) AddReview(from string, review string) error {

	rating, err := parseRating(review)
	if err != nil {
		log.Printf("Error parsing scores in rating: %s", review)
		return errors.New("Error parsing scores")
	}
	cleanedReview := strings.LastIndex(review, " ")
	if cleanedReview == -1 {
		log.Printf("Couldn't find the last space char from "+
			"the review: %s", review)
		return errors.New("Bogus review")
	}

	r := Review{Rating: rating,
		FromPlayer: from,
		UserReview: review[0:cleanedReview]}
	p.ReceivedReviews = append(p.ReceivedReviews, r)

	return nil
}
