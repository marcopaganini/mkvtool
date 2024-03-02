// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"fmt"
	"testing"
)

// getTest5Metadata returns a matroska struct representing testdata/test5.mkv.
func getTest5Metadata() matroska {
	m := matroska{
		FileName: "testdata/test5.mkv",
	}
	m.Tracks = []struct {
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
	}{
		{ID: 0, Type: "video", Properties: struct {
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
		}{Language: "und"}},
		{ID: 1, Type: "audio", Properties: struct {
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
		}{Language: "und"}},
		{ID: 2, Type: "subtitles", Properties: struct {
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
		}{Language: "eng", DefaultTrack: true}},
		{ID: 3, Type: "subtitles", Properties: struct {
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
		}{Language: "hun"}},
		{ID: 4, Type: "subtitles", Properties: struct {
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
		}{Language: "ger"}},
		{ID: 5, Type: "subtitles", Properties: struct {
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
		}{Language: "fre"}},
		{ID: 6, Type: "subtitles", Properties: struct {
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
		}{Language: "spa"}},
		{ID: 7, Type: "subtitles", Properties: struct {
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
		}{Language: "ita"}},
		{ID: 8, Type: "audio", Properties: struct {
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
		}{Language: "eng", TrackName: "Commentary"}},
		{ID: 9, Type: "subtitles", Properties: struct {
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
		}{Language: "jpn"}},
		{ID: 10, Type: "subtitles", Properties: struct {
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
		}{Language: "und"}},
	}
	return m
}

func TestTrackByLanguageAndType(t *testing.T) {
	mkv := getTest5Metadata()

	cases := []struct {
		langs     []string
		tracktype string
		ignore    []string
		want      int
		wantErr   bool
	}{
		{langs: []string{"eng"}, tracktype: typeSubtitles, want: 2},
		{langs: []string{"hun"}, tracktype: typeSubtitles, want: 3},
		{langs: []string{"fre", "eng"}, tracktype: typeSubtitles, want: 5},
		{langs: []string{"ita"}, tracktype: typeSubtitles, want: 7},
		{langs: []string{"jpn"}, tracktype: typeSubtitles, want: 9},
		{langs: []string{"und"}, tracktype: typeSubtitles, want: 10},
		{langs: []string{"eng"}, tracktype: typeAudio, ignore: []string{"Commentary"}, wantErr: true},
		{langs: []string{"und"}, tracktype: typeVideo, want: 0},
		{langs: []string{"non-existent"}, tracktype: typeSubtitles, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v-%s", tc.langs, tc.tracktype), func(t *testing.T) {
			got, err := trackByLanguageAndType(mkv, tc.langs, tc.tracktype, tc.ignore)
			if (err != nil) != tc.wantErr {
				t.Fatalf("trackByLanguageAndType() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Fatalf("trackByLanguageAndType() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	cases := []struct {
		s    string
		slc  []string
		want bool
	}{
		{"foo", []string{"foo"}, true},
		{"foo", []string{"FOO"}, true},
		{"foobar", []string{"foo"}, true},
		{"foo", []string{"foobar"}, false},
		{"foo", []string{}, false},
		{"", []string{"foo"}, false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s-%v", tc.s, tc.slc), func(t *testing.T) {
			if got := stringInSlice(tc.s, tc.slc); got != tc.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tc.want)
			}
		})
	}
}

type mockRunner struct {
	cmds [][]string
}

func (m *mockRunner) run(name string, args ...string) error {
	m.cmds = append(m.cmds, append([]string{name}, args...))
	return nil
}

func TestSetDefault(t *testing.T) {
	mkv := getTest5Metadata()
	runner := &mockRunner{}

	// Test setting track 3 to default.
	// In testdata/test5.mkv, subtitle tracks are 2, 3, 4, 5, 6, 7, 9, 10.
	err := setdefault(mkv, 3, runner)
	if err != nil {
		t.Fatalf("setdefault() error = %v", err)
	}

	// Should have:
	// 1 call to mkvpropedit with all tracks reset to flag-default=0
	// 1 call to mkvpropedit (via adddefault) to set flag-default=1 for track 3.
	// Total 2 calls.
	if len(runner.cmds) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(runner.cmds))
	}

	// Verify the last command is setting track 4 (3+1) to 1.
	last := runner.cmds[len(runner.cmds)-1]
	expected := []string{"mkvpropedit", "testdata/test5.mkv", "--edit", "track:4", "--set", "flag-default=1"}
	if fmt.Sprintf("%v", last) != fmt.Sprintf("%v", expected) {
		t.Errorf("last command = %v, want %v", last, expected)
	}

	// Verify the first command contains reset flags for subtitle tracks.
	first := runner.cmds[0]
	// Should have: mkvpropedit testdata/test5.mkv --edit track:3 --set flag-default=0 --edit track:4 ...
	// 2 + 8*4 = 34 arguments total.
	if len(first) != 34 {
		t.Errorf("first command should have 34 args, got %d: %v", len(first), first)
	}
}

func TestRemux(t *testing.T) {
	runner := &mockRunner{}
	infiles := []string{"in1.mkv", "in2.mkv"}
	outfile := "out.mkv"

	// Test remux with subs
	err := remux(infiles, outfile, runner, true)
	if err != nil {
		t.Fatalf("remux(subs=true) error = %v", err)
	}
	expected := []string{"mkvmerge", "in1.mkv", "in2.mkv", "-o", "out.mkv"}
	if fmt.Sprintf("%v", runner.cmds[0]) != fmt.Sprintf("%v", expected) {
		t.Errorf("remux(subs=true) command = %v, want %v", runner.cmds[0], expected)
	}

	// Test remux without subs
	runner.cmds = nil
	err = remux(infiles, outfile, runner, false)
	if err != nil {
		t.Fatalf("remux(subs=false) error = %v", err)
	}
	expected = []string{"mkvmerge", "-S", "in1.mkv", "in2.mkv", "-o", "out.mkv"}
	if fmt.Sprintf("%v", runner.cmds[0]) != fmt.Sprintf("%v", expected) {
		t.Errorf("remux(subs=false) command = %v, want %v", runner.cmds[0], expected)
	}
}

func TestCheckTrackType(t *testing.T) {
	cases := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"a", "audio", false},
		{"v", "video", false},
		{"s", "subtitles", false},
		{"aud", "audio", false},
		{"vid", "video", false},
		{"sub", "subtitles", false},
		{"audio", "audio", false},
		{"video", "video", false},
		{"subtitles", "subtitles", false},
		{"bad", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := checkTrackType(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("checkTrackType() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("checkTrackType() got = %v, want %v", got, tc.want)
			}
		})
	}
}
