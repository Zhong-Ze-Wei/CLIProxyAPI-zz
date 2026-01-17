package util

import (
	"testing"
)

func TestNormalizeImageModel(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantBase         string
		wantAspectRatio  string
		wantImageSize    string
		wantHasMetadata  bool
	}{
		{
			name:            "基础模型 - 无后缀",
			input:           "gemini-3-pro-image",
			wantBase:        "gemini-3-pro-image",
			wantHasMetadata: false,
		},
		{
			name:            "Preview 模型 - 无后缀",
			input:           "gemini-3-pro-image-preview",
			wantBase:        "gemini-3-pro-image-preview",
			wantHasMetadata: false,
		},
		{
			name:            "仅比例后缀 - 16:9 (dash)",
			input:           "gemini-3-pro-image-16-9",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "16:9",
			wantHasMetadata: true,
		},
		{
			name:            "仅比例后缀 - 16:9 (x)",
			input:           "gemini-3-pro-image-16x9",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "16:9",
			wantHasMetadata: true,
		},
		{
			name:            "仅比例后缀 - 9:16",
			input:           "gemini-3-pro-image-9-16",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "9:16",
			wantHasMetadata: true,
		},
		{
			name:            "仅比例后缀 - 21:9",
			input:           "gemini-3-pro-image-21-9",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "21:9",
			wantHasMetadata: true,
		},
		{
			name:            "仅比例后缀 - 4:3",
			input:           "gemini-3-pro-image-4-3",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "4:3",
			wantHasMetadata: true,
		},
		{
			name:            "仅比例后缀 - 1:1",
			input:           "gemini-3-pro-image-1-1",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "1:1",
			wantHasMetadata: true,
		},
		{
			name:            "仅分辨率后缀 - 1K",
			input:           "gemini-3-pro-image-1k",
			wantBase:        "gemini-3-pro-image-preview",
			wantImageSize:   "1K",
			wantHasMetadata: true,
		},
		{
			name:            "仅分辨率后缀 - 2K",
			input:           "gemini-3-pro-image-2k",
			wantBase:        "gemini-3-pro-image-preview",
			wantImageSize:   "2K",
			wantHasMetadata: true,
		},
		{
			name:            "仅分辨率后缀 - 4K",
			input:           "gemini-3-pro-image-4k",
			wantBase:        "gemini-3-pro-image-preview",
			wantImageSize:   "4K",
			wantHasMetadata: true,
		},
		{
			name:            "组合后缀 - 4K + 16:9 (dash)",
			input:           "gemini-3-pro-image-4k-16-9",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "16:9",
			wantImageSize:   "4K",
			wantHasMetadata: true,
		},
		{
			name:            "组合后缀 - 4K + 16:9 (x)",
			input:           "gemini-3-pro-image-4k-16x9",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "16:9",
			wantImageSize:   "4K",
			wantHasMetadata: true,
		},
		{
			name:            "组合后缀 - 16:9 + 4K (反序)",
			input:           "gemini-3-pro-image-16-9-4k",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "16:9",
			wantImageSize:   "4K",
			wantHasMetadata: true,
		},
		{
			name:            "组合后缀 - 2K + 21:9",
			input:           "gemini-3-pro-image-2k-21-9",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "21:9",
			wantImageSize:   "2K",
			wantHasMetadata: true,
		},
		{
			name:            "组合后缀 - 1K + 9:16",
			input:           "gemini-3-pro-image-1k-9x16",
			wantBase:        "gemini-3-pro-image-preview",
			wantAspectRatio: "9:16",
			wantImageSize:   "1K",
			wantHasMetadata: true,
		},
		{
			name:            "非图像模型 - 其他模型",
			input:           "claude-sonnet-4-5",
			wantBase:        "claude-sonnet-4-5",
			wantHasMetadata: false,
		},
		{
			name:            "非图像模型 - Gemini 其他模型",
			input:           "gemini-2.5-flash",
			wantBase:        "gemini-2.5-flash",
			wantHasMetadata: false,
		},
		{
			name:            "无效后缀 - 未知后缀",
			input:           "gemini-3-pro-image-unknown",
			wantBase:        "gemini-3-pro-image-unknown",
			wantHasMetadata: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBase, gotMetadata := NormalizeImageModel(tt.input)

			// 检查基础模型名称
			if gotBase != tt.wantBase {
				t.Errorf("NormalizeImageModel() base = %v, want %v", gotBase, tt.wantBase)
			}

			// 检查是否有 metadata
			hasMetadata := gotMetadata != nil && len(gotMetadata) > 0
			if hasMetadata != tt.wantHasMetadata {
				t.Errorf("NormalizeImageModel() hasMetadata = %v, want %v", hasMetadata, tt.wantHasMetadata)
			}

			if !tt.wantHasMetadata {
				return
			}

			// 检查 aspectRatio
			if tt.wantAspectRatio != "" {
				gotAR, ok := gotMetadata[ImageAspectRatioMetadataKey].(string)
				if !ok {
					t.Errorf("NormalizeImageModel() aspectRatio not found in metadata")
				} else if gotAR != tt.wantAspectRatio {
					t.Errorf("NormalizeImageModel() aspectRatio = %v, want %v", gotAR, tt.wantAspectRatio)
				}
			}

			// 检查 imageSize
			if tt.wantImageSize != "" {
				gotSize, ok := gotMetadata[ImageSizeMetadataKey].(string)
				if !ok {
					t.Errorf("NormalizeImageModel() imageSize not found in metadata")
				} else if gotSize != tt.wantImageSize {
					t.Errorf("NormalizeImageModel() imageSize = %v, want %v", gotSize, tt.wantImageSize)
				}
			}

			// 检查原始模型名称
			gotOriginal, ok := gotMetadata[ImageOriginalModelMetadataKey].(string)
			if !ok {
				t.Errorf("NormalizeImageModel() original model name not found in metadata")
			} else if gotOriginal != tt.input {
				t.Errorf("NormalizeImageModel() original = %v, want %v", gotOriginal, tt.input)
			}
		})
	}
}

