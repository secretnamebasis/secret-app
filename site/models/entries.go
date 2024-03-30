package models

type EntriesData struct {
	App     string
	Dev     string
	Entries []BlogEntry
}

type EntryData struct {
	App   string
	Dev   string
	Entry BlogEntry
}
