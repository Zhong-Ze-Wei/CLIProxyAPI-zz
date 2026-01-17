package util

import (
	"github.com/tidwall/sjson"
)

// ApplyImageConfigFromMetadataCLI applies image configuration from metadata to Gemini CLI format requests.
// The configuration is injected into the request.generationConfig.imageConfig path.
//
// Supported configuration:
//   - aspectRatio: Image aspect ratio (e.g., "16:9", "4:3", "1:1")
//   - imageSize: Image resolution (e.g., "1K", "2K", "4K")
//
// Parameters:
//   - model: The model name (currently unused, reserved for future model-specific logic)
//   - metadata: Metadata containing image configuration extracted from model name suffix
//   - rawJSON: The request payload in JSON format
//
// Returns:
//   - []byte: The modified request payload with image configuration applied
func ApplyImageConfigFromMetadataCLI(model string, metadata map[string]any, rawJSON []byte) []byte {
	return applyImageConfig(metadata, rawJSON, "request.generationConfig.imageConfig")
}

// ApplyImageConfigFromMetadata applies image configuration from metadata to Gemini format requests.
// The configuration is injected into the generationConfig.imageConfig path.
//
// Supported configuration:
//   - aspectRatio: Image aspect ratio (e.g., "16:9", "4:3", "1:1")
//   - imageSize: Image resolution (e.g., "1K", "2K", "4K")
//
// Parameters:
//   - model: The model name (currently unused, reserved for future model-specific logic)
//   - metadata: Metadata containing image configuration extracted from model name suffix
//   - rawJSON: The request payload in JSON format
//
// Returns:
//   - []byte: The modified request payload with image configuration applied
func ApplyImageConfigFromMetadata(model string, metadata map[string]any, rawJSON []byte) []byte {
	return applyImageConfig(metadata, rawJSON, "generationConfig.imageConfig")
}

// applyImageConfig is the internal implementation that applies image configuration to the request payload.
// It extracts aspectRatio and imageSize from metadata and injects them into the specified JSON path.
//
// Parameters:
//   - metadata: Metadata containing image configuration
//   - rawJSON: The request payload in JSON format
//   - basePath: The JSON path where imageConfig should be injected
//
// Returns:
//   - []byte: The modified request payload
func applyImageConfig(metadata map[string]any, rawJSON []byte, basePath string) []byte {
	if metadata == nil || len(metadata) == 0 {
		return rawJSON
	}

	// Extract image configuration from metadata
	aspectRatio, imageSize, found := ImageConfigFromMetadata(metadata)

	// If no image configuration found, return original payload
	if !found {
		return rawJSON
	}

	result := rawJSON

	// Inject aspectRatio if present
	if aspectRatio != "" {
		result, _ = sjson.SetBytes(result, basePath+".aspectRatio", aspectRatio)
	}

	// Inject imageSize if present
	if imageSize != "" {
		result, _ = sjson.SetBytes(result, basePath+".imageSize", imageSize)
	}

	return result
}
