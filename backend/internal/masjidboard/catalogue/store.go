package catalogue

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// Store persists the last-known-good MasjidBoard catalogue partitions. The
// catalogue is disk-first: Load reads it on demand rather than retaining it in
// memory for the lifetime of the process. Save avoids rewriting identical
// state and replaces changed state atomically.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store { return &Store{path: path} }

func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state = normalizeState(state)
	if err := validateState(state); err != nil { return err }
	current, exists, err := s.loadExistingLocked()
	if err != nil { return err }
	if exists && reflect.DeepEqual(current, state) { return nil }
	return s.writeLocked(state)
}

// SavePartition replaces only one location partition while preserving all
// other active last-known-good partitions.
func (s *Store) SavePartition(partition Partition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	partition.Location = partition.Location.Normalized()
	if err := validatePartition(partition); err != nil { return err }
	state, exists, err := s.loadExistingLocked()
	if err != nil { return err }
	current := cloneState(state)
	key := partition.Location.key()
	replaced := false
	for i := range state.Partitions {
		if state.Partitions[i].Location.key() == key {
			state.Partitions[i] = clonePartition(partition)
			replaced = true
			break
		}
	}
	if !replaced { state.Partitions = append(state.Partitions, clonePartition(partition)) }
	state = normalizeState(state)
	if err := validateState(state); err != nil { return err }
	if exists && reflect.DeepEqual(current, state) { return nil }
	return s.writeLocked(state)
}

// RetainLocations removes partitions that no longer belong to the configured
// discovery scope. This prevents boards from removed locations remaining in
// the merged catalogue after a scope change.
func (s *Store) RetainLocations(locations []Location) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	allowed := make(map[string]struct{}, len(locations))
	for _, location := range locations {
		location = location.Normalized()
		if err := location.Validate(); err != nil { return err }
		allowed[location.key()] = struct{}{}
	}
	state, exists, err := s.loadExistingLocked()
	if err != nil { return err }
	if !exists { return nil }
	kept := make([]Partition, 0, len(state.Partitions))
	for _, partition := range state.Partitions {
		if _, ok := allowed[partition.Location.key()]; ok { kept = append(kept, clonePartition(partition)) }
	}
	state.Partitions = kept
	state = normalizeState(state)
	if err := validateState(state); err != nil { return err }
	current, _, _ := s.loadExistingLocked()
	if reflect.DeepEqual(current, state) { return nil }
	return s.writeLocked(state)
}

func (s *Store) loadLocked() (State, error) {
	state, exists, err := s.loadExistingLocked()
	if err != nil { return State{}, err }
	if !exists { return State{}, nil }
	return state, nil
}

func (s *Store) loadExistingLocked() (State, bool, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) { return State{}, false, nil }
	if err != nil { return State{}, false, fmt.Errorf("masjidboard catalogue: read catalogue: %w", err) }
	var state State
	if err := json.Unmarshal(data, &state); err != nil { return State{}, false, fmt.Errorf("masjidboard catalogue: decode catalogue: %w", err) }
	state = normalizeState(state)
	if err := validateState(state); err != nil { return State{}, false, fmt.Errorf("masjidboard catalogue: invalid persisted catalogue: %w", err) }
	return cloneState(state), true, nil
}

func (s *Store) writeLocked(state State) error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil { return fmt.Errorf("masjidboard catalogue: create storage directory: %w", err) }
	data, err := json.Marshal(state)
	if err != nil { return fmt.Errorf("masjidboard catalogue: encode catalogue: %w", err) }
	tmp, err := os.CreateTemp(dir, ".masjidboard-catalogue-*.tmp")
	if err != nil { return fmt.Errorf("masjidboard catalogue: create temporary file: %w", err) }
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0600); err != nil { _ = tmp.Close(); return fmt.Errorf("masjidboard catalogue: set temporary file mode: %w", err) }
	if _, err := tmp.Write(data); err != nil { _ = tmp.Close(); return fmt.Errorf("masjidboard catalogue: write temporary file: %w", err) }
	if err := tmp.Sync(); err != nil { _ = tmp.Close(); return fmt.Errorf("masjidboard catalogue: sync temporary file: %w", err) }
	if err := tmp.Close(); err != nil { return fmt.Errorf("masjidboard catalogue: close temporary file: %w", err) }
	if err := os.Rename(tmpName, s.path); err != nil { return fmt.Errorf("masjidboard catalogue: replace catalogue: %w", err) }
	return nil
}

func validateState(state State) error {
	seenLocations := make(map[string]struct{}, len(state.Partitions))
	for _, partition := range state.Partitions {
		if err := validatePartition(partition); err != nil { return err }
		key := partition.Location.key()
		if _, exists := seenLocations[key]; exists { return fmt.Errorf("masjidboard catalogue: duplicate location partition %q / %q / %q", partition.Location.Country, partition.Location.Region, partition.Location.City) }
		seenLocations[key] = struct{}{}
	}
	return nil
}

func validatePartition(partition Partition) error {
	if err := partition.Location.Validate(); err != nil { return err }
	if partition.RetrievedAt.IsZero() { return fmt.Errorf("masjidboard catalogue: retrieved_at is required for %q", partition.Location.City) }
	if partition.ValidatedAt.IsZero() { return fmt.Errorf("masjidboard catalogue: validated_at is required for %q", partition.Location.City) }
	seen := make(map[string]struct{}, len(partition.Records))
	for _, record := range partition.Records {
		if err := ValidateRecord(record); err != nil { return err }
		if _, exists := seen[record.ID]; exists { return fmt.Errorf("masjidboard catalogue: duplicate record %q in %q", record.ID, partition.Location.City) }
		seen[record.ID] = struct{}{}
	}
	return nil
}

func normalizeState(state State) State {
	out := cloneState(state)
	for i := range out.Partitions { out.Partitions[i].Location = out.Partitions[i].Location.Normalized() }
	sortPartitions(out.Partitions)
	return out
}

func sortPartitions(partitions []Partition) {
	for i := 1; i < len(partitions); i++ {
		for j := i; j > 0 && partitions[j].Location.key() < partitions[j-1].Location.key(); j-- { partitions[j], partitions[j-1] = partitions[j-1], partitions[j] }
	}
}

func cloneState(state State) State {
	copy := state
	if state.Partitions != nil {
		copy.Partitions = make([]Partition, len(state.Partitions))
		for i, partition := range state.Partitions { copy.Partitions[i] = clonePartition(partition) }
	}
	return copy
}

func clonePartition(partition Partition) Partition {
	copy := partition
	if partition.Records != nil {
		copy.Records = make([]Record, len(partition.Records))
		for i, record := range partition.Records { copy.Records[i] = cloneRecord(record) }
	}
	return copy
}
