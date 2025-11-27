package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type RecipeMetadata struct {
	Description  string   `json:"description,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	Notes        []string `json:"notes,omitempty"`
	MiseEnPlace  []string `json:"mise_en_place,omitempty"`
	Allergens    []string `json:"allergens,omitempty"`
	Equipment    []string `json:"equipment,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	CreatedBy    string   `json:"created_by,omitempty"`
	LastUpdated  string   `json:"last_updated,omitempty"`
}

func LoadMetadata(raw string) (*RecipeMetadata, error) {
	if raw == "" {
		return &RecipeMetadata{}, nil
	}

	var meta RecipeMetadata
	err := json.Unmarshal([]byte(raw), &meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}
	return &meta, nil
}

func MergeMetadata(old *RecipeMetadata, new *RecipeMetadata) *RecipeMetadata {
	if old == nil {
		old = &RecipeMetadata{}
	}
	if new == nil {
		return old
	}

	result := *old

	if new.Description != "" {
		result.Description = new.Description
	}
	if new.Instructions != nil {
		result.Instructions = new.Instructions
	}
	if new.Notes != nil {
		result.Notes = new.Notes
	}
	if new.MiseEnPlace != nil {
		result.MiseEnPlace = new.MiseEnPlace
	}
	if new.Allergens != nil {
		result.Allergens = new.Allergens
	}
	if new.Equipment != nil {
		result.Equipment = new.Equipment
	}
	if new.Tags != nil {
		result.Tags = new.Tags
	}
	if new.CreatedBy != "" {
		result.CreatedBy = new.CreatedBy
	}
	if new.LastUpdated != "" {
		result.LastUpdated = new.LastUpdated
	}

	return &result
}

func MetadataToMarkdown(m *RecipeMetadata) string {
	if m == nil {
		return ""
	}

	var md strings.Builder

	if m.Description != "" {
		md.WriteString("# Description\n\n")
		md.WriteString(m.Description)
		md.WriteString("\n\n")
	}

	if len(m.Instructions) > 0 {
		md.WriteString("# Instructions\n\n")
		for _, instr := range m.Instructions {
			md.WriteString("- ")
			md.WriteString(instr)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}

	if len(m.Notes) > 0 {
		md.WriteString("# Notes\n\n")
		for _, note := range m.Notes {
			md.WriteString("- ")
			md.WriteString(note)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}

	if len(m.MiseEnPlace) > 0 {
		md.WriteString("# Mise En Place\n\n")
		for _, item := range m.MiseEnPlace {
			md.WriteString("- ")
			md.WriteString(item)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}

	if len(m.Allergens) > 0 {
		md.WriteString("# Allergens\n\n")
		for _, allergen := range m.Allergens {
			md.WriteString("- ")
			md.WriteString(allergen)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}

	if len(m.Equipment) > 0 {
		md.WriteString("# Equipment\n\n")
		for _, equip := range m.Equipment {
			md.WriteString("- ")
			md.WriteString(equip)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}

	if len(m.Tags) > 0 {
		md.WriteString("# Tags\n\n")
		for _, tag := range m.Tags {
			md.WriteString("- ")
			md.WriteString(tag)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}

	if m.CreatedBy != "" {
		md.WriteString("# Created By\n\n")
		md.WriteString(m.CreatedBy)
		md.WriteString("\n\n")
	}

	if m.LastUpdated != "" {
		md.WriteString("# Last Updated\n\n")
		md.WriteString(m.LastUpdated)
		md.WriteString("\n\n")
	}

	return md.String()
}

func MarkdownToMetadata(md string) (*RecipeMetadata, error) {
	meta := &RecipeMetadata{}
	scanner := bufio.NewScanner(strings.NewReader(md))

	var currentSection string
	var content strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# ") {
			if currentSection != "" {
				processSection(meta, currentSection, content.String())
				content.Reset()
			}
			currentSection = strings.TrimPrefix(line, "# ")
			currentSection = strings.TrimSpace(currentSection)
			continue
		}

		if currentSection != "" {
			if content.Len() > 0 {
				content.WriteString("\n")
			}
			content.WriteString(line)
		}
	}

	if currentSection != "" {
		processSection(meta, currentSection, content.String())
	}

	return meta, nil
}

func processSection(meta *RecipeMetadata, section, content string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}

	switch strings.ToLower(strings.ReplaceAll(section, " ", "")) {
	case "description":
		meta.Description = content
	case "instructions":
		meta.Instructions = parseList(content)
	case "notes":
		meta.Notes = parseList(content)
	case "miseenplace", "mise_en_place":
		meta.MiseEnPlace = parseList(content)
	case "allergens":
		meta.Allergens = parseList(content)
	case "equipment":
		meta.Equipment = parseList(content)
	case "tags":
		meta.Tags = parseList(content)
	case "createdby", "created_by":
		meta.CreatedBy = content
	case "lastupdated", "last_updated":
		meta.LastUpdated = content
	}
}

func parseList(content string) []string {
	lines := strings.Split(content, "\n")
	var items []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "- ") {
			items = append(items, strings.TrimPrefix(line, "- "))
		} else if strings.HasPrefix(line, "* ") {
			items = append(items, strings.TrimPrefix(line, "* "))
		} else {
			items = append(items, line)
		}
	}

	return items
}

func LoadMetadataFromFile(filepath string) (*RecipeMetadata, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath[strings.LastIndex(filepath, ".")+1:])

	switch ext {
	case "json":
		var meta RecipeMetadata
		err := json.Unmarshal(data, &meta)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		return &meta, nil

	case "md":
		return MarkdownToMetadata(string(data))

	case "txt":
		meta := &RecipeMetadata{}
		content := strings.TrimSpace(string(data))
		if content != "" {
			meta.Notes = []string{content}
		}
		return meta, nil

	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

func SaveMetadataToJSON(meta *RecipeMetadata) (string, error) {
	if meta == nil {
		return "", nil
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return string(data), nil
}

func UpdateTimestamp(meta *RecipeMetadata) {
	if meta != nil {
		meta.LastUpdated = time.Now().Format("2006-01-02 15:04:05")
	}
}
