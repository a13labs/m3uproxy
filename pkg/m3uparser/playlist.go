/*
Copyright © 2024 Alexandre Pires

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package m3uparser

import (
	"io"
	"strings"
)

// M3UPlaylist represents the parsed M3U playlist.
type M3UPlaylist struct {
	Version int        // The version of the M3U (EXTM3U).
	Entries M3UEntries // The list of media entries in the playlist.
	Tags    M3UTags    // Additional tags associated with the entry.
	Type    string     // The type of the media (if available).
}

func (playlist *M3UPlaylist) GetVersion() int {
	return playlist.Version
}

func (playlist *M3UPlaylist) GetEntries() M3UEntries {
	return playlist.Entries
}

func (playlist *M3UPlaylist) EntriesString() string {
	var result string
	for _, tag := range playlist.Tags {
		result += "#" + tag.Tag + ":" + tag.Value + "\n"
	}
	for _, entry := range playlist.Entries {
		result += entry.String() + "\n"
	}
	return strings.Trim(result, "\n")
}

func (playlist *M3UPlaylist) String() string {
	var result string
	result += "#EXTM3U\n"
	result += playlist.EntriesString()
	return result
}

func (playlist *M3UPlaylist) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write([]byte("#EXTM3U\n"))
	if err != nil {
		return int64(n), err
	}
	for _, tag := range playlist.Tags {
		nBytes, _ := writer.Write([]byte("#" + tag.Tag + ":" + tag.Value + "\n"))
		n += nBytes
	}
	for _, entry := range playlist.Entries {
		nBytes, _ := entry.WriteTo(writer)
		n += int(nBytes)
	}
	return int64(n), err
}

func (playlist *M3UPlaylist) SearchEntryByTitle(title string) *M3UEntry {
	for _, entry := range playlist.Entries {
		if entry.Title == title {
			return &entry
		}
	}
	return nil
}

func (playlist *M3UPlaylist) SearchEntryByURI(uri string) *M3UEntry {
	for _, entry := range playlist.Entries {
		if entry.URI == uri {
			return &entry
		}
	}
	return nil
}

func (playlist *M3UPlaylist) StreamCount() int {
	return len(playlist.Entries)
}

func (playlist *M3UPlaylist) SearchEntryByTvgTag(tag, value string) *M3UEntry {
	return playlist.Entries.SearchByTvgTag(tag, value)
}

func (playlist *M3UPlaylist) SearchEntryIndexByTvgTag(tag, value string) int {
	return playlist.Entries.SearchIndexByTvgTag(tag, value)
}

func (playlist *M3UPlaylist) RemoveEntryByTvgTag(tag, value string) {
	playlist.Entries.RemoveByTvgTag(tag, value)
}
