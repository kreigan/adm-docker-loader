package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// StackRepository handles stack discovery and retrieval.
type StackRepository struct {
	logger    *Logger
	compose   *ComposeClient
	stacksDir string
}

// NewStackRepository creates a new stack repository.
func NewStackRepository(baseDir string, logger *Logger, compose *ComposeClient) *StackRepository {
	return &StackRepository{
		logger:    logger,
		compose:   compose,
		stacksDir: filepath.Join(baseDir, "stacks"),
	}
}

// FindAll discovers all stacks in the stacks directory.
func (r *StackRepository) FindAll() ([]*Stack, error) {
	if _, err := os.Stat(r.stacksDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("stacks directory does not exist: %s", r.stacksDir)
	}

	entries, err := os.ReadDir(r.stacksDir)
	if err != nil {
		return nil, fmt.Errorf("reading stacks directory: %w", err)
	}

	pattern := regexp.MustCompile(`^\d{2}-.+$`)
	stacks := make([]*Stack, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !pattern.MatchString(name) {
			r.logger.Warning("Skipping directory with invalid name: %s", name)
			continue
		}

		stackName := extractStackName(name)
		if stackName == "" {
			r.logger.Warning("Skipping directory with invalid name pattern: %s", name)
			continue
		}

		stacks = append(stacks, &Stack{
			Name: stackName,
			Dir:  filepath.Join(r.stacksDir, name),
		})
	}

	// Sort stacks lexically by directory name
	sort.Slice(stacks, func(i, j int) bool {
		return filepath.Base(stacks[i].Dir) < filepath.Base(stacks[j].Dir)
	})

	// Populate status for all stacks
	r.populateStatuses(stacks)

	return stacks, nil
}

// FindByName finds a stack by name or directory name.
func (r *StackRepository) FindByName(name string) (*Stack, error) {
	stacks, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, stack := range stacks {
		if stack.Name == name || filepath.Base(stack.Dir) == name {
			return stack, nil
		}
	}

	return nil, fmt.Errorf("stack not found: %s", name)
}

func (r *StackRepository) populateStatuses(stacks []*Stack) {
	statuses := r.compose.GetProjectStatuses()
	for _, stack := range stacks {
		stack.Status = statuses[stack.Name]
		if stack.Status == "" {
			stack.Status = StackStatusDown
		}
	}
}

// extractStackName extracts the stack name from directory name.
// Format: NN-stack-name -> stack-name
func extractStackName(dirName string) string {
	// Find the first dash
	idx := -1
	for i := 0; i < len(dirName); i++ {
		if dirName[i] == '-' {
			idx = i
			break
		}
	}

	if idx == -1 || idx == len(dirName)-1 {
		return ""
	}

	return dirName[idx+1:]
}

// CheckDuplicates returns an error if duplicate stack names are found.
func CheckDuplicates(stacks []*Stack) error {
	nameMap := make(map[string][]string)

	for _, stack := range stacks {
		nameMap[stack.Name] = append(nameMap[stack.Name], filepath.Base(stack.Dir))
	}

	var duplicates []string
	for name, dirs := range nameMap {
		if len(dirs) > 1 {
			sort.Strings(dirs)
			duplicates = append(duplicates, fmt.Sprintf("  - '%s' in: %s", name, joinStrings(dirs, ", ")))
		}
	}

	if len(duplicates) > 0 {
		sort.Strings(duplicates)
		return fmt.Errorf("duplicate stack names found:\n%s", joinStrings(duplicates, "\n"))
	}

	return nil
}

// WarnDuplicates prints a warning if duplicate stack names are found.
func WarnDuplicates(stacks []*Stack) {
	nameMap := make(map[string][]string)

	for _, stack := range stacks {
		nameMap[stack.Name] = append(nameMap[stack.Name], filepath.Base(stack.Dir))
	}

	var duplicates []string
	for name, dirs := range nameMap {
		if len(dirs) > 1 {
			sort.Strings(dirs)
			duplicates = append(duplicates, fmt.Sprintf("  - '%s' in: %s", name, joinStrings(dirs, ", ")))
		}
	}

	if len(duplicates) > 0 {
		sort.Strings(duplicates)
		fmt.Fprintf(os.Stderr, "WARNING: Duplicate stack names found:\n%s\n\n", joinStrings(duplicates, "\n"))
	}
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	n := len(sep) * (len(strs) - 1)
	for _, s := range strs {
		n += len(s)
	}

	result := make([]byte, n)
	pos := copy(result, strs[0])
	for _, s := range strs[1:] {
		pos += copy(result[pos:], sep)
		pos += copy(result[pos:], s)
	}

	return string(result)
}
