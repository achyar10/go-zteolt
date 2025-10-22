package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// TemplateManager handles all template operations
type TemplateManager struct {
	templates map[string]*template.Template
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() (*TemplateManager, error) {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	// Load all templates
	templates := map[string]string{
		"add-onu":            "templates/add-onu.tmpl",
		"delete-onu":         "templates/delete-onu.tmpl",
		"check-attenuation":  "templates/check-attenuation.tmpl",
		"check-unconfigured": "templates/check-unconfigured.tmpl",
	}

	for name, path := range templates {
		if err := tm.loadTemplate(name, path); err != nil {
			return nil, fmt.Errorf("failed to load template %s: %w", name, err)
		}
	}

	return tm, nil
}

// loadTemplate loads a template from file
func (tm *TemplateManager) loadTemplate(name, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return err
	}

	tm.templates[name] = tmpl
	return nil
}

// RenderTemplate renders a template with the given data
func (tm *TemplateManager) RenderTemplate(templateName string, data interface{}) ([]string, string, error) {
	tmpl, exists := tm.templates[templateName]
	if !exists {
		return nil, "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, "", err
	}

	// Split into individual commands
	lines := strings.Split(buf.String(), "\n")
	commands := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) != "" {
			commands = append(commands, line)
		}
	}

	return commands, buf.String(), nil
}

// GetAvailableTemplates returns list of available template names
func (tm *TemplateManager) GetAvailableTemplates() []string {
	templates := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		templates = append(templates, name)
	}
	return templates
}