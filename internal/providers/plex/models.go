package plex

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
	KeyTag         string  `json:"keyTag,omitempty"`
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
