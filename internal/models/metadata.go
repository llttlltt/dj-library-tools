package models

import (
	"fmt"
	"strconv"
)

// FieldKind defines if a field is treated as a string or a number in queries.
type FieldKind int

const (
	KindString FieldKind = iota
	KindNumeric
)

// TrackFields is the single source of truth for queryable track metadata.
var TrackFields = map[string]FieldDefinition[Track]{
	"id":         {Kind: KindString, Accessor: func(t Track) string { return t.ID }},
	"title":      {Kind: KindString, Accessor: func(t Track) string { return t.Title }},
	"artist":     {Kind: KindString, Accessor: func(t Track) string { return t.Artist }},
	"album":      {Kind: KindString, Accessor: func(t Track) string { return t.Album }},
	"genre":      {Kind: KindString, Accessor: func(t Track) string { return t.Genre }},
	"comment":    {Kind: KindString, Accessor: func(t Track) string { return t.Comment }},
	"label":      {Kind: KindString, Accessor: func(t Track) string { return t.Label }},
	"year":       {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Year) }},
	"color":      {Kind: KindString, Accessor: func(t Track) string { return t.Color }},
	"bpm":        {Kind: KindNumeric, Accessor: func(t Track) string { return fmt.Sprintf("%.2f", t.BPM) }},
	"key":        {Kind: KindString, Accessor: func(t Track) string { return t.Key }},
	"location":   {Kind: KindString, Accessor: func(t Track) string { return t.Location }},
	"display":    {Kind: KindString, Accessor: func(t Track) string { return t.Display }},
	"rating":     {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Rating) }},
	"plays":      {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Plays) }},
	"added":      {Kind: KindString, Accessor: func(t Track) string { return t.DateAdded }},
	"modified":   {Kind: KindString, Accessor: func(t Track) string { return t.DateModified }},
	"bitrate":    {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Bitrate) }},
	"samplerate": {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.SampleRate) }},
	"size":       {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.FormatInt(t.Size, 10) }},
	"remixer":    {Kind: KindString, Accessor: func(t Track) string { return t.Remixer }},
	"mix":        {Kind: KindString, Accessor: func(t Track) string { return t.Mix }},
	"duration":   {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Duration) }},
	"hotcues":    {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Hotcues()) }},
	"memorycues": {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Memorycues()) }},
	"beatgrids":  {Kind: KindNumeric, Accessor: func(t Track) string { return strconv.Itoa(t.Beatgrids()) }},
	"playlists":  {Kind: KindNumeric, Accessor: func(t Track) string { return "0" }}, // Handled specially in evaluator
}

// GroupFields is the single source of truth for queryable group metadata.
var GroupFields = map[string]FieldDefinition[ResourceGroup]{
	"id":     {Kind: KindString, Accessor: func(g ResourceGroup) string { return g.ID }},
	"name":   {Kind: KindString, Accessor: func(g ResourceGroup) string { return g.Name }},
	"parent": {Kind: KindString, Accessor: func(g ResourceGroup) string { return g.ParentFolder }},
	"folder": {Kind: KindString, Accessor: func(g ResourceGroup) string { return g.ParentFolder }},
	"items":  {Kind: KindNumeric, Accessor: func(g ResourceGroup) string { return strconv.Itoa(g.Items) }},
	"kind":   {Kind: KindString, Accessor: func(g ResourceGroup) string { return string(g.Kind) }},
}

type FieldDefinition[T any] struct {
	Kind     FieldKind
	Accessor func(T) string
}
