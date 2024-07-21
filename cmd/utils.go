package cmd

import "os"

// getEnv fetches the value of an environment variable or returns a fallback value if the variable is not set.
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
