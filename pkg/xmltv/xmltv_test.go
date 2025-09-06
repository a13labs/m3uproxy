package xmltv

import (
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/a13labs/m3uproxy/pkg/xmltv/internal/types"
	"github.com/google/go-cmp/cmp"
)

func makePointer[T any](t T) *T {
	return &t
}

func parseTime(t *testing.T, layout string, value string) time.Time {
	t.Helper()

	tt, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return tt
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	want, err := os.ReadFile("testdata/marshal/epg.xml")
	if err != nil {
		t.Fatal(err)
	}

	epg := TV{
		XMLName: xml.Name{
			Space: "",
			Local: "tv",
		},
		Date: &types.XMLTVTime{
			Time: parseTime(t, "20060102150405 -0700", "20220401000000 +0000"),
		},
		SourceInfoURL:     makePointer("example.com"),
		SourceInfoName:    makePointer("example"),
		SourceDataURL:     makePointer("example.com/a"),
		GeneratorInfoName: makePointer("Example Generator"),
		GeneratorInfoURL:  makePointer("https://example.com"),
		Channels: []Channel{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-one.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("en"),
						Text: "Channel One",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("fr"),
						Text: "Chaîne un",
					},
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/channel_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("example"),
						Text:   "https://example.com/channel_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("other_system"),
						Text:   "https://example.com/channel_one_alternate",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-two.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: nil,
						Text: "Channel Two: Minimum valid channel",
					},
				},
				Icons: nil,
				URLs:  nil,
			},
		},
		Programmes: []Programme{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331190000 +0000"),
				},
				PDCStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				VPSStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				ShowView:   makePointer("12345"),
				VideoPlus:  makePointer("67890"),
				Channel:    "channel-one.tv",
				ClumpIndex: makePointer("0/1"),
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: makePointer("en"),
						Text: "Programme One",
					},
				},
				SubTitles: []SubTitle{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "sub-title",
						},
						Lang: makePointer("en"),
						Text: "Pilot",
					},
				},
				Descriptions: []Description{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "This programme entry showcases all possible features of the DTD",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "Short description",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("cy"),
						Text: "Mae'r cofnod rhaglen hwn yn arddangos holl nodweddion posibl y DTD",
					},
				},
				Credits: &Credits{
					XMLName: xml.Name{
						Space: "",
						Local: "credits",
					},
					Directors: []Director{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "director",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Actors: []Actor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Walter Johnson"),
							IsGuest: nil,
							Images:  nil,
							URLs:    nil,
							Text:    "David Thompson",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Karl James"),
							IsGuest: makePointer(types.XMLTVBool(true)),
							Images: []Image{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "image",
									},
									Type:        makePointer(ImageTypePerson),
									Size:        nil,
									Orientation: nil,
									System:      nil,
									Text:        "\n                https://example.com/xxx.jpg",
								},
							},
							URLs: []URL{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "url",
									},
									System: makePointer("moviedb"),
									Text:   "\n                https://example.com/person/204",
								},
							},
							Text: " Ryan Lee \n        \n            ",
						},
					},
					Writers: []Writer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "writer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Adapters: []Adapter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "adapter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "William Brown",
						},
					},
					Producers: []Producer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "producer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Davis",
						},
					},
					Composers: []Composer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "composer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Max Wright",
						},
					},
					Editors: []Editor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "editor",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Nora Peterson",
						},
					},
					Presenters: []Presenter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "presenter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Amanda Johnson",
						},
					},
					Commentators: []Commentator{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "commentator",
							},
							Images: nil,
							URLs:   nil,
							Text:   "James Wilson",
						},
					},
					Guests: []Guest{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Lucas Martin",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Parker",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Oliver Nelson",
						},
					},
				},
				Date: &types.XMLTVTime{
					Time: parseTime(t, "20060102", "19901011"),
				},
				Categories: []Category{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Drama",
					},
				},
				Keywords: []Keyword{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "physical-comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "romantic",
					},
				},
				Language: &Language{
					XMLName: xml.Name{
						Space: "",
						Local: "language",
					},
					Lang: nil,
					Text: "English",
				},
				OriginalLanguage: &OriginalLanguage{
					XMLName: xml.Name{
						Space: "",
						Local: "orig-language",
					},
					Lang: makePointer("en"),
					Text: "French",
				},
				Length: &Length{
					XMLName: xml.Name{
						Space: "",
						Local: "length",
					},
					Units: LengthUnitsMinutes,
					Text:  makePointer(60),
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/programme_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("imdb"),
						Text:   "https://example.com/programme_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: nil,
						Text:   "https://example.com/programme_one_2",
					},
				},
				Countries: []Country{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "country",
						},
						Lang: nil,
						Text: "US",
					},
				},
				EpisodeNumbers: []EpisodeNumber{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "onscreen",
						Text:   "S01E01",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "xmltv_ns",
						Text:   "1 . 1 . 0/1",
					},
				},
				Video: &Video{
					XMLName: xml.Name{
						Space: "",
						Local: "video",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Colour:  makePointer(types.XMLTVBool(false)),
					Aspect: &Aspect{
						XMLName: xml.Name{
							Space: "",
							Local: "aspect",
						},
						Text: "16:9",
					},
					Quality: &Quality{
						XMLName: xml.Name{
							Space: "",
							Local: "quality",
						},
						Text: "HDTV",
					},
				},
				Audio: &Audio{
					XMLName: xml.Name{
						Space: "",
						Local: "audio",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Stereo: &Stereo{
						XMLName: xml.Name{
							Space: "",
							Local: "stereo",
						},
						Text: "Dolby Digital",
					},
				},
				PreviouslyShown: &PreviouslyShown{
					XMLName: xml.Name{
						Space: "",
						Local: "previously-shown",
					},
					Start: &types.XMLTVTime{
						Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
					},
					Channel: makePointer("channel-two.tv"),
				},
				Premiere: &Premiere{
					XMLName: xml.Name{
						Space: "",
						Local: "premiere",
					},
					Lang: nil,
					Text: "First time on British TV",
				},
				Lastchance: &LastChance{
					XMLName: xml.Name{
						Space: "",
						Local: "last-chance",
					},
					Lang: makePointer("en"),
					Text: "Last time on this channel",
				},
				IsNew: types.XMLTVBool(true),
				Subtitles: []Subtitles{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeTeletext),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: nil,
							Text: "English",
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeOnScreen),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: makePointer("en"),
							Text: "Spanish",
						},
					},
				},
				Ratings: []Rating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("BBFC"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "15",
						},
						Icons: nil,
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("MPAA"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "NC-17",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "NC-17_symbol.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
				},
				StarRatings: []StarRating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("TV Guide"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "4/5",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "stars.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("IMDB"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "8/10",
						},
						Icons: nil,
					},
				},
				Reviews: []Review{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "This is a\n            fantastic show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("IDMB"),
						Reviewer: makePointer("Jane Doe"),
						Lang:     makePointer("en"),
						Text:     "I love this show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeURL),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "\n            https://example.com/programme_one_review",
					},
				},
				Images: []Image{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeSmall),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_poster_1.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeMedium),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_poster_2.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_backdrop_3.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_backdrop_3.jpg",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop:       nil,
				PDCStart:   nil,
				VPSStart:   nil,
				ShowView:   nil,
				VideoPlus:  nil,
				Channel:    "channel-two.tv",
				ClumpIndex: nil,
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: nil,
						Text: "Programme Two: The minimum valid programme",
					},
				},
				SubTitles:        nil,
				Descriptions:     nil,
				Credits:          nil,
				Date:             nil,
				Categories:       nil,
				Keywords:         nil,
				Language:         nil,
				OriginalLanguage: nil,
				Length:           nil,
				Icons:            nil,
				URLs:             nil,
				Countries:        nil,
				EpisodeNumbers:   nil,
				Video:            nil,
				Audio:            nil,
				PreviouslyShown:  nil,
				Premiere:         nil,
				Lastchance:       nil,
				IsNew:            types.XMLTVBool(false),
				Subtitles:        nil,
				Ratings:          nil,
				StarRatings:      nil,
				Reviews:          nil,
				Images:           nil,
			},
		},
	}
	got, err := xml.Marshal(epg)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/unmarshal/epg.xml")
	if err != nil {
		t.Fatal(err)
	}

	want := TV{
		XMLName: xml.Name{
			Space: "",
			Local: "tv",
		},
		Date: &types.XMLTVTime{
			Time: parseTime(t, "20060102150405 -0700", "20220401000000 +0000"),
		},
		SourceInfoURL:     makePointer("example.com"),
		SourceInfoName:    makePointer("example"),
		SourceDataURL:     makePointer("example.com/a"),
		GeneratorInfoName: makePointer("Example Generator"),
		GeneratorInfoURL:  makePointer("https://example.com"),
		Channels: []Channel{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-one.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("en"),
						Text: "Channel One",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("fr"),
						Text: "Chaîne un",
					},
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/channel_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("example"),
						Text:   "https://example.com/channel_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("other_system"),
						Text:   "https://example.com/channel_one_alternate",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-two.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: nil,
						Text: "Channel Two: Minimum valid channel",
					},
				},
				Icons: nil,
				URLs:  nil,
			},
		},
		Programmes: []Programme{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331190000 +0000"),
				},
				PDCStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				VPSStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				ShowView:   makePointer("12345"),
				VideoPlus:  makePointer("67890"),
				Channel:    "channel-one.tv",
				ClumpIndex: makePointer("0/1"),
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: makePointer("en"),
						Text: "Programme One",
					},
				},
				SubTitles: []SubTitle{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "sub-title",
						},
						Lang: makePointer("en"),
						Text: "Pilot",
					},
				},
				Descriptions: []Description{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "This programme entry showcases all possible features of the DTD",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "Short description",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("cy"),
						Text: "Mae'r cofnod rhaglen hwn yn arddangos holl nodweddion posibl y DTD",
					},
				},
				Credits: &Credits{
					XMLName: xml.Name{
						Space: "",
						Local: "credits",
					},
					Directors: []Director{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "director",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Actors: []Actor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Walter Johnson"),
							IsGuest: nil,
							Images:  nil,
							URLs:    nil,
							Text:    "David Thompson",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Karl James"),
							IsGuest: makePointer(types.XMLTVBool(true)),
							Images: []Image{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "image",
									},
									Type:        makePointer(ImageTypePerson),
									Size:        nil,
									Orientation: nil,
									System:      nil,
									Text:        "\n                https://example.com/xxx.jpg",
								},
							},
							URLs: []URL{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "url",
									},
									System: makePointer("moviedb"),
									Text:   "\n                https://example.com/person/204",
								},
							},
							Text: " Ryan Lee \n        \n            ",
						},
					},
					Writers: []Writer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "writer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Adapters: []Adapter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "adapter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "William Brown",
						},
					},
					Producers: []Producer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "producer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Davis",
						},
					},
					Composers: []Composer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "composer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Max Wright",
						},
					},
					Editors: []Editor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "editor",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Nora Peterson",
						},
					},
					Presenters: []Presenter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "presenter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Amanda Johnson",
						},
					},
					Commentators: []Commentator{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "commentator",
							},
							Images: nil,
							URLs:   nil,
							Text:   "James Wilson",
						},
					},
					Guests: []Guest{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Lucas Martin",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Parker",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Oliver Nelson",
						},
					},
				},
				Date: &types.XMLTVTime{
					Time: parseTime(t, "20060102", "19901011"),
				},
				Categories: []Category{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Drama",
					},
				},
				Keywords: []Keyword{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "physical-comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "romantic",
					},
				},
				Language: &Language{
					XMLName: xml.Name{
						Space: "",
						Local: "language",
					},
					Lang: nil,
					Text: "English",
				},
				OriginalLanguage: &OriginalLanguage{
					XMLName: xml.Name{
						Space: "",
						Local: "orig-language",
					},
					Lang: makePointer("en"),
					Text: "French",
				},
				Length: &Length{
					XMLName: xml.Name{
						Space: "",
						Local: "length",
					},
					Units: LengthUnitsMinutes,
					Text:  makePointer(60),
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/programme_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("imdb"),
						Text:   "https://example.com/programme_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: nil,
						Text:   "https://example.com/programme_one_2",
					},
				},
				Countries: []Country{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "country",
						},
						Lang: nil,
						Text: "US",
					},
				},
				EpisodeNumbers: []EpisodeNumber{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "onscreen",
						Text:   "S01E01",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "xmltv_ns",
						Text:   "1 . 1 . 0/1",
					},
				},
				Video: &Video{
					XMLName: xml.Name{
						Space: "",
						Local: "video",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Colour:  makePointer(types.XMLTVBool(false)),
					Aspect: &Aspect{
						XMLName: xml.Name{
							Space: "",
							Local: "aspect",
						},
						Text: "16:9",
					},
					Quality: &Quality{
						XMLName: xml.Name{
							Space: "",
							Local: "quality",
						},
						Text: "HDTV",
					},
				},
				Audio: &Audio{
					XMLName: xml.Name{
						Space: "",
						Local: "audio",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Stereo: &Stereo{
						XMLName: xml.Name{
							Space: "",
							Local: "stereo",
						},
						Text: "Dolby Digital",
					},
				},
				PreviouslyShown: &PreviouslyShown{
					XMLName: xml.Name{
						Space: "",
						Local: "previously-shown",
					},
					Start: &types.XMLTVTime{
						Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
					},
					Channel: makePointer("channel-two.tv"),
				},
				Premiere: &Premiere{
					XMLName: xml.Name{
						Space: "",
						Local: "premiere",
					},
					Lang: nil,
					Text: "First time on British TV",
				},
				Lastchance: &LastChance{
					XMLName: xml.Name{
						Space: "",
						Local: "last-chance",
					},
					Lang: makePointer("en"),
					Text: "Last time on this channel",
				},
				IsNew: types.XMLTVBool(true),
				Subtitles: []Subtitles{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeTeletext),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: nil,
							Text: "English",
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeOnScreen),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: makePointer("en"),
							Text: "Spanish",
						},
					},
				},
				Ratings: []Rating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("BBFC"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "15",
						},
						Icons: nil,
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("MPAA"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "NC-17",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "NC-17_symbol.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
				},
				StarRatings: []StarRating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("TV Guide"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "4/5",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "stars.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("IMDB"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "8/10",
						},
						Icons: nil,
					},
				},
				Reviews: []Review{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "This is a\n            fantastic show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("IDMB"),
						Reviewer: makePointer("Jane Doe"),
						Lang:     makePointer("en"),
						Text:     "I love this show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeURL),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "\n            https://example.com/programme_one_review",
					},
				},
				Images: []Image{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeSmall),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_poster_1.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeMedium),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_poster_2.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_backdrop_3.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_backdrop_3.jpg",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop:       nil,
				PDCStart:   nil,
				VPSStart:   nil,
				ShowView:   nil,
				VideoPlus:  nil,
				Channel:    "channel-two.tv",
				ClumpIndex: nil,
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: nil,
						Text: "Programme Two: The minimum valid programme",
					},
				},
				SubTitles:        nil,
				Descriptions:     nil,
				Credits:          nil,
				Date:             nil,
				Categories:       nil,
				Keywords:         nil,
				Language:         nil,
				OriginalLanguage: nil,
				Length:           nil,
				Icons:            nil,
				URLs:             nil,
				Countries:        nil,
				EpisodeNumbers:   nil,
				Video:            nil,
				Audio:            nil,
				PreviouslyShown:  nil,
				Premiere:         nil,
				Lastchance:       nil,
				IsNew:            types.XMLTVBool(false),
				Subtitles:        nil,
				Ratings:          nil,
				StarRatings:      nil,
				Reviews:          nil,
				Images:           nil,
			},
		},
	}
	var got TV
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestMarshalEmptyDate(t *testing.T) {
	t.Parallel()

	want, err := os.ReadFile("testdata/marshal/epg-empty-date.xml")
	if err != nil {
		t.Fatal(err)
	}

	epg := TV{
		XMLName: xml.Name{
			Space: "",
			Local: "tv",
		},
		Date:              &types.XMLTVTime{},
		SourceInfoURL:     makePointer("example.com"),
		SourceInfoName:    makePointer("example"),
		SourceDataURL:     makePointer("example.com/a"),
		GeneratorInfoName: makePointer("Example Generator"),
		GeneratorInfoURL:  makePointer("https://example.com"),
		Channels: []Channel{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-one.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("en"),
						Text: "Channel One",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("fr"),
						Text: "Chaîne un",
					},
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/channel_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("example"),
						Text:   "https://example.com/channel_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("other_system"),
						Text:   "https://example.com/channel_one_alternate",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-two.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: nil,
						Text: "Channel Two: Minimum valid channel",
					},
				},
				Icons: nil,
				URLs:  nil,
			},
		},
		Programmes: []Programme{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331190000 +0000"),
				},
				PDCStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				VPSStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				ShowView:   makePointer("12345"),
				VideoPlus:  makePointer("67890"),
				Channel:    "channel-one.tv",
				ClumpIndex: makePointer("0/1"),
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: makePointer("en"),
						Text: "Programme One",
					},
				},
				SubTitles: []SubTitle{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "sub-title",
						},
						Lang: makePointer("en"),
						Text: "Pilot",
					},
				},
				Descriptions: []Description{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "This programme entry showcases all possible features of the DTD",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "Short description",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("cy"),
						Text: "Mae'r cofnod rhaglen hwn yn arddangos holl nodweddion posibl y DTD",
					},
				},
				Credits: &Credits{
					XMLName: xml.Name{
						Space: "",
						Local: "credits",
					},
					Directors: []Director{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "director",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Actors: []Actor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Walter Johnson"),
							IsGuest: nil,
							Images:  nil,
							URLs:    nil,
							Text:    "David Thompson",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Karl James"),
							IsGuest: makePointer(types.XMLTVBool(true)),
							Images: []Image{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "image",
									},
									Type:        makePointer(ImageTypePerson),
									Size:        nil,
									Orientation: nil,
									System:      nil,
									Text:        "\n                https://example.com/xxx.jpg",
								},
							},
							URLs: []URL{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "url",
									},
									System: makePointer("moviedb"),
									Text:   "\n                https://example.com/person/204",
								},
							},
							Text: " Ryan Lee \n        \n            ",
						},
					},
					Writers: []Writer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "writer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Adapters: []Adapter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "adapter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "William Brown",
						},
					},
					Producers: []Producer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "producer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Davis",
						},
					},
					Composers: []Composer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "composer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Max Wright",
						},
					},
					Editors: []Editor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "editor",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Nora Peterson",
						},
					},
					Presenters: []Presenter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "presenter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Amanda Johnson",
						},
					},
					Commentators: []Commentator{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "commentator",
							},
							Images: nil,
							URLs:   nil,
							Text:   "James Wilson",
						},
					},
					Guests: []Guest{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Lucas Martin",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Parker",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Oliver Nelson",
						},
					},
				},
				Date: &types.XMLTVTime{},
				Categories: []Category{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Drama",
					},
				},
				Keywords: []Keyword{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "physical-comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "romantic",
					},
				},
				Language: &Language{
					XMLName: xml.Name{
						Space: "",
						Local: "language",
					},
					Lang: nil,
					Text: "English",
				},
				OriginalLanguage: &OriginalLanguage{
					XMLName: xml.Name{
						Space: "",
						Local: "orig-language",
					},
					Lang: makePointer("en"),
					Text: "French",
				},
				Length: &Length{
					XMLName: xml.Name{
						Space: "",
						Local: "length",
					},
					Units: LengthUnitsMinutes,
					Text:  makePointer(60),
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/programme_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("imdb"),
						Text:   "https://example.com/programme_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: nil,
						Text:   "https://example.com/programme_one_2",
					},
				},
				Countries: []Country{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "country",
						},
						Lang: nil,
						Text: "US",
					},
				},
				EpisodeNumbers: []EpisodeNumber{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "onscreen",
						Text:   "S01E01",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "xmltv_ns",
						Text:   "1 . 1 . 0/1",
					},
				},
				Video: &Video{
					XMLName: xml.Name{
						Space: "",
						Local: "video",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Colour:  makePointer(types.XMLTVBool(false)),
					Aspect: &Aspect{
						XMLName: xml.Name{
							Space: "",
							Local: "aspect",
						},
						Text: "16:9",
					},
					Quality: &Quality{
						XMLName: xml.Name{
							Space: "",
							Local: "quality",
						},
						Text: "HDTV",
					},
				},
				Audio: &Audio{
					XMLName: xml.Name{
						Space: "",
						Local: "audio",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Stereo: &Stereo{
						XMLName: xml.Name{
							Space: "",
							Local: "stereo",
						},
						Text: "Dolby Digital",
					},
				},
				PreviouslyShown: &PreviouslyShown{
					XMLName: xml.Name{
						Space: "",
						Local: "previously-shown",
					},
					Start: &types.XMLTVTime{
						Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
					},
					Channel: makePointer("channel-two.tv"),
				},
				Premiere: &Premiere{
					XMLName: xml.Name{
						Space: "",
						Local: "premiere",
					},
					Lang: nil,
					Text: "First time on British TV",
				},
				Lastchance: &LastChance{
					XMLName: xml.Name{
						Space: "",
						Local: "last-chance",
					},
					Lang: makePointer("en"),
					Text: "Last time on this channel",
				},
				IsNew: types.XMLTVBool(true),
				Subtitles: []Subtitles{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeTeletext),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: nil,
							Text: "English",
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeOnScreen),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: makePointer("en"),
							Text: "Spanish",
						},
					},
				},
				Ratings: []Rating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("BBFC"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "15",
						},
						Icons: nil,
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("MPAA"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "NC-17",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "NC-17_symbol.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
				},
				StarRatings: []StarRating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("TV Guide"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "4/5",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "stars.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("IMDB"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "8/10",
						},
						Icons: nil,
					},
				},
				Reviews: []Review{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "This is a\n            fantastic show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("IDMB"),
						Reviewer: makePointer("Jane Doe"),
						Lang:     makePointer("en"),
						Text:     "I love this show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeURL),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "\n            https://example.com/programme_one_review",
					},
				},
				Images: []Image{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeSmall),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_poster_1.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeMedium),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_poster_2.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_backdrop_3.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_backdrop_3.jpg",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop:       nil,
				PDCStart:   nil,
				VPSStart:   nil,
				ShowView:   nil,
				VideoPlus:  nil,
				Channel:    "channel-two.tv",
				ClumpIndex: nil,
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: nil,
						Text: "Programme Two: The minimum valid programme",
					},
				},
				SubTitles:        nil,
				Descriptions:     nil,
				Credits:          nil,
				Date:             nil,
				Categories:       nil,
				Keywords:         nil,
				Language:         nil,
				OriginalLanguage: nil,
				Length:           nil,
				Icons:            nil,
				URLs:             nil,
				Countries:        nil,
				EpisodeNumbers:   nil,
				Video:            nil,
				Audio:            nil,
				PreviouslyShown:  nil,
				Premiere:         nil,
				Lastchance:       nil,
				IsNew:            types.XMLTVBool(false),
				Subtitles:        nil,
				Ratings:          nil,
				StarRatings:      nil,
				Reviews:          nil,
				Images:           nil,
			},
		},
	}
	got, err := xml.Marshal(epg)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestUnmarshalEmptyDate(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/unmarshal/epg-empty-date.xml")
	if err != nil {
		t.Fatal(err)
	}

	want := TV{
		XMLName: xml.Name{
			Space: "",
			Local: "tv",
		},
		Date:              &types.XMLTVTime{},
		SourceInfoURL:     makePointer("example.com"),
		SourceInfoName:    makePointer("example"),
		SourceDataURL:     makePointer("example.com/a"),
		GeneratorInfoName: makePointer("Example Generator"),
		GeneratorInfoURL:  makePointer("https://example.com"),
		Channels: []Channel{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-one.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("en"),
						Text: "Channel One",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: makePointer("fr"),
						Text: "Chaîne un",
					},
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/channel_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("example"),
						Text:   "https://example.com/channel_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("other_system"),
						Text:   "https://example.com/channel_one_alternate",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				ID: "channel-two.tv",
				DisplayNames: []DisplayName{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "display-name",
						},
						Lang: nil,
						Text: "Channel Two: Minimum valid channel",
					},
				},
				Icons: nil,
				URLs:  nil,
			},
		},
		Programmes: []Programme{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331190000 +0000"),
				},
				PDCStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				VPSStart: &types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				ShowView:   makePointer("12345"),
				VideoPlus:  makePointer("67890"),
				Channel:    "channel-one.tv",
				ClumpIndex: makePointer("0/1"),
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: makePointer("en"),
						Text: "Programme One",
					},
				},
				SubTitles: []SubTitle{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "sub-title",
						},
						Lang: makePointer("en"),
						Text: "Pilot",
					},
				},
				Descriptions: []Description{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "This programme entry showcases all possible features of the DTD",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("en"),
						Text: "Short description",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "desc",
						},
						Lang: makePointer("cy"),
						Text: "Mae'r cofnod rhaglen hwn yn arddangos holl nodweddion posibl y DTD",
					},
				},
				Credits: &Credits{
					XMLName: xml.Name{
						Space: "",
						Local: "credits",
					},
					Directors: []Director{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "director",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Actors: []Actor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Walter Johnson"),
							IsGuest: nil,
							Images:  nil,
							URLs:    nil,
							Text:    "David Thompson",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "actor",
							},
							Role:    makePointer("Karl James"),
							IsGuest: makePointer(types.XMLTVBool(true)),
							Images: []Image{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "image",
									},
									Type:        makePointer(ImageTypePerson),
									Size:        nil,
									Orientation: nil,
									System:      nil,
									Text:        "\n                https://example.com/xxx.jpg",
								},
							},
							URLs: []URL{
								{
									XMLName: xml.Name{
										Space: "",
										Local: "url",
									},
									System: makePointer("moviedb"),
									Text:   "\n                https://example.com/person/204",
								},
							},
							Text: " Ryan Lee \n        \n            ",
						},
					},
					Writers: []Writer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "writer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Samuel Jones",
						},
					},
					Adapters: []Adapter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "adapter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "William Brown",
						},
					},
					Producers: []Producer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "producer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Davis",
						},
					},
					Composers: []Composer{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "composer",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Max Wright",
						},
					},
					Editors: []Editor{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "editor",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Nora Peterson",
						},
					},
					Presenters: []Presenter{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "presenter",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Amanda Johnson",
						},
					},
					Commentators: []Commentator{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "commentator",
							},
							Images: nil,
							URLs:   nil,
							Text:   "James Wilson",
						},
					},
					Guests: []Guest{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Lucas Martin",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Emily Parker",
						},
						{
							XMLName: xml.Name{
								Space: "",
								Local: "guest",
							},
							Images: nil,
							URLs:   nil,
							Text:   "Oliver Nelson",
						},
					},
				},
				Date: &types.XMLTVTime{},
				Categories: []Category{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "category",
						},
						Lang: makePointer("en"),
						Text: "Drama",
					},
				},
				Keywords: []Keyword{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "physical-comedy",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "keyword",
						},
						Lang: makePointer("en"),
						Text: "romantic",
					},
				},
				Language: &Language{
					XMLName: xml.Name{
						Space: "",
						Local: "language",
					},
					Lang: nil,
					Text: "English",
				},
				OriginalLanguage: &OriginalLanguage{
					XMLName: xml.Name{
						Space: "",
						Local: "orig-language",
					},
					Lang: makePointer("en"),
					Text: "French",
				},
				Length: &Length{
					XMLName: xml.Name{
						Space: "",
						Local: "length",
					},
					Units: LengthUnitsMinutes,
					Text:  makePointer(60),
				},
				Icons: []Icon{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "icon",
						},
						Source: "https://example.com/programme_one_icon.jpg",
						Width:  makePointer(100),
						Height: makePointer(100),
					},
				},
				URLs: []URL{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: makePointer("imdb"),
						Text:   "https://example.com/programme_one",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "url",
						},
						System: nil,
						Text:   "https://example.com/programme_one_2",
					},
				},
				Countries: []Country{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "country",
						},
						Lang: nil,
						Text: "US",
					},
				},
				EpisodeNumbers: []EpisodeNumber{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "onscreen",
						Text:   "S01E01",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "episode-num",
						},
						System: "xmltv_ns",
						Text:   "1 . 1 . 0/1",
					},
				},
				Video: &Video{
					XMLName: xml.Name{
						Space: "",
						Local: "video",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Colour:  makePointer(types.XMLTVBool(false)),
					Aspect: &Aspect{
						XMLName: xml.Name{
							Space: "",
							Local: "aspect",
						},
						Text: "16:9",
					},
					Quality: &Quality{
						XMLName: xml.Name{
							Space: "",
							Local: "quality",
						},
						Text: "HDTV",
					},
				},
				Audio: &Audio{
					XMLName: xml.Name{
						Space: "",
						Local: "audio",
					},
					Present: makePointer(types.XMLTVBool(true)),
					Stereo: &Stereo{
						XMLName: xml.Name{
							Space: "",
							Local: "stereo",
						},
						Text: "Dolby Digital",
					},
				},
				PreviouslyShown: &PreviouslyShown{
					XMLName: xml.Name{
						Space: "",
						Local: "previously-shown",
					},
					Start: &types.XMLTVTime{
						Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
					},
					Channel: makePointer("channel-two.tv"),
				},
				Premiere: &Premiere{
					XMLName: xml.Name{
						Space: "",
						Local: "premiere",
					},
					Lang: nil,
					Text: "First time on British TV",
				},
				Lastchance: &LastChance{
					XMLName: xml.Name{
						Space: "",
						Local: "last-chance",
					},
					Lang: makePointer("en"),
					Text: "Last time on this channel",
				},
				IsNew: types.XMLTVBool(true),
				Subtitles: []Subtitles{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeTeletext),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: nil,
							Text: "English",
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "subtitles",
						},
						Type: makePointer(SubtitlesTypeOnScreen),
						Language: &Language{
							XMLName: xml.Name{
								Space: "",
								Local: "language",
							},
							Lang: makePointer("en"),
							Text: "Spanish",
						},
					},
				},
				Ratings: []Rating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("BBFC"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "15",
						},
						Icons: nil,
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "rating",
						},
						System: makePointer("MPAA"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "NC-17",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "NC-17_symbol.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
				},
				StarRatings: []StarRating{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("TV Guide"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "4/5",
						},
						Icons: []Icon{
							{
								XMLName: xml.Name{
									Space: "",
									Local: "icon",
								},
								Source: "stars.png",
								Width:  nil,
								Height: nil,
							},
						},
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "star-rating",
						},
						System: makePointer("IMDB"),
						Value: &Value{
							XMLName: xml.Name{
								Space: "",
								Local: "value",
							},
							Text: "8/10",
						},
						Icons: nil,
					},
				},
				Reviews: []Review{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "This is a\n            fantastic show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeText),
						Source:   makePointer("IDMB"),
						Reviewer: makePointer("Jane Doe"),
						Lang:     makePointer("en"),
						Text:     "I love this show!",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "review",
						},
						Type:     makePointer(ReviewTypeURL),
						Source:   makePointer("Rotten Tomatoes"),
						Reviewer: makePointer("Joe Bloggs"),
						Lang:     makePointer("en"),
						Text:     "\n            https://example.com/programme_one_review",
					},
				},
				Images: []Image{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeSmall),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_poster_1.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypePoster),
						Size:        makePointer(ImageSizeMedium),
						Orientation: makePointer(ImageOrientationPortrait),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_poster_2.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tvdb"),
						Text:        "\n            https://tvdb.com/programme_one_backdrop_3.jpg",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "image",
						},
						Type:        makePointer(ImageTypeBackdrop),
						Size:        makePointer(ImageSizeLarge),
						Orientation: makePointer(ImageOrientationLandscape),
						System:      makePointer("tmdb"),
						Text:        "\n            https://tmdb.com/programme_one_backdrop_3.jpg",
					},
				},
			},
			{
				XMLName: xml.Name{
					Space: "",
					Local: "programme",
				},
				Start: types.XMLTVTime{
					Time: parseTime(t, "20060102150405 -0700", "20220331180000 +0000"),
				},
				Stop:       nil,
				PDCStart:   nil,
				VPSStart:   nil,
				ShowView:   nil,
				VideoPlus:  nil,
				Channel:    "channel-two.tv",
				ClumpIndex: nil,
				Titles: []Title{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "title",
						},
						Lang: nil,
						Text: "Programme Two: The minimum valid programme",
					},
				},
				SubTitles:        nil,
				Descriptions:     nil,
				Credits:          nil,
				Date:             nil,
				Categories:       nil,
				Keywords:         nil,
				Language:         nil,
				OriginalLanguage: nil,
				Length:           nil,
				Icons:            nil,
				URLs:             nil,
				Countries:        nil,
				EpisodeNumbers:   nil,
				Video:            nil,
				Audio:            nil,
				PreviouslyShown:  nil,
				Premiere:         nil,
				Lastchance:       nil,
				IsNew:            types.XMLTVBool(false),
				Subtitles:        nil,
				Ratings:          nil,
				StarRatings:      nil,
				Reviews:          nil,
				Images:           nil,
			},
		},
	}
	var got TV
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}
