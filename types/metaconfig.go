package types

// MetaConfig is used to determine meta-operations on the csv file, e.g. Rename.
type MetaConfig struct {
	// DestFileNamePattern is used as a string format pattern which takes one
	// string argument: the file name of the source CSV file.
	DestFileNamePattern string `json:"dest_file_name_pattern,omitempty"`
}
