package mail

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"sync"
)

// Embed the templates directory
//
//go:embed templates/*.html
var templateFS embed.FS

// TemplateManager manages email templates
type TemplateManager struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

// NewTemplateManager creates a new template manager and loads templates
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}
	tm.loadTemplates()
	return tm
}

// loadTemplates loads all email templates from the embedded filesystem
func (tm *TemplateManager) loadTemplates() {
	// Define template functions
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// Read all template files
	entries, err := fs.ReadDir(templateFS, "templates")
	if err != nil {
		log.Printf("Warning: failed to read templates directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Skip partial templates (starting with _)
		if len(name) > 0 && name[0] == '_' {
			continue
		}

		// Load template content
		content, err := templateFS.ReadFile("templates/" + name)
		if err != nil {
			log.Printf("Warning: failed to read template %s: %v", name, err)
			continue
		}

		// Parse template
		tmpl, err := template.New(name).Funcs(funcMap).Parse(string(content))
		if err != nil {
			log.Printf("Warning: failed to parse template %s: %v", name, err)
			continue
		}

		// Store template (without .html extension)
		templateName := name
		if len(templateName) > 5 && templateName[len(templateName)-5:] == ".html" {
			templateName = templateName[:len(templateName)-5]
		}

		tm.mu.Lock()
		tm.templates[templateName] = tmpl
		tm.mu.Unlock()
	}

	log.Printf("Loaded %d email templates", len(tm.templates))
}

// Render renders a template with the given data
func (tm *TemplateManager) Render(templateName string, data map[string]interface{}) (string, error) {
	tm.mu.RLock()
	tmpl, exists := tm.templates[templateName]
	tm.mu.RUnlock()

	if !exists {
		return "", &EmailError{Message: "template not found: " + templateName}
	}

	// Set default values
	if data == nil {
		data = make(map[string]interface{})
	}

	// Add default data
	defaultData := map[string]interface{}{
		"AppName":      "App",
		"AppURL":       "#",
		"SupportEmail": "support@example.com",
	}

	for k, v := range defaultData {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", &EmailError{Message: "failed to render template", Cause: err}
	}

	return buf.String(), nil
}

// AddTemplate adds or updates a template
func (tm *TemplateManager) AddTemplate(name, content string) error {
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(content)
	if err != nil {
		return &EmailError{Message: "failed to parse template", Cause: err}
	}

	tm.mu.Lock()
	tm.templates[name] = tmpl
	tm.mu.Unlock()

	return nil
}

// GetTemplateNames returns all available template names
func (tm *TemplateManager) GetTemplateNames() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}
