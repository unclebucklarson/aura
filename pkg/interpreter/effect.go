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
        "time"
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

// --- Time Provider Interface ---

// TimeProvider defines the interface for time-related operations.
// Implementations can be real (os-backed) or mock (controllable for tests).
type TimeProvider interface {
        // Now returns the current time as a Unix timestamp in seconds.
        Now() int64
        // NowNano returns the current time as a Unix timestamp in nanoseconds.
        NowNano() int64
        // Sleep pauses execution for the given number of milliseconds.
        Sleep(ms int)
}

// --- Env Provider Interface ---

// EnvProvider defines the interface for environment variable operations.
// Implementations can be real (os-backed) or mock (in-memory for tests).
type EnvProvider interface {
        // Get returns the value of an environment variable and whether it exists.
        Get(key string) (string, bool)
        // Set sets the value of an environment variable.
        Set(key, value string)
        // Has returns true if the environment variable exists.
        Has(key string) bool
        // List returns all environment variables as a map.
        List() map[string]string
        // Cwd returns the current working directory.
        Cwd() (string, error)
        // Args returns the command-line arguments.
        Args() []string
}

// --- Effect Context ---

// EffectContext holds all effect capability providers.
// It is threaded through the interpreter and can be swapped for testing.
type EffectContext struct {
        file FileProvider
        time TimeProvider
        envp EnvProvider
}

// NewEffectContext creates a new EffectContext with default (real) providers.
func NewEffectContext() *EffectContext {
        return &EffectContext{
                file: &RealFileProvider{},
                time: &RealTimeProvider{},
                envp: &RealEnvProvider{},
        }
}

// NewMockEffectContext creates an EffectContext with mock providers for testing.
func NewMockEffectContext() *EffectContext {
        return &EffectContext{
                file: NewMockFileProvider(),
                time: NewMockTimeProvider(),
                envp: NewMockEnvProvider(),
        }
}

// File returns the file provider.
func (ec *EffectContext) File() FileProvider {
        return ec.file
}

// Time returns the time provider.
func (ec *EffectContext) Time() TimeProvider {
        return ec.time
}

// Env returns the env provider.
func (ec *EffectContext) Env() EnvProvider {
        return ec.envp
}

// WithFile returns a new EffectContext with the given file provider.
func (ec *EffectContext) WithFile(fp FileProvider) *EffectContext {
        return &EffectContext{
                file: fp,
                time: ec.time,
                envp: ec.envp,
        }
}

// WithTime returns a new EffectContext with the given time provider.
func (ec *EffectContext) WithTime(tp TimeProvider) *EffectContext {
        return &EffectContext{
                file: ec.file,
                time: tp,
                envp: ec.envp,
        }
}

