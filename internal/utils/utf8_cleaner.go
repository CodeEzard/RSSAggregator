package utils

import (
    "regexp"
    "strings"
    "unicode/utf8"
)

// Aggressively clean UTF-8 issues
func sanitizeUTF8(s string) string {
    if s == "" {
        return s
    }
    
    s = strings.ToValidUTF8(s, "")
    
    s = strings.ReplaceAll(s, "\x00", "")
    
    // Remove problematic byte sequences
    s = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`).ReplaceAllString(s, "")
    
    // Fix specific problematic sequences
    s = strings.ReplaceAll(s, "\u00c2", "")  
    s = strings.ReplaceAll(s, "\u00a0", " ") 
    s = strings.ReplaceAll(s, "\u2019", "'") 
    s = strings.ReplaceAll(s, "\u201c", "\"") 
    s = strings.ReplaceAll(s, "\u201d", "\"") 
    s = strings.ReplaceAll(s, "\u2013", "-") 
    s = strings.ReplaceAll(s, "\u2014", "-") 
    
    // Remove any remaining invalid UTF-8
    if !utf8.ValidString(s) {
        // Convert to bytes, filter invalid ones, convert back
        var validBytes []byte
        for _, b := range []byte(s) {
            if b < 128 || utf8.ValidRune(rune(b)) {
                validBytes = append(validBytes, b)
            }
        }
        s = string(validBytes)
    }
    
    return s
}