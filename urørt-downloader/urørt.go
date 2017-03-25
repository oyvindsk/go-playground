package main

import "time"

//  https://mholt.github.io/json-to-go/

type urørtResponse struct {
	ID      string `json:"$id"`
	Type    string `json:"$type"`
	Results []struct {
		IDs                string        `json:"$id"`
		Type               string        `json:"$type"`
		ID                 int           `json:"Id"`
		Title              string        `json:"Title"`
		Composer           string        `json:"Composer"`
		Songwriter         string        `json:"Songwriter"`
		Released           time.Time     `json:"Released"`
		TrackState         int           `json:"TrackState"`
		Recommended        time.Time     `json:"Recommended"`
		LikeCount          int           `json:"LikeCount"`
		PlayCount          int           `json:"PlayCount"`
		BandID             int           `json:"BandId"`
		BandName           string        `json:"BandName"`
		InternalBandURL    string        `json:"InternalBandUrl"`
		Image              string        `json:"Image"`
		PlayedOnRadioCount int           `json:"PlayedOnRadioCount"`
		IsPlayable         bool          `json:"IsPlayable"`
		Tags               []interface{} `json:"Tags"`
		Files              []struct {
			IDs      string `json:"$id"`
			Type     string `json:"$type"`
			ID       int    `json:"Id"`
			TrackID  int    `json:"TrackId"`
			FileRef  string `json:"FileRef"`
			FileType string `json:"FileType"`
		} `json:"Files"`
		CommentCount int `json:"CommentCount,omitempty"`
	} `json:"Results"`
	InlineCount int `json:"InlineCount"`
}
