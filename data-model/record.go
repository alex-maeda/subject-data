package datamodel

// RecordMetadata describes how a record was obtained.
type RecordMetadata struct {
	Website     *string   `json:"website,omitempty" jsonschema_description:"Website where the record was obtained"`
	URL         *string   `json:"url,omitempty" jsonschema_description:"Full URL where the record was obtained"`
	CaptureDate *DateTime `json:"capture_date,omitempty" jsonschema_description:"When the record was obtained (ISO 8601)"`
	Tags        []string  `json:"tags" jsonschema_description:"Free-form text tags"`
	Identifiers []string  `json:"identifiers" jsonschema_description:"Additional identifiers"`
	ImagePath   *string   `json:"image_path,omitempty" jsonschema_description:"Local path to an image file for this record"`
}

// NewRecordMetadata creates a RecordMetadata with initialized slices.
func NewRecordMetadata() RecordMetadata {
	return RecordMetadata{
		Tags:        []string{},
		Identifiers: []string{},
	}
}

// DataSource describes what type of record this is.
type DataSource struct {
	Platform    string `json:"platform" jsonschema_description:"Source platform for the record"`
	ContentType string `json:"content_type" jsonschema_description:"Classification of the content type"`
}

// MediaAnalysis holds the analysis of an image or video.
type MediaAnalysis struct {
	Description       *string `json:"description,omitempty" jsonschema_description:"Description of the media"`
	Subjects          *string `json:"subjects,omitempty" jsonschema_description:"People in the media"`
	Setting           *string `json:"setting,omitempty" jsonschema_description:"Location of the subjects in the media"`
	ApparentActivity  *string `json:"apparent_activity,omitempty" jsonschema_description:"Activities the subjects are engaged in"`
	PresentationStyle *string `json:"presentation_style,omitempty" jsonschema_description:"How the media is presented"`
	MoodTone          *string `json:"mood_tone,omitempty" jsonschema_description:"Mood of the media"`
	TextInMedia       *string `json:"text_in_media,omitempty" jsonschema_description:"Text features in the media"`
}

// ImageEntry represents an image with its analysis.
type ImageEntry struct {
	URL      string         `json:"url" jsonschema_description:"URL to access the image content"`
	Context  string         `json:"context" jsonschema_description:"How the image was used"`
	Analysis *MediaAnalysis `json:"analysis,omitempty" jsonschema_description:"Extracted information from the image"`
}

// VideoEntry represents a video with its analysis.
type VideoEntry struct {
	URL      string         `json:"url" jsonschema_description:"URL to access the video content"`
	Context  string         `json:"context" jsonschema_description:"How the video was used"`
	Analysis *MediaAnalysis `json:"analysis,omitempty" jsonschema_description:"Extracted information from the video"`
}

// PlatformData holds the content of the record, split into normalized and raw data.
type PlatformData struct {
	Normalized map[string]interface{} `json:"normalized,omitempty" jsonschema_description:"Normalized schema for the content type"`
	Raw        map[string]interface{} `json:"raw,omitempty" jsonschema_description:"Platform and content_type specific data"`
}

// Record is the complete data for a subject from a single source.
type Record struct {
	Identifiable
	SubjectID    *string        `json:"subject_id,omitempty" jsonschema_description:"Subject identifier"`
	Metadata     RecordMetadata `json:"metadata" jsonschema_description:"How the record was obtained"`
	DataSource   DataSource     `json:"data_source" jsonschema_description:"What type of record this is"`
	PlatformData PlatformData   `json:"platform_data" jsonschema_description:"Content of the record"`
	Images       []ImageEntry   `json:"images" jsonschema_description:"Images in the record"`
	Videos       []VideoEntry   `json:"videos" jsonschema_description:"Videos in the record"`
}

// NewRecord creates a Record with initialized slices.
func NewRecord() Record {
	return Record{
		Metadata: NewRecordMetadata(),
		Images:   []ImageEntry{},
		Videos:   []VideoEntry{},
	}
}

// PlatformStr returns the platform as a plain string.
func (s *Record) PlatformStr() string {
	return s.DataSource.Platform
}

// ContentText extracts primary text content from this record.
func (s *Record) ContentText() *string {
	if s.PlatformData.Raw != nil {
		for _, key := range []string{"caption_text", "text", "bio"} {
			if val, ok := s.PlatformData.Raw[key]; ok {
				if str, ok := val.(string); ok && len(str) > 0 {
					return &str
				}
			}
		}
	}
	if s.PlatformData.Normalized != nil {
		if val, ok := s.PlatformData.Normalized["text"]; ok {
			if str, ok := val.(string); ok && len(str) > 0 {
				return &str
			}
		}
	}
	return nil
}
