// Package interpreter implements a tree-walk interpreter for the Aura language.
//
// effect.go defines the effect system infrastructure for Aura.
// Effects are explicit capabilities that functions can request (file I/O, network, time, etc.).
// The effect system enables:
//   - Explicit side effects: functions declare what capabilities they need
//   - Testability: mock providers can replace real I/O in tests
//   - Safety: AI-generated code cannot accidentally introduce side effects
package interpreter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// --- Effect Provider Interfaces ---

// FileProvider defines the interface for file system operations.
// Implementations can be real (os-backed) or mock (in-memory).
type FileProvider interface {
	// ReadFile reads the entire contents of a file.
	ReadFile(path string) (string, error)
	// WriteFile writes content to a file, creating it if necessary.
	WriteFile(path string, content string) error
	// AppendFile appends content to an existing file.
	AppendFile(path string, content string) error
	// Exists checks if a path exists.
	Exists(path string) bool
	// Delete removes a file or empty directory.
	Delete(path string) error
	// ListDir returns the names of entries in a directory.
	ListDir(path string) ([]string, error)
	// CreateDir creates a directory (and parents if needed).
	CreateDir(path string) error
	// IsFile checks if the path is a regular file.
	IsFile(path string) bool
	// IsDir checks if the path is a directory.
	IsDir(path string) bool
}

// --- Effect Context ---

// EffectContext holds all effect capability providers.
// It is threaded through the interpreter and can be swapped for testing.
type EffectContext struct {
	file FileProvider
	// Future providers: net, time, db, log, etc.
}

// NewEffectContext creates a new EffectContext with default (real) providers.
func NewEffectContext() *EffectContext {
	return &EffectContext{
		file: &RealFileProvider{},
	}
}

// NewMockEffectContext creates an EffectContext with mock providers for testing.
func NewMockEffectContext() *EffectContext {
	return &EffectContext{
		file: NewMockFileProvider(),
	}
}

// File returns the file provider.
func (ec *EffectContext) File() FileProvider {
	return ec.file
}

// WithFile returns a new EffectContext with the given file provider.
func (ec *EffectContext) WithFile(fp FileProvider) *EffectContext {
	return &EffectContext{
		file: fp,
	}
}

// --- Real File Provider ---

// RealFileProvider implements FileProvider using the real filesystem.
type RealFileProvider struct{}

func (r *RealFileProvider) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *RealFileProvider) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *RealFileProvider) AppendFile(path string, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func (r *RealFileProvider) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (r *RealFileProvider) Delete(path string) error {
	return os.Remove(path)
}

func (r *RealFileProvider) ListDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

func (r *RealFileProvider) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (r *RealFileProvider) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (r *RealFileProvider) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// --- Mock File Provider ---

// MockFileProvider implements FileProvider with an in-memory filesystem.
// Used for testing to avoid actual filesystem access.
type MockFileProvider struct {
	mu    sync.RWMutex
	files map[string]string // path -> content
	dirs  map[string]bool   // path -> exists
}

// NewMockFileProvider creates a new empty mock filesystem.
func NewMockFileProvider() *MockFileProvider {
	return &MockFileProvider{
		files: make(map[string]string),
		dirs:  make(map[string]bool),
	}
}

// AddFile adds a file to the mock filesystem.
func (m *MockFileProvider) AddFile(path string, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[filepath.Clean(path)] = content
	// Ensure parent directories exist
	dir := filepath.Dir(filepath.Clean(path))
	for dir != "." && dir != "/" {
		m.dirs[dir] = true
		dir = filepath.Dir(dir)
	}
}

// AddDir adds a directory to the mock filesystem.
func (m *MockFileProvider) AddDir(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dirs[filepath.Clean(path)] = true
}

func (m *MockFileProvider) ReadFile(path string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	content, ok := m.files[filepath.Clean(path)]
	if !ok {
		return "", fmt.Errorf("file not found: %s", path)
	}
	return content, nil
}

func (m *MockFileProvider) WriteFile(path string, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cleaned := filepath.Clean(path)
	m.files[cleaned] = content
	// Ensure parent dir exists
	dir := filepath.Dir(cleaned)
	for dir != "." && dir != "/" {
		m.dirs[dir] = true
		dir = filepath.Dir(dir)
	}
	return nil
}

func (m *MockFileProvider) AppendFile(path string, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cleaned := filepath.Clean(path)
	existing, ok := m.files[cleaned]
	if !ok {
		// Create new file if it doesn't exist (matches os behavior with O_CREATE)
		m.files[cleaned] = content
	} else {
		m.files[cleaned] = existing + content
	}
	return nil
}

func (m *MockFileProvider) Exists(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cleaned := filepath.Clean(path)
	if _, ok := m.files[cleaned]; ok {
		return true
	}
	if _, ok := m.dirs[cleaned]; ok {
		return true
	}
	return false
}

func (m *MockFileProvider) Delete(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cleaned := filepath.Clean(path)
	if _, ok := m.files[cleaned]; ok {
		delete(m.files, cleaned)
		return nil
	}
	if _, ok := m.dirs[cleaned]; ok {
		// Check if directory is empty
		prefix := cleaned + string(filepath.Separator)
		for p := range m.files {
			if strings.HasPrefix(p, prefix) {
				return fmt.Errorf("directory not empty: %s", path)
			}
		}
		for p := range m.dirs {
			if p != cleaned && strings.HasPrefix(p, prefix) {
				return fmt.Errorf("directory not empty: %s", path)
			}
		}
		delete(m.dirs, cleaned)
		return nil
	}
	return fmt.Errorf("no such file or directory: %s", path)
}

func (m *MockFileProvider) ListDir(path string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cleaned := filepath.Clean(path)
	if _, ok := m.dirs[cleaned]; !ok {
		// Check if it's root or a path that has children
		hasChildren := false
		prefix := cleaned + string(filepath.Separator)
		for p := range m.files {
			if strings.HasPrefix(p, prefix) {
				hasChildren = true
				break
			}
		}
		if !hasChildren {
			for p := range m.dirs {
				if strings.HasPrefix(p, prefix) {
					hasChildren = true
					break
				}
			}
		}
		if !hasChildren {
			return nil, fmt.Errorf("directory not found: %s", path)
		}
	}

	nameSet := make(map[string]bool)
	prefix := cleaned + string(filepath.Separator)
	for p := range m.files {
		if strings.HasPrefix(p, prefix) {
			rest := p[len(prefix):]
			// Only direct children
			parts := strings.SplitN(rest, string(filepath.Separator), 2)
			nameSet[parts[0]] = true
		}
	}
	for p := range m.dirs {
		if strings.HasPrefix(p, prefix) {
			rest := p[len(prefix):]
			parts := strings.SplitN(rest, string(filepath.Separator), 2)
			nameSet[parts[0]] = true
		}
	}

	names := make([]string, 0, len(nameSet))
	for n := range nameSet {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}

func (m *MockFileProvider) CreateDir(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cleaned := filepath.Clean(path)
	m.dirs[cleaned] = true
	// Create parent directories
	dir := filepath.Dir(cleaned)
	for dir != "." && dir != "/" {
		m.dirs[dir] = true
		dir = filepath.Dir(dir)
	}
	return nil
}

func (m *MockFileProvider) IsFile(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.files[filepath.Clean(path)]
	return ok
}

func (m *MockFileProvider) IsDir(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cleaned := filepath.Clean(path)
	if _, ok := m.dirs[cleaned]; ok {
		return true
	}
	// Check if any files/dirs have this as a prefix (implicit directory)
	prefix := cleaned + string(filepath.Separator)
	for p := range m.files {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}
