package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stubMount replaces mountBackup for one test and restores it afterwards.
func stubMount(t *testing.T, fn func() error) {
	t.Helper()
	old := mountBackup
	mountBackup = fn
	t.Cleanup(func() { mountBackup = old })
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

/* Save */

func TestSaveLoadRoundTrip(t *testing.T) {
	useTempFiles(t)

	orig := sampleWalks()
	if err := orig.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded := Walks{}
	if err := loaded.Load(LocalFile); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Walks) != len(orig.Walks) {
		t.Fatalf("len = %d, want %d", len(loaded.Walks), len(orig.Walks))
	}
	for i := range orig.Walks {
		if loaded.Walks[i] != orig.Walks[i] {
			t.Errorf("walk %d = %+v, want %+v", i, loaded.Walks[i], orig.Walks[i])
		}
	}
}

func TestSaveWritesIndentedJSON(t *testing.T) {
	useTempFiles(t)

	w := sampleWalks()
	if err := w.Save(); err != nil {
		t.Fatal(err)
	}

	content := readFile(t, LocalFile)
	if !strings.Contains(content, "\n  \"walks\"") {
		t.Errorf("expected 2-space indented JSON, got:\n%s", content)
	}
}

func TestSaveWritesBackup(t *testing.T) {
	useTempFiles(t)

	mounts := 0
	stubMount(t, func() error { mounts++; return nil })

	w := sampleWalks()
	if err := w.Save(); err != nil {
		t.Fatal(err)
	}

	if mounts != 1 {
		t.Errorf("mountBackup called %d times, want 1", mounts)
	}
	if readFile(t, BackupFile) != readFile(t, LocalFile) {
		t.Error("backup content differs from local content")
	}
}

func TestSaveBackupMountFailureStillSucceeds(t *testing.T) {
	useTempFiles(t)
	stubMount(t, func() error { return errors.New("drive not found") })

	w := sampleWalks()
	if err := w.Save(); err != nil {
		t.Fatalf("Save should succeed when only the backup fails, got: %v", err)
	}

	if _, err := os.Stat(LocalFile); err != nil {
		t.Errorf("local file should exist: %v", err)
	}
	if _, err := os.Stat(BackupFile); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("backup file should not exist, stat err = %v", err)
	}
}

func TestSaveBackupWriteFailureStillSucceeds(t *testing.T) {
	useTempFiles(t)

	// Make the backup location unwritable: its parent "dir" is a regular file.
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	BackupFile = filepath.Join(blocker, "sub", "walkings.json")

	w := sampleWalks()
	if err := w.Save(); err != nil {
		t.Fatalf("Save should succeed when only the backup fails, got: %v", err)
	}
	if _, err := os.Stat(LocalFile); err != nil {
		t.Errorf("local file should exist: %v", err)
	}
}

func TestSaveLocalFailureReturnsError(t *testing.T) {
	useTempFiles(t)

	mounts := 0
	stubMount(t, func() error { mounts++; return nil })

	// Make the local location unwritable: its parent "dir" is a regular file.
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	LocalFile = filepath.Join(blocker, "sub", "walkings.json")

	w := sampleWalks()
	err := w.Save()
	if err == nil {
		t.Fatal("expected error when local save fails, got nil")
	}
	if !strings.Contains(err.Error(), "local save") {
		t.Errorf("error should mention local save, got: %v", err)
	}
	if mounts != 0 {
		t.Errorf("backup should not be attempted after a local failure, mounts = %d", mounts)
	}
}

func TestSaveOverwritesExistingFile(t *testing.T) {
	useTempFiles(t)

	w := sampleWalks()
	if err := w.Save(); err != nil {
		t.Fatal(err)
	}

	w.Walks = w.Walks[:1]
	if err := w.Save(); err != nil {
		t.Fatal(err)
	}

	loaded := Walks{}
	if err := loaded.Load(LocalFile); err != nil {
		t.Fatal(err)
	}
	if len(loaded.Walks) != 1 {
		t.Errorf("len = %d, want 1", len(loaded.Walks))
	}
}

/* Load */

func TestLoadMissingFileIsEmpty(t *testing.T) {
	w := Walks{}
	err := w.Load(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err != nil {
		t.Fatalf("missing file should not be an error, got: %v", err)
	}
	if len(w.Walks) != 0 {
		t.Errorf("len = %d, want 0", len(w.Walks))
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte(`{"walks": [`), 0o600); err != nil {
		t.Fatal(err)
	}

	w := Walks{}
	err := w.Load(path)
	if err == nil {
		t.Fatal("expected error on corrupt JSON, got nil")
	}
	if !strings.Contains(err.Error(), "parsing JSON") {
		t.Errorf("error should mention parsing JSON, got: %v", err)
	}
}

func TestLoadUnreadableSource(t *testing.T) {
	// Reading a directory as a file fails with an error other than ErrNotExist.
	w := Walks{}
	err := w.Load(t.TempDir())
	if err == nil {
		t.Fatal("expected error reading a directory, got nil")
	}
	if !strings.Contains(err.Error(), "reading file") {
		t.Errorf("error should mention reading file, got: %v", err)
	}
}

func TestLoadKnownJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "walkings.json")
	data := `{"walks":[{"id":7,"distance":"4.8","duration":"1:02:10","calories":"310","date":"29-10-2026"}]}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	w := Walks{}
	if err := w.Load(path); err != nil {
		t.Fatal(err)
	}

	want := Walk{Id: 7, Distance: "4.8", Duration: "1:02:10", Calories: "310", Date: "29-10-2026"}
	if len(w.Walks) != 1 || w.Walks[0] != want {
		t.Errorf("got %+v, want [%+v]", w.Walks, want)
	}
}

/* writeFileAtomic */

func TestWriteFileAtomicCreatesParentDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c", "file.json")

	if err := writeFileAtomic(path, []byte("hello"), 0o600); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	if got := readFile(t, path); got != "hello" {
		t.Errorf("content = %q, want hello", got)
	}
}

func TestWriteFileAtomicOverwrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.json")

	if err := writeFileAtomic(path, []byte("old content that is long"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(path, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != "new" {
		t.Errorf("content = %q, want new (no leftovers from old content)", got)
	}
}

func TestWriteFileAtomicSetsPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.json")

	if err := writeFileAtomic(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("perm = %o, want 600", got)
	}
}

func TestWriteFileAtomicLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.json")

	for i := 0; i < 3; i++ {
		if err := writeFileAtomic(path, []byte("data"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "file.json" {
		names := []string{}
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("dir should contain only file.json, got %v", names)
	}
}

func TestWriteFileAtomicFailureLeavesOriginalIntact(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.json")

	if err := writeFileAtomic(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Make the target a non-empty directory so the final rename fails.
	bad := filepath.Join(dir, "target")
	if err := os.MkdirAll(filepath.Join(bad, "inner"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(bad, []byte("new"), 0o600); err == nil {
		t.Error("expected rename onto a non-empty directory to fail")
	}

	// The unrelated original is untouched and no temp file is left behind.
	if got := readFile(t, path); got != "original" {
		t.Errorf("original content changed: %q", got)
	}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}
