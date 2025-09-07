// Package xmltv provides Go structures for parsing and generating XMLTV formatted data by implementing
// all elements of the XMLTV DTD.
// XMLTV is an XML-based format widely used in electronic program guides (EPG) for describing TV listings.
package xmltv

import (
	"encoding/xml"

	"github.com/a13labs/m3uproxy/pkg/xmltv/types"
)

// EPG is an alias for TV.
type EPG = TV

// TV represents the root element of an XMLTV document.
type TV struct {
	XMLName           xml.Name         `xml:"tv"`
	Date              *types.XMLTVTime `xml:"date,attr,omitempty"`
	SourceInfoURL     *string          `xml:"source-info-url,attr,omitempty"`
	SourceInfoName    *string          `xml:"source-info-name,attr,omitempty"`
	SourceDataURL     *string          `xml:"source-data-url,attr,omitempty"`
	GeneratorInfoName *string          `xml:"generator-info-name,attr,omitempty"`
	GeneratorInfoURL  *string          `xml:"generator-info-url,attr,omitempty"`
	Channels          []Channel        `xml:"channel,omitempty"`
	Programmes        []Programme      `xml:"programme,omitempty"`
}

// Channel represents a channel.
type Channel struct {
	XMLName      xml.Name      `xml:"channel"`
	ID           string        `xml:"id,attr"`
	DisplayNames []DisplayName `xml:"display-name"`
	Icons        []Icon        `xml:"icon,omitempty"`
	URLs         []URL         `xml:"url,omitempty"`
}