func TestImageConfigFromMetadata(t *testing.T) {
	tests := []struct {
		name            string
		metadata        map[string]any
		wantAspectRatio string
		wantImageSize   string
		wantFound       bool
	}{
		{
			name:            "完整配置",
			metadata:        map[string]any{ImageAspectRatioMetadataKey: "16:9", ImageSizeMetadataKey: "4K"},
			wantAspectRatio: "16:9",
			wantImageSize:   "4K",
			wantFound:       true,
		},
		{
			name:            "仅比例",
			metadata:        map[string]any{ImageAspectRatioMetadataKey: "21:9"},
			wantAspectRatio: "21:9",
			wantImageSize:   "",
			wantFound:       true,
		},
		{
			name:          "仅分辨率",
			metadata:      map[string]any{ImageSizeMetadataKey: "2K"},
			wantImageSize: "2K",
			wantFound:     true,
		},
		{
			name:      "空 metadata",
			metadata:  map[string]any{},
			wantFound: false,
		},
		{
			name:      "nil metadata",
			metadata:  nil,
			wantFound: false,
		},
		{
			name:      "无相关字段",
			metadata:  map[string]any{"other_key": "value"},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAR, gotSize, gotFound := ImageConfigFromMetadata(tt.metadata)

			if gotFound != tt.wantFound {
				t.Errorf("ImageConfigFromMetadata() found = %v, want %v", gotFound, tt.wantFound)
			}

			if gotAR != tt.wantAspectRatio {
				t.Errorf("ImageConfigFromMetadata() aspectRatio = %v, want %v", gotAR, tt.wantAspectRatio)
			}

			if gotSize != tt.wantImageSize {
				t.Errorf("ImageConfigFromMetadata() imageSize = %v, want %v", gotSize, tt.wantImageSize)
			}
		})
	}
}
