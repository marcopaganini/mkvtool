// This file is part of mkvtool (http://github.com/marcopaganini/mkvtool))
// See instructions in the README.md file that accompanies this program.
// (C) 2022-2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestErrorFromSlice(t *testing.T) {
	cases := []struct {
		input []string
		want  string
	}{
		{input: []string{}, want: ""},
		{input: []string{"error1"}, want: "error1"},
		{input: []string{"error1", "error2"}, want: "error1\nerror2"},
	}

	for _, tc := range cases {
		err := errorFromSlice(tc.input)
		if tc.want == "" {
			if err != nil {
				t.Errorf("errorFromSlice(%v) = %v, want nil", tc.input, err)
			}
		} else {
			if err == nil || err.Error() != tc.want {
				t.Errorf("errorFromSlice(%v) = %v, want %q", tc.input, err, tc.want)
			}
		}
	}
}

func TestCheckTwoArgs(t *testing.T) {
	cases := []struct {
		args    []string
		wantErr bool
	}{
		{args: []string{}, wantErr: true},
		{args: []string{"one"}, wantErr: true},
		{args: []string{"one", "two"}, wantErr: false},
		{args: []string{"one", "two", "three"}, wantErr: true},
	}

	for _, tc := range cases {
		// We'll use a more standard way by running a dummy app.
		testApp := &cli.App{
			Commands: []*cli.Command{
				{
					Name: "test",
					Action: func(c *cli.Context) error {
						err := checkTwoArgs(c)
						if (err != nil) != tc.wantErr {
							t.Errorf("checkTwoArgs(%v) error = %v, wantErr %v", tc.args, err, tc.wantErr)
						}
						return nil
					},
				},
			},
		}
		_ = testApp.Run(append([]string{"mkvtool", "test"}, tc.args...))
	}
}

func TestCheckMultiArgs(t *testing.T) {
	cases := []struct {
		args    []string
		wantErr bool
	}{
		{args: []string{}, wantErr: true},
		{args: []string{"one"}, wantErr: false},
		{args: []string{"one", "two"}, wantErr: false},
	}

	for _, tc := range cases {
		testApp := &cli.App{
			Commands: []*cli.Command{
				{
					Name: "test",
					Action: func(c *cli.Context) error {
						err := checkMultiArgs(c)
						if (err != nil) != tc.wantErr {
							t.Errorf("checkMultiArgs(%v) error = %v, wantErr %v", tc.args, err, tc.wantErr)
						}
						return nil
					},
				},
			},
		}
		_ = testApp.Run(append([]string{"mkvtool", "test"}, tc.args...))
	}
}

func TestRunnerFromContext(t *testing.T) {
	var r runner = &mockRunner{}

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), runnerKey, &r)
		got := runnerFromContext(ctx)
		if got != &r {
			t.Errorf("runnerFromContext() = %v, want %v", got, &r)
		}
	})

	t.Run("Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("runnerFromContext() did not panic")
			}
		}()
		runnerFromContext(context.Background())
	})
}

func TestActionMerge(t *testing.T) {
	// Mock remux? No, actionMerge calls remux (which is in mkvtool.go).
	// remux is already tested in tracks_test.go, but here we test the action wrapper.
	mRunner := &mockRunner{}
	var r runner = mRunner
	ctx := context.WithValue(context.Background(), runnerKey, &r)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "merge",
				Action: actionMerge,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}},
					&cli.BoolFlag{Name: "subs", Value: true},
				},
			},
		},
	}

	err := app.RunContext(ctx, []string{"mkvtool", "merge", "-o", "out.mkv", "in1.mkv", "in2.mkv"})
	if err != nil {
		t.Fatalf("actionMerge() error = %v", err)
	}

	expected := []string{"mkvmerge", "in1.mkv", "in2.mkv", "-o", "out.mkv"}
	if !reflect.DeepEqual(mRunner.cmds[0], expected) {
		t.Errorf("command = %v, want %v", mRunner.cmds[0], expected)
	}
}

