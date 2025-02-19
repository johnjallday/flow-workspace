// MusicDetails represents additional metadata for music production projects.
package project

type MusicDetails struct {
	BPM         int      `toml:"bpm,omitempty"`
	Artist      string   `toml:"artist,omitempty"`
	Writers     []string `toml:"writers,omitempty"`
	Genre       string   `toml:"genre,omitempty"`
	Key         string   `toml:"key,omitempty"`
	KeyChange   bool     `toml:"key_change,omitempty"`
	TempoChange bool     `toml:"tempo_change,omitempty"`
}
