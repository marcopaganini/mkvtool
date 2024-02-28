// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"time"
)

// Source: https://mkvtoolnix.download/doc/mkvmerge-identification-output-schema-v14.json
// Converted to json with: https://json-schema-faker.js.org/ (enable all optionals!)
// Converted to Go Struct with: https://mholt.github.io/json-to-go/
// Numerical UIDs changed to uint64, or they would overflow int/int64 at runtime.
type matroska struct {
	Tracks []struct {
		Codec      string `json:"codec"`
		ID         int    `json:"id"`
		Type       string `json:"type"`
		Properties struct {
			MinimumTimestamp          int    `json:"minimum_timestamp"`
			PixelDimensions           string `json:"pixel_dimensions"`
			UID                       uint64 `json:"uid"`
			CodecPrivateLength        int    `json:"codec_private_length"`
			TagBps                    string `json:"tag_bps"`
			TextSubtitles             bool   `json:"text_subtitles"`
			DefaultTrack              bool   `json:"default_track"`
			EnabledTrack              bool   `json:"enabled_track"`
			CodecDelay                int    `json:"codec_delay"`
			TagFps                    string `json:"tag_fps"`
			Number                    int    `json:"number"`
			SubStreamID               int    `json:"sub_stream_id"`
			FlagHearingImpaired       bool   `json:"flag_hearing_impaired"`
			ProgramNumber             int    `json:"program_number"`
			FlagVisualImpaired        bool   `json:"flag_visual_impaired"`
			CodecPrivateData          string `json:"codec_private_data"`
			CodecID                   string `json:"codec_id"`
			FlagOriginal              bool   `json:"flag_original"`
			TagTitle                  string `json:"tag_title"`
			TagArtist                 string `json:"tag_artist"`
			StreamID                  int    `json:"stream_id"`
			DisplayUnit               int    `json:"display_unit"`
			ContentEncodingAlgorithms string `json:"content_encoding_algorithms"`
			StereoMode                int    `json:"stereo_mode"`
			CodecName                 string `json:"codec_name"`
			AacIsSbr                  string `json:"aac_is_sbr"`
			DisplayDimensions         string `json:"display_dimensions"`
			TeletextPage              int    `json:"teletext_page"`
			DefaultDuration           int    `json:"default_duration"`
			Language                  string `json:"language"`
			TrackName                 string `json:"track_name"`
			MultiplexedTracks         []int  `json:"multiplexed_tracks"`
			FlagCommentary            bool   `json:"flag_commentary"`
			FlagTextDescriptions      bool   `json:"flag_text_descriptions"`
			TagBitsps                 string `json:"tag_bitsps"`
			AudioBitsPerSample        int    `json:"audio_bits_per_sample"`
			AudioChannels             int    `json:"audio_channels"`
			AudioSamplingFrequency    int    `json:"audio_sampling_frequency"`
			Encoding                  string `json:"encoding"`
			ForcedTrack               bool   `json:"forced_track"`
			Packetizer                string `json:"packetizer"`
			LanguageIetf              string `json:"language_ietf"`
		} `json:"properties"`
	} `json:"tracks"`
	IdentificationFormatVersion int    `json:"identification_format_version"`
	FileName                    string `json:"file_name"`
	GlobalTags                  []struct {
		NumEntries int `json:"num_entries"`
	} `json:"global_tags"`
	Container struct {
		Recognized bool `json:"recognized"`
		Supported  bool `json:"supported"`
		Properties struct {
			PreviousSegmentUID string `json:"previous_segment_uid"`
			SegmentUID         string `json:"segment_uid"`
			Playlist           bool   `json:"playlist"`
			Programs           []struct {
				ProgramNumber   int    `json:"program_number"`
				ServiceName     string `json:"service_name"`
				ServiceProvider string `json:"service_provider"`
			} `json:"programs"`
			Duration              int       `json:"duration"`
			PlaylistChapters      int       `json:"playlist_chapters"`
			DateLocal             time.Time `json:"date_local"`
			Title                 string    `json:"title"`
			PlaylistDuration      int       `json:"playlist_duration"`
			MuxingApplication     string    `json:"muxing_application"`
			IsProvidingTimestamps bool      `json:"is_providing_timestamps"`
			PlaylistFile          []string  `json:"playlist_file"`
			ContainerType         int       `json:"container_type"`
			PlaylistSize          int       `json:"playlist_size"`
			NextSegmentUID        string    `json:"next_segment_uid"`
			OtherFile             []string  `json:"other_file"`
			DateUtc               time.Time `json:"date_utc"`
			WritingApplication    string    `json:"writing_application"`
		} `json:"properties"`
		Type string `json:"type"`
	} `json:"container"`
	Warnings    []string `json:"warnings"`
	Attachments []struct {
		FileName   string `json:"file_name"`
		ID         int    `json:"id"`
		Properties struct {
			UID uint64 `json:"uid"`
		} `json:"properties"`
		Size        int    `json:"size"`
		Description string `json:"description"`
		ContentType string `json:"content_type"`
		Type        string `json:"type"`
	} `json:"attachments"`
	Errors    []string `json:"errors"`
	TrackTags []struct {
		NumEntries int `json:"num_entries"`
		TrackID    int `json:"track_id"`
	} `json:"track_tags"`
	Chapters []struct {
		NumEntries int `json:"num_entries"`
	} `json:"chapters"`
}
