#!/bin/bash

echo "Fixing package declarations..."

# Fix all package declarations
find internal -name "*.go" -exec sed -i 's/package main/package handlers/g' {} \; 2>/dev/null
find internal/utils -name "*.go" -exec sed -i 's/package handlers/package utils/g' {} \; 2>/dev/null
find internal/auth -name "*.go" -exec sed -i 's/package handlers/package auth/g' {} \; 2>/dev/null
find internal/models -name "*.go" -exec sed -i 's/package handlers/package models/g' {} \; 2>/dev/null
find internal/scraper -name "*.go" -exec sed -i 's/package handlers/package scraper/g' {} \; 2>/dev/null
find internal/middleware -name "*.go" -exec sed -i 's/package handlers/package middleware/g' {} \; 2>/dev/null

echo "Fixing function names..."
# Fix function names to be exported (capitalize first letter)
find internal -name "*.go" -exec sed -i 's/func handler/func Handler/g' {} \; 2>/dev/null
find internal -name "*.go" -exec sed -i 's/func respond/func Respond/g' {} \; 2>/dev/null

echo "Project fixed!"