// WithEnv returns a new EffectContext with the given env provider.
func (ec *EffectContext) WithEnv(ep EnvProvider) *EffectContext {
        return &EffectContext{
                file: ec.file,
                time: ec.time,
                envp: ep,
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



// --- Real Time Provider ---

// RealTimeProvider implements TimeProvider using Go's time package.
type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() int64 {
        return time.Now().Unix()
}

func (r *RealTimeProvider) NowNano() int64 {
        return time.Now().UnixNano()
}

func (r *RealTimeProvider) Sleep(ms int) {
        time.Sleep(time.Duration(ms) * time.Millisecond)
}

// --- Mock Time Provider ---

// MockTimeProvider implements TimeProvider with controllable time for deterministic tests.
type MockTimeProvider struct {
        mu         sync.RWMutex
        currentSec int64
        currentNs  int64
        sleepLog   []int // records sleep durations in ms
}

// NewMockTimeProvider creates a new mock time provider starting at Unix epoch 1000000.
func NewMockTimeProvider() *MockTimeProvider {
        return &MockTimeProvider{
                currentSec: 1000000,
                currentNs:  1000000 * 1e9,
                sleepLog:   make([]int, 0),
        }
}

// SetTime sets the current time for the mock provider.
func (m *MockTimeProvider) SetTime(sec int64) {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.currentSec = sec
        m.currentNs = sec * 1e9
}

// SleepLog returns the recorded sleep durations.
func (m *MockTimeProvider) SleepLog() []int {
        m.mu.RLock()
        defer m.mu.RUnlock()
        result := make([]int, len(m.sleepLog))
        copy(result, m.sleepLog)
        return result
}

func (m *MockTimeProvider) Now() int64 {
        m.mu.RLock()
        defer m.mu.RUnlock()
        return m.currentSec
}

func (m *MockTimeProvider) NowNano() int64 {
        m.mu.RLock()
        defer m.mu.RUnlock()
        return m.currentNs
}

func (m *MockTimeProvider) Sleep(ms int) {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.sleepLog = append(m.sleepLog, ms)
        // Advance mock time by sleep duration
        m.currentSec += int64(ms) / 1000
        m.currentNs += int64(ms) * 1e6
}

// --- Real Env Provider ---

// RealEnvProvider implements EnvProvider using Go's os package.
type RealEnvProvider struct{}

func (r *RealEnvProvider) Get(key string) (string, bool) {
        return os.LookupEnv(key)
}

func (r *RealEnvProvider) Set(key, value string) {
        os.Setenv(key, value)
}

func (r *RealEnvProvider) Has(key string) bool {
        _, ok := os.LookupEnv(key)
        return ok
}

func (r *RealEnvProvider) List() map[string]string {
        result := make(map[string]string)
        for _, e := range os.Environ() {
                parts := strings.SplitN(e, "=", 2)
                if len(parts) == 2 {
                        result[parts[0]] = parts[1]
                }
        }
        return result
}

func (r *RealEnvProvider) Cwd() (string, error) {
        return os.Getwd()
}

func (r *RealEnvProvider) Args() []string {
        return os.Args
}

// --- Mock Env Provider ---

// MockEnvProvider implements EnvProvider with in-memory environment for testing.
type MockEnvProvider struct {
        mu   sync.RWMutex
        vars map[string]string
        cwd  string
        args []string
}

// NewMockEnvProvider creates a new mock environment provider.
func NewMockEnvProvider() *MockEnvProvider {
        return &MockEnvProvider{
                vars: make(map[string]string),
                cwd:  "/mock/cwd",
                args: []string{"aura"},
        }
}

// SetVar sets a variable in the mock environment.
func (m *MockEnvProvider) SetVar(key, value string) {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.vars[key] = value
}

// SetCwd sets the current working directory in the mock environment.
func (m *MockEnvProvider) SetCwd(cwd string) {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.cwd = cwd
}

// SetArgs sets the command-line arguments in the mock environment.
func (m *MockEnvProvider) SetArgs(args []string) {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.args = make([]string, len(args))
        copy(m.args, args)
}

func (m *MockEnvProvider) Get(key string) (string, bool) {
        m.mu.RLock()
        defer m.mu.RUnlock()
        val, ok := m.vars[key]
        return val, ok
}

func (m *MockEnvProvider) Set(key, value string) {
        m.mu.Lock()
        defer m.mu.Unlock()
        m.vars[key] = value
}

func (m *MockEnvProvider) Has(key string) bool {
        m.mu.RLock()
        defer m.mu.RUnlock()
        _, ok := m.vars[key]
        return ok
}

func (m *MockEnvProvider) List() map[string]string {
        m.mu.RLock()
        defer m.mu.RUnlock()
        result := make(map[string]string, len(m.vars))
        for k, v := range m.vars {
                result[k] = v
        }
        return result
}

func (m *MockEnvProvider) Cwd() (string, error) {
        m.mu.RLock()
        defer m.mu.RUnlock()
        return m.cwd, nil
}

func (m *MockEnvProvider) Args() []string {
        m.mu.RLock()
        defer m.mu.RUnlock()
        result := make([]string, len(m.args))
        copy(result, m.args)
        return result
}

// --- Effect Composition Infrastructure ---

// Clone creates a deep copy of the EffectContext with all the same providers.
func (ec *EffectContext) Clone() *EffectContext {
        return &EffectContext{
                file: ec.file,
                time: ec.time,
                envp: ec.envp,
        }
}

// Derive creates a new EffectContext that inherits from the current one,
// but allows overriding specific providers. nil arguments keep the parent's provider.
func (ec *EffectContext) Derive(file FileProvider, time TimeProvider, env EnvProvider) *EffectContext {
        derived := ec.Clone()
        if file != nil {
                derived.file = file
        }
        if time != nil {
                derived.time = time
        }
        if env != nil {
                derived.envp = env
        }
        return derived
}

// EffectStack manages a stack of EffectContexts, enabling nested effect scopes.
type EffectStack struct {
        mu     sync.RWMutex
        stack  []*EffectContext
}

// NewEffectStack creates a new EffectStack with an initial context.
func NewEffectStack(initial *EffectContext) *EffectStack {
        return &EffectStack{
                stack: []*EffectContext{initial},
        }
}

// Current returns the top-most EffectContext on the stack.
func (es *EffectStack) Current() *EffectContext {
        es.mu.RLock()
        defer es.mu.RUnlock()
        if len(es.stack) == 0 {
                return nil
        }
        return es.stack[len(es.stack)-1]
}

// Push adds a new EffectContext to the stack.
func (es *EffectStack) Push(ctx *EffectContext) {
        es.mu.Lock()
        defer es.mu.Unlock()
        es.stack = append(es.stack, ctx)
}

// Pop removes and returns the top EffectContext from the stack.
// Returns nil if the stack has only one element (the base context is preserved).
func (es *EffectStack) Pop() *EffectContext {
        es.mu.Lock()
        defer es.mu.Unlock()
        if len(es.stack) <= 1 {
                return nil // never pop the base context
        }
        top := es.stack[len(es.stack)-1]
        es.stack = es.stack[:len(es.stack)-1]
        return top
}

// Depth returns the number of contexts on the stack.
func (es *EffectStack) Depth() int {
        es.mu.RLock()
        defer es.mu.RUnlock()
        return len(es.stack)
}

// --- Mock Builder (Fluent API) ---

// MockBuilder provides a fluent API for building mock EffectContexts.
type MockBuilder struct {
        ctx *EffectContext
}

// NewMockBuilder creates a new MockBuilder starting with a fresh mock context.
func NewMockBuilder() *MockBuilder {
        return &MockBuilder{
                ctx: NewMockEffectContext(),
        }
}

// WithFile adds a file to the mock filesystem.
func (mb *MockBuilder) WithFile(path, content string) *MockBuilder {
        if fp, ok := mb.ctx.file.(*MockFileProvider); ok {
                fp.AddFile(path, content)
        }
        return mb
}

// WithDir adds a directory to the mock filesystem.
func (mb *MockBuilder) WithDir(path string) *MockBuilder {
        if fp, ok := mb.ctx.file.(*MockFileProvider); ok {
                fp.AddDir(path)
        }
        return mb
}

// WithTime sets the mock time to a specific Unix timestamp (seconds).
func (mb *MockBuilder) WithTime(sec int64) *MockBuilder {
        if tp, ok := mb.ctx.time.(*MockTimeProvider); ok {
                tp.SetTime(sec)
        }
        return mb
}

// WithEnvVar sets an environment variable in the mock environment.
func (mb *MockBuilder) WithEnvVar(key, value string) *MockBuilder {
        if ep, ok := mb.ctx.envp.(*MockEnvProvider); ok {
                ep.SetVar(key, value)
        }
        return mb
}

// WithCwd sets the current working directory in the mock environment.
func (mb *MockBuilder) WithCwd(cwd string) *MockBuilder {
        if ep, ok := mb.ctx.envp.(*MockEnvProvider); ok {
                ep.SetCwd(cwd)
        }
        return mb
}

// WithArgs sets the command-line arguments in the mock environment.
func (mb *MockBuilder) WithArgs(args []string) *MockBuilder {
        if ep, ok := mb.ctx.envp.(*MockEnvProvider); ok {
                ep.SetArgs(args)
        }
        return mb
}

// WithFiles adds multiple files to the mock filesystem.
func (mb *MockBuilder) WithFiles(files map[string]string) *MockBuilder {
        if fp, ok := mb.ctx.file.(*MockFileProvider); ok {
                for path, content := range files {
                        fp.AddFile(path, content)
                }
        }
        return mb
}

// WithEnvVars sets multiple environment variables in the mock environment.
func (mb *MockBuilder) WithEnvVars(vars map[string]string) *MockBuilder {
        if ep, ok := mb.ctx.envp.(*MockEnvProvider); ok {
                for k, v := range vars {
                        ep.SetVar(k, v)
                }
        }
        return mb
}

// WithFileProvider replaces the file provider entirely.
func (mb *MockBuilder) WithFileProvider(fp FileProvider) *MockBuilder {
        mb.ctx.file = fp
        return mb
}

// WithTimeProvider replaces the time provider entirely.
func (mb *MockBuilder) WithTimeProvider(tp TimeProvider) *MockBuilder {
        mb.ctx.time = tp
        return mb
}

// WithEnvProvider replaces the env provider entirely.
func (mb *MockBuilder) WithEnvProvider(ep EnvProvider) *MockBuilder {
        mb.ctx.envp = ep
        return mb
}

// Build returns the configured EffectContext.
func (mb *MockBuilder) Build() *EffectContext {
        return mb.ctx
}

// --- Pre-configured Fixtures ---

// EmptyMockContext creates a mock context with no files, default time, no env vars.
// This is the same as NewMockEffectContext().
func EmptyMockContext() *EffectContext {
        return NewMockEffectContext()
}

// FixtureWithFiles creates a mock context pre-populated with files.
func FixtureWithFiles(files map[string]string) *EffectContext {
        mb := NewMockBuilder()
        mb.WithFiles(files)
        return mb.Build()
}

// FixtureWithTime creates a mock context with a specific time.
func FixtureWithTime(sec int64) *EffectContext {
        mb := NewMockBuilder()
        mb.WithTime(sec)
        return mb.Build()
}

// FixtureWithEnv creates a mock context with specific environment variables.
func FixtureWithEnv(vars map[string]string) *EffectContext {
        mb := NewMockBuilder()
        mb.WithEnvVars(vars)
        return mb.Build()
}

// FixtureComplete creates a mock context with files, time, and env vars all configured.
func FixtureComplete(files map[string]string, timeSec int64, envVars map[string]string) *EffectContext {
        mb := NewMockBuilder()
        mb.WithFiles(files).WithTime(timeSec).WithEnvVars(envVars)
        return mb.Build()
}

// --- Effect Assertion Helpers (for Go tests) ---

// AssertFileExists checks that a file exists in the mock filesystem.
// Returns true if it exists, false otherwise.
func AssertFileExists(ctx *EffectContext, path string) bool {
        return ctx.File().Exists(path)
}

// AssertFileContent checks that a file has the expected content.
// Returns true if the content matches, false otherwise.
func AssertFileContent(ctx *EffectContext, path string, expected string) bool {
        content, err := ctx.File().ReadFile(path)
        if err != nil {
                return false
        }
        return content == expected
}

// AssertEnvVar checks that an environment variable has the expected value.
func AssertEnvVar(ctx *EffectContext, key string, expected string) bool {
        val, ok := ctx.Env().Get(key)
        if !ok {
                return false
        }
        return val == expected
}

// AssertMockTime checks that the mock time matches the expected value.
func AssertMockTime(ctx *EffectContext, expected int64) bool {
        return ctx.Time().Now() == expected
}

// GetMockFileProvider returns the MockFileProvider from a context, or nil if not mock.
func GetMockFileProvider(ctx *EffectContext) *MockFileProvider {
        if fp, ok := ctx.file.(*MockFileProvider); ok {
                return fp
        }
        return nil
}

// GetMockTimeProvider returns the MockTimeProvider from a context, or nil if not mock.
func GetMockTimeProvider(ctx *EffectContext) *MockTimeProvider {
        if tp, ok := ctx.time.(*MockTimeProvider); ok {
                return tp
        }
        return nil
}

// GetMockEnvProvider returns the MockEnvProvider from a context, or nil if not mock.
func GetMockEnvProvider(ctx *EffectContext) *MockEnvProvider {
        if ep, ok := ctx.envp.(*MockEnvProvider); ok {
                return ep
        }
        return nil
}
