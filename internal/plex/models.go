package plex

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type Resource struct {
	Name                   string       `json:"name"`
	ClientIdentifier       string       `json:"clientIdentifier"`
	Product                string       `json:"product"`
	ProductVersion         string       `json:"productVersion"`
	Platform               string       `json:"platform"`
	PlatformVerson         string       `json:"platformVersion"`
	Device                 string       `json:"device"`
	CreatedAt              string       `json:"createdAt"`
	LastSeenAt             string       `json:"lastSeenAt"`
	Provides               string       `json:"provides"`
	OwnerID                int          `json:"ownerId"`
	SourceTitle            string       `json:"sourceTitle"`
	PublicAddress          string       `json:"publicAddress"`
	AccessToken            string       `json:"accessToken"`
	Owned                  bool         `json:"owned"`
	Home                   bool         `json:"home"`
	Synced                 bool         `json:"synced"`
	Relay                  bool         `json:"relay"`
	Presence               bool         `json:"presence"`
	HTTPSRequired          bool         `json:"httpsRequired"`
	PublicAddressMatches   bool         `json:"publicAddressMatches"`
	DNSRebindingProtection bool         `json:"dnsRebindingProtection"`
	NATLoopbackSupported   bool         `json:"natLoopbackSupported"`
	Connections            []Connection `json:"connections"`
}

type Connection struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	URI      string `json:"uri"`
	Local    bool   `json:"local"`
	Relay    bool   `json:"relay"`
	IPv6     bool   `json:"ipv6"`
}

type Playlist struct {
	RatingKey string `json:"ratingKey"`
	Key       string `json:"key"`
	GUID      string `json:"guid"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Smart     bool   `json:"smart"`
	LeafCount int    `json:"leafCount"`
}

type Track struct {
	RatingKey      string  `json:"ratingKey"`
	Key            string  `json:"key"`
	Type           string  `json:"type"`
	Title          string  `json:"title"`
	Summary        string  `json:"summary"`
	Artist         string  `json:"grandparentTitle"`
	Album          string  `json:"parentTitle"`
	BPM            float64 `json:"bpm"`
	KeyTag         string  `json:"key"`
	UserRating     float64 `json:"userRating"` // Plex usually 0-10 or 0-5
	Media          []Media `json:"Media"`
}

type Media struct {
	ID   int    `json:"id"`
	Part []Part `json:"Part"`
}

type Part struct {
	ID   int    `json:"id"`
	File string `json:"file"`
}

type MediaContainer struct {
	Metadata []Playlist `json:"Metadata"`
}

type TrackContainer struct {
	Metadata []Track `json:"Metadata"`
}

func (t Track) ToNeutral() models.Track {
	mt := models.Track{
		ID:     t.RatingKey,
		Title:  t.Title,
		Artist: t.Artist,
		Album:  t.Album,
		BPM:    t.BPM,
		Key:    t.KeyTag,
		Rating: models.NormalizeRating(t.UserRating, 10.0), // Plex uses a 10-point internal scale
		Raw:    t,
	}
	if len(t.Media) > 0 && len(t.Media[0].Part) > 0 {
		mt.Location = t.Media[0].Part[0].File
	}
	return mt
}

func (p Playlist) ToNeutralGroup() models.ResourceGroup {
	return models.ResourceGroup{
		Name:    p.Title,
		Items:   p.LeafCount,
		Type:    1,
		Raw:     p,
	}
}