func TestActionRemux(t *testing.T) {
	// Save the original and restore it after the test.
	oldMustParseFile := mustParseFile
	defer func() { mustParseFile = oldMustParseFile }()

	// Mock metadata for testing.
	mockMetadata := matroska{}
	mockMetadata.Tracks = []struct {
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
		}{Language: "eng"}},
		{ID: 2, Type: "audio", Properties: struct {
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
		{ID: 3, Type: "audio", Properties: struct {
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
		}{Language: ""}},
		{ID: 4, Type: "audio", Properties: struct {
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
		}{Language: "eng"}},
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
		}{Language: "rus"}},
	}

	mustParseFile = func(fname string) matroska {
		return mockMetadata
	}

	cases := []struct {
		name    string
		args    []string
		flags   map[string]interface{}
		wantCmd []string
		wantErr bool
	}{
		{
			name:    "Error: missing args",
			args:    []string{"input.mkv"},
			wantErr: true,
		},
		{
			name: "Default: keep all",
			args: []string{"input.mkv", "output.mkv"},
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-d", "0", "-a", "1,2,3,4", "-s", "5,6", "input.mkv"},
		},
		{
			name: "Tracks: 1,2,5",
			args: []string{"input.mkv", "output.mkv"},
			flags: map[string]interface{}{
				"lang":   "v:none,a:none,s:none",
				"tracks": "1,2,5",
			},
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-D", "-a", "1,2", "-s", "5", "input.mkv"},
		},
		{
			name: "Lang example: v:all,a:eng,!s:rus",
			args: []string{"input.mkv", "output.mkv"},
			flags: map[string]interface{}{
				"lang": "v:all,a:eng,!s:rus",
			},
			// v:all (0), a:eng (1, 4), !s:rus (all - rus = 5)
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-d", "0", "-a", "1,4", "-s", "5", "input.mkv"},
		},
		{
			name: "Lang special: a:first,s:none",
			args: []string{"input.mkv", "output.mkv"},
			flags: map[string]interface{}{
				"lang": "a:first,s:none",
			},
			// v:default (0), a:first (1), s:none
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-d", "0", "-a", "1", "-S", "input.mkv"},
		},
		{
			name: "Ignore supersedes: lang=v:all,a:all,s:all ignore=Commentary",
			args: []string{"input.mkv", "output.mkv"},
			flags: map[string]interface{}{
				"ignore": []string{"Commentary"},
			},
			// v:all (0), a:all - 4 (1,2,3), s:all (5,6)
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-d", "0", "-a", "1,2,3", "-s", "5,6", "input.mkv"},
		},
		{
			name: "Lang default: a:default",
			args: []string{"input.mkv", "output.mkv"},
			flags: map[string]interface{}{
				"lang": "a:default",
			},
			// v:default (0), a:default (3), s:default (5,6)
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-d", "0", "-a", "3", "-s", "5,6", "input.mkv"},
		},
		{
			name: "Edge: exclude first, exclude none, exclude all, invalid parts",
			args: []string{"input.mkv", "output.mkv"},
			flags: map[string]interface{}{
				"lang": "v:all,!v:first,a:none,!a:none,s:all,!s:all,invalid,v:invalid",
			},
			// v: all(0) - first(0) = none
			// a: none - none = none
			// s: all(5,6) - all(5,6) = none
			wantCmd: []string{"mkvmerge", "-o", "output.mkv", "-D", "-A", "-S", "input.mkv"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mRunner := &mockRunner{}
			app := &cli.App{
				Commands: []*cli.Command{
					{
						Name:   "remux",
						Action: actionRemux,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "lang"},
							&cli.StringFlag{Name: "tracks"},
							&cli.StringSliceFlag{Name: "ignore"},
						},
					},
				},
			}

			// Prepare context and flags.
			var r runner = mRunner
			ctx := context.WithValue(context.Background(), runnerKey, &r)

			fullArgs := []string{"mkvtool", "remux"}
			for k, v := range tc.flags {
				switch val := v.(type) {
				case string:
					fullArgs = append(fullArgs, "--"+k, val)
				case bool:
					if val {
						fullArgs = append(fullArgs, "--"+k)
					}
				case []string:
					for _, s := range val {
						fullArgs = append(fullArgs, "--"+k, s)
					}
				}
			}
			fullArgs = append(fullArgs, tc.args...)

			err := app.RunContext(ctx, fullArgs)

			if (err != nil) != tc.wantErr {
				t.Fatalf("actionRemux() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				if len(mRunner.cmds) != 1 {
					t.Fatalf("expected 1 command, got %d", len(mRunner.cmds))
				}
				gotCmd := mRunner.cmds[0]
				if !reflect.DeepEqual(gotCmd, tc.wantCmd) {
					t.Errorf("command = %v, want %v", gotCmd, tc.wantCmd)
				}
			}
		})
	}
}

func TestActionPrint(t *testing.T) {
	oldFormat := format
	defer func() { format = oldFormat }()

	format = func(mask, fname string) (string, error) {
		if fname == "error.mkv" {
			return "", errors.New("format error")
		}
		return "formatted", nil
	}

	mRunner := &mockRunner{}
	var r runner = mRunner
	ctx := context.WithValue(context.Background(), runnerKey, &r)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "print",
				Action: actionPrint,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format"},
				},
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "print", "in.mkv"})
		if err != nil {
			t.Errorf("actionPrint() error = %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "print", "error.mkv"})
		if err == nil {
			t.Errorf("actionPrint() want error, got nil")
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "print"})
		if err == nil {
			t.Errorf("actionPrint() want error, got nil")
		}
	})
}

