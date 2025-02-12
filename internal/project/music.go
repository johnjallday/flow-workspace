// MusicDetails represents additional metadata for music production projects.
package project

type MusicDetails struct {
	BPM     int      `toml:"bpm,omitempty"`
	Artist  string   `toml:"artist,omitempty"`
	Writers []string `toml:"writers,omitempty"`
}
