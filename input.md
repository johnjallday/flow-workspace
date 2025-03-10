I'm creating a project in go
It's a device manager, which handles connections to multiple computers via ssh

give me name, alias, ProjectType, and tags to identify this project

for project type there's only music, coding or general.

// Project represents a single project's metadata from project_info.toml.
type Project struct {
 Name         string        `toml:"name"`
 Alias        string        `toml:"alias"`
 ProjectType  string        `toml:"project_type"`
 Tags         []string      `toml:"tags"`
 DateCreated  time.Time     `toml:"date_created"`
 DateModified time.Time     `toml:"date_modified"`
 Notes        []string      `toml:"notes"`
 Path         string        `toml:"path"`
 GitURL       string        `toml:"git_url,omitempty"`
 MusicDetails *MusicDetails `toml:"music_details,omitempty"`
}