func TestActionRename(t *testing.T) {
	oldReadable := readable
	oldRename := rename
	defer func() {
		readable = oldReadable
		rename = oldRename
	}()

	readable = func(fnames []string) []string {
		return fnames
	}
	rename = func(mask, fname string, dryrun bool) error {
		if fname == "error.mkv" {
			return errors.New("rename error")
		}
		return nil
	}

	mRunner := &mockRunner{}
	var r runner = mRunner
	ctx := context.WithValue(context.Background(), runnerKey, &r)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "rename",
				Action: actionRename,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format"},
					&cli.BoolFlag{Name: "dry-run"},
				},
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "rename", "in.mkv"})
		if err != nil {
			t.Errorf("actionRename() error = %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "rename", "error.mkv"})
		if err == nil {
			t.Errorf("actionRename() want error, got nil")
		}
	})
}

func TestActionSetDefault(t *testing.T) {
	oldMustParseFile := mustParseFile
	oldSetdefault := setdefault
	oldReadable := readable
	defer func() {
		mustParseFile = oldMustParseFile
		setdefault = oldSetdefault
		readable = oldReadable
	}()

	readable = func(fnames []string) []string { return fnames }
	mustParseFile = func(fname string) matroska { return matroska{} }
	setdefault = func(mkv matroska, tracknum int, cmd runner) error {
		if tracknum == -1 {
			return errors.New("setdefault error")
		}
		return nil
	}

	mRunner := &mockRunner{}
	var r runner = mRunner
	ctx := context.WithValue(context.Background(), runnerKey, &r)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "setdefault",
				Action: actionSetDefault,
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "track"},
				},
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "setdefault", "--track", "1", "in.mkv"})
		if err != nil {
			t.Errorf("actionSetDefault() error = %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "setdefault", "--track", "-1", "in.mkv"})
		if err == nil {
			t.Errorf("actionSetDefault() want error, got nil")
		}
	})
}

func TestActionSetDefaultByLang(t *testing.T) {
	oldMustParseFile := mustParseFile
	oldTrackByLanguageAndType := trackByLanguageAndType
	oldSetdefault := setdefault
	oldReadable := readable
	defer func() {
		mustParseFile = oldMustParseFile
		trackByLanguageAndType = oldTrackByLanguageAndType
		setdefault = oldSetdefault
		readable = oldReadable
	}()

	readable = func(fnames []string) []string { return fnames }
	mustParseFile = func(fname string) matroska { return matroska{} }
	trackByLanguageAndType = func(mkv matroska, languages []string, tracktype string, ignore []string) (int, error) {
		if languages[0] == "error" {
			return 0, errors.New("track error")
		}
		return 1, nil
	}
	setdefault = func(mkv matroska, tracknum int, cmd runner) error {
		if tracknum == 2 {
			return errors.New("setdefault error")
		}
		return nil
	}

	mRunner := &mockRunner{}
	var r runner = mRunner
	ctx := context.WithValue(context.Background(), runnerKey, &r)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "setdefaultbylang",
				Action: actionSetDefaultByLang,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "lang"},
					&cli.StringSliceFlag{Name: "ignore"},
				},
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "setdefaultbylang", "--lang", "eng", "in.mkv"})
		if err != nil {
			t.Errorf("actionSetDefaultByLang() error = %v", err)
		}
	})

	t.Run("TrackError", func(t *testing.T) {
		err := app.RunContext(ctx, []string{"mkvtool", "setdefaultbylang", "--lang", "error", "in.mkv"})
		if err == nil {
			t.Errorf("actionSetDefaultByLang() want error, got nil")
		}
	})

	t.Run("SetDefaultError", func(t *testing.T) {
		trackByLanguageAndType = func(mkv matroska, languages []string, tracktype string, ignore []string) (int, error) {
			return 2, nil
		}
		err := app.RunContext(ctx, []string{"mkvtool", "setdefaultbylang", "--lang", "eng", "in.mkv"})
		if err == nil {
			t.Errorf("actionSetDefaultByLang() want error, got nil")
		}
	})
}

func TestActionShow(t *testing.T) {
	oldMustParseFile := mustParseFile
	oldShow := show
	oldReadable := readable
	defer func() {
		mustParseFile = oldMustParseFile
		show = oldShow
		readable = oldReadable
	}()

	readable = func(fnames []string) []string { return fnames }
	mustParseFile = func(fname string) matroska { return matroska{} }
	showCalled := false
	show = func(mkv matroska, showUID bool) {
		showCalled = true
	}

	mRunner := &mockRunner{}
	var r runner = mRunner
	ctx := context.WithValue(context.Background(), runnerKey, &r)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "show",
				Action: actionShow,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "uid"},
				},
			},
		},
	}

	err := app.RunContext(ctx, []string{"mkvtool", "show", "in.mkv"})
	if err != nil {
		t.Errorf("actionShow() error = %v", err)
	}
	if !showCalled {
		t.Errorf("show() was not called")
	}
}
