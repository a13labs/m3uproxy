# XMLTV

A Go library for parsing and generating XMLTV formatted data.

## Overview

This package provides Go structures for parsing and generating XMLTV formatted data by implementing all elements of the [XMLTV DTD](https://raw.githubusercontent.com/XMLTV/xmltv/refs/heads/master/xmltv.dtd).

XMLTV is an XML-based format widely used in electronic program guides (EPG) for describing TV listings.

## Installation

```bash
go get github.com/sherif-fanous/xmltv
```

## Usage Examples

### Parsing XMLTV Data

```go
package main

import (
    "encoding/xml"
    "log"
    "os"

    "github.com/sherif-fanous/xmltv"
)

func main() {
    // Read XMLTV file
    data, err := os.ReadFile("epg.xml")
    if err != nil {
        log.Fatalf("Error reading file: %v", err)
    }

    // Parse XMLTV data
    var epg xmltv.EPG
    if err := xml.Unmarshal(data, &epg); err != nil {
        log.Fatalf("Error parsing XMLTV: %v", err)
    }

    // Access the data
    log.Printf("Found %d channels and %d programmes\n", len(epg.Channels), len(epg.Programmes))

    for _, channel := range epg.Channels {
        for _, displayName := range channel.DisplayNames {
            log.Printf("Channel name: %s\n", displayName.Text)
        }
    }
}
```

### Creating XMLTV Data

```go
package main

import (
    "encoding/xml"
    "io/fs"
    "log"
    "os"
    "time"

    "github.com/sherif-fanous/xmltv"
)

func makePointer[T any](t T) *T {
    return &t
}

func main() {
    now := time.Now()

    // Create a new XMLTV document
    epg := xmltv.EPG{
        GeneratorInfoName: makePointer("My Example Generator"),
        Channels: []xmltv.Channel{
            {
                ID: "channel1.example.com",
                DisplayNames: []xmltv.DisplayName{
                    {Text: "Example TV"},
                    {Text: "ETV", Lang: makePointer("en")},
                },
                Icons: []xmltv.Icon{
                    {Source: "https://example.com/logo.png"},
                },
            },
        },
        Programmes: []xmltv.Programme{
            {
                Start:   xmltv.Time{Time: now},
                Stop:    &xmltv.Time{Time: now.Add(30 * time.Minute)},
                Channel: "channel1.example.com",
                Titles: []xmltv.Title{
                    {Text: "Sample Program"},
                },
                Descriptions: []xmltv.Description{
                    {Text: "This is a sample program description."},
                },
                Categories: []xmltv.Category{
                    {Text: "Entertainment"},
                },
            },
        },
    }

    // Marshal to XML
    output, err := xml.MarshalIndent(epg, "", "  ")
    if err != nil {
        log.Fatalf("Error marshaling XMLTV: %v\n", err)
    }

    // Add XML header
    result := []byte(xml.Header + string(output))

    // Write to file
    if err := os.WriteFile("output.xml", result, fs.FileMode(0644)); err != nil {
        log.Fatalf("Error writing file: %v\n", err)
    }

    log.Println("XMLTV file created successfully!")
}

```
