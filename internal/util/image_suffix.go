package util

import (
	"regexp"
	"strings"
)

const (
	ImageAspectRatioMetadataKey   = "image_aspect_ratio"
	ImageSizeMetadataKey          = "image_size"
	ImageOriginalModelMetadataKey = "image_original_model"
)

var (
	// aspectRatioMap maps aspect ratio suffixes to standardized format
	// Supports both dash (-) and x formats
	aspectRatioMap = map[string]string{
		"1-1":  "1:1",
		"1x1":  "1:1",
		"16-9": "16:9",
		"16x9": "16:9",
		"9-16": "9:16",
		"9x16": "9:16",
		"21-9": "21:9",
		"21x9": "21:9",
		"4-3":  "4:3",
		"4x3":  "4:3",
		"3-4":  "3:4",
		"3x4":  "3:4",
		"3-2":  "3:2",
		"3x2":  "3:2",
		"2-3":  "2:3",
		"2x3":  "2:3",
	}

	// imageSizeMap maps image size suffixes to standardized format
	imageSizeMap = map[string]string{
		"1k": "1K",
		"2k": "2K",
		"4k": "4K",
	}
)

// NormalizeImageModel parses image model name suffixes and extracts configuration.
// It handles dynamic suffixes for aspect ratio and image size.
//
// Supported formats:
//   - gemini-3-pro-image-16-9       → aspectRatio: "16:9"
//   - gemini-3-pro-image-4k         → imageSize: "4K"
//   - gemini-3-pro-image-4k-16x9    → aspectRatio: "16:9", imageSize: "4K"
//   - gemini-3-pro-image-21-9-2k    → aspectRatio: "21:9", imageSize: "2K"
//   - gemini-3-pro-image-16x9-4k    → aspectRatio: "16:9", imageSize: "4K"
//
// Parameters:
//   - modelName: The model name to parse
//
// Returns:
//   - string: The normalized base model name (e.g., "gemini-3-pro-image-preview")
//   - map[string]any: Metadata containing extracted configuration (aspectRatio, imageSize)
func NormalizeImageModel(modelName string) (string, map[string]any) {
	if modelName == "" {
		return modelName, nil
	}

	// Only process gemini-3-pro-image models
	if !strings.Contains(modelName, "gemini-3-pro-image") {
		return modelName, nil
	}

	// Base models without suffixes - return as-is
	if modelName == "gemini-3-pro-image" || modelName == "gemini-3-pro-image-preview" {
		return modelName, nil
	}

	// Parse suffix and extract configuration
	baseModel, aspectRatio, imageSize := parseImageSuffix(modelName)

	// If no configuration was extracted, return original model name
	if aspectRatio == "" && imageSize == "" {
		return modelName, nil
	}

	// Build metadata with extracted configuration
	metadata := map[string]any{
		ImageOriginalModelMetadataKey: modelName,
	}

	if aspectRatio != "" {
		metadata[ImageAspectRatioMetadataKey] = aspectRatio
	}
	if imageSize != "" {
		metadata[ImageSizeMetadataKey] = imageSize
	}

	return baseModel, metadata
}

// parseImageSuffix parses the model name suffix to extract aspect ratio and image size.
// Returns the base model name and extracted configuration.
func parseImageSuffix(modelName string) (base, aspectRatio, imageSize string) {
	// Match pattern: gemini-3-pro-image(-preview)?-<suffix>
	pattern := regexp.MustCompile(`^(gemini-3-pro-image(?:-preview)?)-(.+)$`)
	matches := pattern.FindStringSubmatch(modelName)

	if len(matches) < 3 {
		return modelName, "", ""
	}

	base = matches[1]
	suffix := matches[2]

	// Split suffix into parts: e.g., "4k-16-9" → ["4k", "16", "9"]
	parts := strings.Split(suffix, "-")

	// Parse each part to identify aspect ratio and image size
	for i := 0; i < len(parts); i++ {
		part := strings.ToLower(parts[i])

		// Check if this part is an image size (1k, 2k, 4k)
		if size, ok := imageSizeMap[part]; ok {
			imageSize = size
			continue
		}

		// Check if this is part of an aspect ratio (e.g., "16-9")
		if i+1 < len(parts) {
			ratioKey := part + "-" + parts[i+1]
			if ratio, ok := aspectRatioMap[ratioKey]; ok {
				aspectRatio = ratio
				i++ // Skip the next part as it's part of the ratio
				continue
			}
		}

		// Check for "x" format aspect ratio (e.g., "16x9")
		if strings.Contains(part, "x") {
			ratioKey := strings.ReplaceAll(part, "x", "-")
			if ratio, ok := aspectRatioMap[ratioKey]; ok {
				aspectRatio = ratio
			}
		}
	}

	// If nothing was parsed, return original model name
	if aspectRatio == "" && imageSize == "" {
		return modelName, "", ""
	}

	// Keep the base model name as-is (don't force -preview suffix)
	// The upstream may return either gemini-3-pro-image or gemini-3-pro-image-preview

	return base, aspectRatio, imageSize
}

// ImageConfigFromMetadata extracts image configuration from metadata.
// Returns aspect ratio and image size if present.
func ImageConfigFromMetadata(metadata map[string]any) (aspectRatio, imageSize string, found bool) {
	if len(metadata) == 0 {
		return "", "", false
	}

	if ar, ok := metadata[ImageAspectRatioMetadataKey].(string); ok && strings.TrimSpace(ar) != "" {
		aspectRatio = strings.TrimSpace(ar)
		found = true
	}

	if size, ok := metadata[ImageSizeMetadataKey].(string); ok && strings.TrimSpace(size) != "" {
		imageSize = strings.TrimSpace(size)
		found = true
	}

	return aspectRatio, imageSize, found
}
