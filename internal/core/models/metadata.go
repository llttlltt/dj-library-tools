package models

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldKind defines if a field is treated as a string or a number in queries.
type FieldKind int

const (
	KindString FieldKind = iota
	KindNumeric
)

// Capability defines a specific feature required to serve or update a field.
type Capability string

const (
	CapNone      Capability = ""
	CapMetadata  Capability = "Metadata"
	CapCues      Capability = "Cues"
	CapBeatgrids Capability = "Beatgrids"
)

// TrackFields is the single source of truth for queryable track metadata.
var TrackFields = map[string]FieldDefinition[Track]{
	"id":         {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string { return t.ID }},
	"title":      {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string { return t.Title }},
	"artist":     {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string { return t.Artist }},
	"album":      {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string { return t.Album }},
	"genre":      {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Genre }},
	"comment":    {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Comment }},
	"label":      {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Label }},
	"year":       {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return strconv.Itoa(t.Year) }},
	"color":      {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Color }},
	"bpm":        {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return fmt.Sprintf("%.2f", t.BPM) }},
	"key":        {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Key }},
	"location":   {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string { return t.Location }},
	"display":    {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string { return t.Display }},
	"rating":     {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return strconv.Itoa(t.Rating) }},
	"plays":      {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return strconv.Itoa(t.Plays) }},
	"added":      {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.DateAdded }},
	"modified":   {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.DateModified }},
	"bitrate":    {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return strconv.Itoa(t.Bitrate) }},
	"samplerate": {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return strconv.Itoa(t.SampleRate) }},
	"size":       {Kind: KindNumeric, RequiredCap: CapMetadata, Accessor: func(t Track) string { return strconv.FormatInt(t.Size, 10) }},
	"remixer":    {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Remixer }},
	"mix":        {Kind: KindString, RequiredCap: CapMetadata, Accessor: func(t Track) string { return t.Mix }},
	"duration":   {Kind: KindNumeric, RequiredCap: CapNone, Accessor: func(t Track) string { return strconv.Itoa(t.Duration) }},
	"playlists": {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string {
		names := make([]string, len(t.Playlists))
		for i, p := range t.Playlists {
			names[i] = p.Name
		}
		return strings.Join(names, ",")
	}},
	"missing": {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string {
		if t.FileExists == nil {
			return ""
		}
		if !*t.FileExists {
			return "true"
		}
		return "false"
	}},
	"exists": {Kind: KindString, RequiredCap: CapNone, Accessor: func(t Track) string {
		if t.FileExists == nil {
			return ""
		}
		if *t.FileExists {
			return "true"
		}
		return "false"
	}},
}

// CollectionFields defines the properties available within specific track collections.
var CollectionFields = map[string]map[string]FieldKind{
	"hotcues": {
		"color":    KindString,
		"name":     KindString,
		"position": KindNumeric,
	},
	"memorycues": {
		"color":    KindString,
		"name":     KindString,
		"position": KindNumeric,
	},
	"beatgrids": {
		"bpm":      KindNumeric,
		"position": KindNumeric,
	},
	"playlists": {
		"name":   KindString,
		"folder": KindString,
	},
}

// CollectionCapabilities maps collections to the required capability to use them.
var CollectionCapabilities = map[string]Capability{
	"hotcues":    CapCues,
	"memorycues": CapCues,
	"beatgrids":  CapBeatgrids,
	"playlists":  CapNone,
}

// GroupFields is the single source of truth for queryable group metadata.
var GroupFields = map[string]FieldDefinition[ResourceGroup]{
	"id":     {Kind: KindString, RequiredCap: CapNone, Accessor: func(g ResourceGroup) string { return g.ID }},
	"name":   {Kind: KindString, RequiredCap: CapNone, Accessor: func(g ResourceGroup) string { return g.Name }},
	"parent": {Kind: KindString, RequiredCap: CapNone, Accessor: func(g ResourceGroup) string { return g.ParentFolder }},
	"folder": {Kind: KindString, RequiredCap: CapNone, Accessor: func(g ResourceGroup) string { return g.ParentFolder }},
	"items":  {Kind: KindNumeric, RequiredCap: CapNone, Accessor: func(g ResourceGroup) string { return strconv.Itoa(g.Items) }},
	"kind":   {Kind: KindString, RequiredCap: CapNone, Accessor: func(g ResourceGroup) string { return string(g.Kind) }},
}

type FieldDefinition[T any] struct {
	Kind        FieldKind
	RequiredCap Capability
	Accessor    func(T) string
}