type DisplayName struct {
	XMLName xml.Name `xml:"display-name"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Icon struct {
	XMLName xml.Name `xml:"icon"`
	Source  string   `xml:"src,attr"`
	Width   *int     `xml:"width,attr,omitempty"`
	Height  *int     `xml:"height,attr,omitempty"`
}

type URL struct {
	XMLName xml.Name `xml:"url"`
	System  *string  `xml:"system,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

// Programme represents a programme.
type Programme struct {
	XMLName          xml.Name          `xml:"programme"`
	Start            types.XMLTVTime   `xml:"start,attr"`
	Stop             *types.XMLTVTime  `xml:"stop,attr,omitempty"`
	PDCStart         *types.XMLTVTime  `xml:"pdc-start,attr,omitempty"`
	VPSStart         *types.XMLTVTime  `xml:"vps-start,attr,omitempty"`
	ShowView         *string           `xml:"showview,attr,omitempty"`
	VideoPlus        *string           `xml:"videoplus,attr,omitempty"`
	Channel          string            `xml:"channel,attr"`
	ClumpIndex       *string           `xml:"clumpidx,attr,omitempty"`
	Titles           []Title           `xml:"title"`
	SubTitles        []SubTitle        `xml:"sub-title,omitempty"`
	Descriptions     []Description     `xml:"desc,omitempty"`
	Credits          *Credits          `xml:"credits,omitempty"`
	Date             *types.XMLTVTime  `xml:"date,omitempty"`
	Categories       []Category        `xml:"category,omitempty"`
	Keywords         []Keyword         `xml:"keyword,omitempty"`
	Language         *Language         `xml:"language,omitempty"`
	OriginalLanguage *OriginalLanguage `xml:"orig-language,omitempty"`
	Length           *Length           `xml:"length,omitempty"`
	Icons            []Icon            `xml:"icon,omitempty"`
	URLs             []URL             `xml:"url,omitempty"`
	Countries        []Country         `xml:"country,omitempty"`
	EpisodeNumbers   []EpisodeNumber   `xml:"episode-num,omitempty"`
	Video            *Video            `xml:"video,omitempty"`
	Audio            *Audio            `xml:"audio,omitempty"`
	PreviouslyShown  *PreviouslyShown  `xml:"previously-shown,omitempty"`
	Premiere         *Premiere         `xml:"premiere,omitempty"`
	Lastchance       *LastChance       `xml:"last-chance,omitempty"`
	IsNew            types.XMLTVBool   `xml:"new,omitempty"`
	Subtitles        []Subtitles       `xml:"subtitles,omitempty"`
	Ratings          []Rating          `xml:"rating,omitempty"`
	StarRatings      []StarRating      `xml:"star-rating,omitempty"`
	Reviews          []Review          `xml:"review,omitempty"`
	Images           []Image           `xml:"image,omitempty"`
}

type Title struct {
	XMLName xml.Name `xml:"title"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type SubTitle struct {
	XMLName xml.Name `xml:"sub-title"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Description struct {
	XMLName xml.Name `xml:"desc"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Credits struct {
	XMLName      xml.Name      `xml:"credits"`
	Directors    []Director    `xml:"director,omitempty"`
	Actors       []Actor       `xml:"actor,omitempty"`
	Writers      []Writer      `xml:"writer,omitempty"`
	Adapters     []Adapter     `xml:"adapter,omitempty"`
	Producers    []Producer    `xml:"producer,omitempty"`
	Composers    []Composer    `xml:"composer,omitempty"`
	Editors      []Editor      `xml:"editor,omitempty"`
	Presenters   []Presenter   `xml:"presenter,omitempty"`
	Commentators []Commentator `xml:"commentator,omitempty"`
	Guests       []Guest       `xml:"guest,omitempty"`
}

type Director struct {
	XMLName xml.Name `xml:"director"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Actor struct {
	XMLName xml.Name         `xml:"actor"`
	Role    *string          `xml:"role,attr,omitempty"`
	IsGuest *types.XMLTVBool `xml:"guest,attr,omitempty"`
	Images  []Image          `xml:"image,omitempty"`
	URLs    []URL            `xml:"url,omitempty"`
	Text    string           `xml:",chardata"`
}

type Writer struct {
	XMLName xml.Name `xml:"writer"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Adapter struct {
	XMLName xml.Name `xml:"adapter"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Producer struct {
	XMLName xml.Name `xml:"producer"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Composer struct {
	XMLName xml.Name `xml:"composer"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Editor struct {
	XMLName xml.Name `xml:"editor"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Presenter struct {
	XMLName xml.Name `xml:"presenter"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Commentator struct {
	XMLName xml.Name `xml:"commentator"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Guest struct {
	XMLName xml.Name `xml:"guest"`
	Images  []Image  `xml:"image,omitempty"`
	URLs    []URL    `xml:"url,omitempty"`
	Text    string   `xml:",chardata"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Keyword struct {
	XMLName xml.Name `xml:"keyword"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Language struct {
	XMLName xml.Name `xml:"language"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type OriginalLanguage struct {
	XMLName xml.Name `xml:"orig-language"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type LengthUnits string

const (
	LengthUnitsSeconds LengthUnits = "seconds"
	LengthUnitsMinutes LengthUnits = "minutes"
	LengthUnitsHours   LengthUnits = "hours"
)

type Length struct {
	XMLName xml.Name    `xml:"length"`
	Units   LengthUnits `xml:"units,attr"`
	Text    *int        `xml:",chardata"`
}

type Country struct {
	XMLName xml.Name `xml:"country"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type EpisodeNumber struct {
	XMLName xml.Name `xml:"episode-num"`
	System  string   `xml:"system,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Video struct {
	XMLName xml.Name         `xml:"video"`
	Present *types.XMLTVBool `xml:"present,omitempty"`
	Colour  *types.XMLTVBool `xml:"colour,omitempty"`
	Aspect  *Aspect          `xml:"aspect,omitempty"`
	Quality *Quality         `xml:"quality,omitempty"`
}

type Aspect struct {
	XMLName xml.Name `xml:"aspect"`
	Text    string   `xml:",chardata"`
}

type Quality struct {
	XMLName xml.Name `xml:"quality"`
	Text    string   `xml:",chardata"`
}

type Audio struct {
	XMLName xml.Name         `xml:"audio"`
	Present *types.XMLTVBool `xml:"present,omitempty"`
	Stereo  *Stereo          `xml:"stereo,omitempty"`
}

type Stereo struct {
	XMLName xml.Name `xml:"stereo"`
	Text    string   `xml:",chardata"`
}

type PreviouslyShown struct {
	XMLName xml.Name         `xml:"previously-shown"`
	Start   *types.XMLTVTime `xml:"start,attr,omitempty"`
	Channel *string          `xml:"channel,attr,omitempty"`
}

type Premiere struct {
	XMLName xml.Name `xml:"premiere"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type LastChance struct {
	XMLName xml.Name `xml:"last-chance"`
	Lang    *string  `xml:"lang,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type SubtitlesType string

const (
	SubtitlesTypeTeletext   SubtitlesType = "teletext"
	SubtitlesTypeOnScreen   SubtitlesType = "onscreen"
	SubtitlesTypeDeafSigned SubtitlesType = "deaf-signed"
)

type Subtitles struct {
	XMLName  xml.Name       `xml:"subtitles"`
	Type     *SubtitlesType `xml:"type,attr,omitempty"`
	Language *Language      `xml:"language,omitempty"`
}

type Rating struct {
	XMLName xml.Name `xml:"rating"`
	System  *string  `xml:"system,attr,omitempty"`
	Value   *Value   `xml:"value,omitempty"`
	Icons   []Icon   `xml:"icon,omitempty"`
}

type Value struct {
	XMLName xml.Name `xml:"value"`
	Text    string   `xml:",chardata"`
}

type StarRating struct {
	XMLName xml.Name `xml:"star-rating"`
	System  *string  `xml:"system,attr,omitempty"`
	Value   *Value   `xml:"value,omitempty"`
	Icons   []Icon   `xml:"icon,omitempty"`
}

type ReviewType string

const (
	ReviewTypeText ReviewType = "text"
	ReviewTypeURL  ReviewType = "url"
)

type Review struct {
	XMLName  xml.Name    `xml:"review"`
	Type     *ReviewType `xml:"type,attr,omitempty"`
	Source   *string     `xml:"source,attr,omitempty"`
	Reviewer *string     `xml:"reviewer,attr,omitempty"`
	Lang     *string     `xml:"lang,attr,omitempty"`
	Text     string      `xml:",chardata"`
}

type ImageType string

const (
	ImageTypePoster    ImageType = "poster"
	ImageTypeBackdrop  ImageType = "backdrop"
	ImageTypeStill     ImageType = "still"
	ImageTypePerson    ImageType = "person"
	ImageTypeCharacter ImageType = "character"
)

type ImageSize int

const (
	ImageSizeSmall  ImageSize = 1
	ImageSizeMedium ImageSize = 2
	ImageSizeLarge  ImageSize = 3
)

type ImageOrientation string

const (
	ImageOrientationPortrait  ImageOrientation = "P"
	ImageOrientationLandscape ImageOrientation = "L"
)

type Image struct {
	XMLName     xml.Name          `xml:"image"`
	Type        *ImageType        `xml:"type,attr,omitempty"`
	Size        *ImageSize        `xml:"size,attr,omitempty"`
	Orientation *ImageOrientation `xml:"orient,attr,omitempty"`
	System      *string           `xml:"system,attr,omitempty"`
	Text        string            `xml:",chardata"`
}
