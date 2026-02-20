package lib

import "time"

type Copyright struct {
	Items      []map[string]interface{} `json:"items"`
	TotalCount int64                    `json:"totalCount"`
}

type ColorRaw struct {
	Hex string `json:"hex"`
}

type ExtractedColors struct {
	ColorRaw ColorRaw `json:"colorRaw"`
}

type CoverArt struct {
	ExtractedColors ExtractedColors `json:"extractedColors"`
	Sources         []struct {
		Height int    `json:"height"`
		Width  int    `json:"width"`
		Url    string `json:"url"`
	} `json:"sources"`
}

type Date struct {
	IsoString time.Time `json:"isoString"`
	Precision string    `json:"precision"`
	Year      int64     `json:"year"`
}

type SharingInfo struct {
	ShareId  string `json:"shareId"`
	ShareUrl string `json:"shareUrl"`
}

type Tracks struct {
	Items      []map[string]interface{} `json:"items"`
	TotalCount int64                    `json:"totalCount"`
}

type AlbumOfTrack struct {
	Copyright    Copyright   `json:"copyright"`
	CourtesyLine string      `json:"courtesyLine"`
	CoverArt     CoverArt    `json:"coverArt"`
	Date         Date        `json:"date"`
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Playability  Playability `json:"playability"`
	SharingInfo  SharingInfo `json:"sharingInfo"`
	Tracks       Tracks      `json:"tracks"`
	Type         string      `json:"type"`
	Uri          string      `json:"uri"`
}

type AudioAssociations struct {
	TypeName string        `json:"__typename"`
	Items    []interface{} `json:"items"`
}

type VideoAssociations struct {
	TotalCount int64 `json:"totalCount"`
}

type Associations struct {
	AudioAssociations AudioAssociations `json:"audioAssociations"`
	VideoAssociations VideoAssociations `json:"videoAssociations"`
}

type ContentRating struct {
	Label string `json:"label"`
}

type Duration struct {
	TotalMilliseconds int64 `json:"totalMilliseconds"`
}

type FirstArtist struct {
	Items      []ArtistItems `json:"items"`
	TotalCount int64         `json:"totalCount"`
}

type ArtistItems struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
}

type OtherArtists struct {
	Items []ArtistItems `json:"items"`
}

type Playability struct {
	Playable bool   `json:"playable"`
	Reason   string `json:"reason"`
}

type TrackUnion struct {
	TypeName       string        `json:"__typename"`
	AlbumOfTrack   AlbumOfTrack  `json:"albumOfTrack"`
	Associations   Associations  `json:"associationsV3"`
	ContentRating  ContentRating `json:"contentRating"`
	Duration       Duration      `json:"duration"`
	FirstArtist    FirstArtist   `json:"firstArtist"`
	Id             string        `json:"id"`
	MediaType      string        `json:"mediaType"`
	Name           string        `json:"name"`
	OtherArtists   OtherArtists  `json:"otherArtists"`
	Playability    Playability   `json:"playability"`
	Playcount      string        `json:"playcount"`
	Saved          bool          `json:"saved"`
	SharingInfo    interface{}   `json:"sharingInfo"`
	TrackNumber    int64         `json:"trackNumber"`
	Uri            string        `json:"uri"`
	VisualIdentity interface{}   `json:"visualIdentity"`
}

type Data struct {
	TrackUnion TrackUnion `json:"trackUnion"`
}

type TrackMetadata struct {
	Data Data `json:"data"`
}

type PlaylistMetadata struct {
	Data struct {
		Playlist struct {
			Name    string `json:"name"`
			Content struct {
				Items []struct {
					Item struct {
						Data struct {
							IdentityTrait struct {
								Name string `json:"name"`
							} `json:"identityTrait"`
							Uri string `json:"uri"`
						} `json:"data"`
					} `json:"itemV3"`
				} `json:"items"`
			} `json:"content"`
		} `json:"playlistV2"`
	} `json:"data"`
}
