package config

import "os"

func GetGeminiAPIKey() string {
	return os.Getenv("GEMINI_API_KEY")
}