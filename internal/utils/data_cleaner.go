package main

import (
    "html"
    "regexp"
    "strings"
)

func cleanJobData(title, description string) (string, string) {
    // Clean title
    cleanTitle := cleanTitle(title)
    
    // Clean description  
    cleanDescription := cleanDescription(description)
    
    return cleanTitle, cleanDescription
}

func cleanTitle(title string) string {
    title = strings.ReplaceAll(title, "\n", "")
    title = strings.ReplaceAll(title, "\t", "")

    title = strings.TrimSpace(title)
    
    return title
}

func cleanDescription(description string) string {
    description = html.UnescapeString(description)
    
    htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
    description = htmlTagRegex.ReplaceAllString(description, "")
    
    spaceRegex := regexp.MustCompile(`\s+`)
    description = spaceRegex.ReplaceAllString(description, " ")
    
    description = strings.TrimSpace(description)
    
    if len(description) > 500 {
        description = description[:500] + "..."
    }
    
    return description
}