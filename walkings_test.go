package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

/* Test helpers (shared with storage_test.go and report_test.go) */

// setStdin replaces the package-level stdin reader for the duration of a test.
func setStdin(t *testing.T, input string) {
	t.Helper()
	old := stdin
	stdin = bufio.NewReader(strings.NewReader(input))
	t.Cleanup(func() { stdin = old })
}

// useTempFiles points LocalFile and BackupFile at a temp dir so tests
// never touch the real database or the backup drive. It also stubs the
// drive mount and clears the global notices before and after the test.
func useTempFiles(t *testing.T) {
	t.Helper()
	dir := t.TempDir()

	oldLocal, oldBackup := LocalFile, BackupFile
	LocalFile = filepath.Join(dir, "local", "walkings.json")
	BackupFile = filepath.Join(dir, "backup", "walkings.json")

	for _, f := range []string{LocalFile, BackupFile} {
		if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Never mount the real backup drive from tests.
	oldMount := mountBackup
	mountBackup = func() error { return nil }

	notices = nil

	t.Cleanup(func() {
		LocalFile, BackupFile = oldLocal, oldBackup
		mountBackup = oldMount
		notices = nil
	})
}

// breakLocalFile makes the local save fail: the parent "dir" is a regular file.
// useTempFiles restores LocalFile afterwards.
func breakLocalFile(t *testing.T) {
	t.Helper()
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	LocalFile = filepath.Join(blocker, "sub", "walkings.json")
}

func sampleWalks() Walks {
	return Walks{Walks: []Walk{
		{Id: 1, Distance: "3.2", Duration: "0:45:00", Calories: "200", Date: "01-10-2026"},
		{Id: 2, Distance: "4.8", Duration: "1:02:10", Calories: "310", Date: "02-10-2026"},
		{Id: 5, Distance: "2.0", Duration: "0:30:00", Calories: "120", Date: "03-10-2026"},
	}}
}

// reload reads the saved file back, to verify what was persisted.
func reload(t *testing.T) Walks {
	t.Helper()
	w := Walks{}
	if err := w.Load(LocalFile); err != nil {
		t.Fatalf("reload: %v", err)
	}
	return w
}

// assertNotices checks the global notices exactly.
func assertNotices(t *testing.T, want ...string) {
	t.Helper()
	if len(notices) != len(want) {
		t.Fatalf("notices = %q, want %q", notices, want)
	}
	for i := range want {
		if notices[i] != want[i] {
			t.Errorf("notice %d = %q, want %q", i, notices[i], want[i])
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

/* Pure helpers */

func TestParseID(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"3", 3, false},
		{" 7 ", 7, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"1.5", 0, true},
		{"", 0, true},
	}
	for _, tt := range tests {
		got, err := parseID(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseID(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("parseID(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestNextID(t *testing.T) {
	empty := Walks{}
	if got := empty.nextID(); got != 1 {
		t.Errorf("empty nextID = %d, want 1", got)
	}

	w := sampleWalks() // highest ID is 5, with gaps
	if got := w.nextID(); got != 6 {
		t.Errorf("nextID = %d, want 6", got)
	}
}

func TestIndexOf(t *testing.T) {
	w := sampleWalks()

	idx, err := w.indexOf(5)
	if err != nil || idx != 2 {
		t.Errorf("indexOf(5) = %d, %v; want 2, nil", idx, err)
	}

	if _, err := w.indexOf(99); err == nil {
		t.Error("indexOf(99) expected error, got nil")
	}
}

func TestValidateNumber(t *testing.T) {
	valid := []string{"4.8", "0.1", "10", "100.25"}
	invalid := []string{"", "0", "-1", "-0.5", "abc", "1.2.3"}

	for _, s := range valid {
		if err := validateNumber(s); err != nil {
			t.Errorf("validateNumber(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range invalid {
		if err := validateNumber(s); err == nil {
			t.Errorf("validateNumber(%q) expected error, got nil", s)
		}
	}
}

func TestValidateInt(t *testing.T) {
	valid := []string{"1", "310", "9999"}
	invalid := []string{"", "0", "-5", "3.5", "abc"}

	for _, s := range valid {
		if err := validateInt(s); err != nil {
			t.Errorf("validateInt(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range invalid {
		if err := validateInt(s); err == nil {
			t.Errorf("validateInt(%q) expected error, got nil", s)
		}
	}
}

func TestValidateDuration(t *testing.T) {
	valid := []string{"1:02:10", "0:45", "0:45:00", "0:00:01", "500:10"}
	invalid := []string{
		"",        // empty
		"45",      // single part
		"1:2:3:4", // too many parts
		"0:00",    // zero
		"0:00:00", // zero
		"1:60",    // minutes out of range
		"1:30:60", // seconds out of range
		"-1:30",   // negative
		"a:30",    // not a number
		"1:xx:10", // not a number
	}

	for _, s := range valid {
		if err := validateDuration(s); err != nil {
			t.Errorf("validateDuration(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range invalid {
		if err := validateDuration(s); err == nil {
			t.Errorf("validateDuration(%q) expected error, got nil", s)
		}
	}
}

/* Prompts */

func TestPromptValid(t *testing.T) {
	t.Run("accepts valid input", func(t *testing.T) {
		setStdin(t, "4.8\n")
		got, err := promptValid("Distance:", "", validateNumber)
		if err != nil || got != "4.8" {
			t.Errorf("got %q, %v; want 4.8, nil", got, err)
		}
	})

	t.Run("converts comma to dot", func(t *testing.T) {
		setStdin(t, "4,8\n")
		got, err := promptValid("Distance:", "", validateNumber)
		if err != nil || got != "4.8" {
			t.Errorf("got %q, %v; want 4.8, nil", got, err)
		}
	})

	t.Run("empty input keeps current value", func(t *testing.T) {
		setStdin(t, "\n")
		got, err := promptValid("Distance:", "3.2", validateNumber)
		if err != nil || got != "3.2" {
			t.Errorf("got %q, %v; want 3.2, nil", got, err)
		}
	})

	t.Run("retries after invalid input", func(t *testing.T) {
		setStdin(t, "abc\n-1\n5.5\n")
		got, err := promptValid("Distance:", "", validateNumber)
		if err != nil || got != "5.5" {
			t.Errorf("got %q, %v; want 5.5, nil", got, err)
		}
	})

	t.Run("empty input with no current value retries", func(t *testing.T) {
		setStdin(t, "\n2.5\n")
		got, err := promptValid("Distance:", "", validateNumber)
		if err != nil || got != "2.5" {
			t.Errorf("got %q, %v; want 2.5, nil", got, err)
		}
	})

	t.Run("EOF returns error", func(t *testing.T) {
		setStdin(t, "")
		if _, err := promptValid("Distance:", "", validateNumber); err == nil {
			t.Error("expected error on EOF, got nil")
		}
	})
}

func TestPromptWalkInput(t *testing.T) {
	t.Run("new walk", func(t *testing.T) {
		setStdin(t, "4,8\n1:02:10\n310\n")
		got, err := promptWalkInput(Walk{})
		if err != nil {
			t.Fatal(err)
		}
		want := Walk{Distance: "4.8", Duration: "1:02:10", Calories: "310"}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("edit keeps unchanged fields", func(t *testing.T) {
		setStdin(t, "\n0:50:00\n\n")
		current := Walk{Id: 2, Distance: "4.8", Duration: "1:02:10", Calories: "310", Date: "02-10-2026"}
		got, err := promptWalkInput(current)
		if err != nil {
			t.Fatal(err)
		}
		want := Walk{Id: 2, Distance: "4.8", Duration: "0:50:00", Calories: "310", Date: "02-10-2026"}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}

/* Add */

func TestAdd(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "4.8\n1:02:10\n310\n")

	w := sampleWalks()
	if err := w.Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(w.Walks) != 4 {
		t.Fatalf("len = %d, want 4", len(w.Walks))
	}
	added := w.Walks[3]
	if added.Id != 6 {
		t.Errorf("Id = %d, want 6", added.Id)
	}
	if added.Distance != "4.8" || added.Duration != "1:02:10" || added.Calories != "310" {
		t.Errorf("unexpected walk: %+v", added)
	}
	if want := time.Now().Format("02-01-2006"); added.Date != want {
		t.Errorf("Date = %q, want %q", added.Date, want)
	}

	assertNotices(t, "Added walk #6")

	// Verify it was persisted
	if got := reload(t); len(got.Walks) != 4 {
		t.Errorf("saved len = %d, want 4", len(got.Walks))
	}
}

func TestAddToEmptyList(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "2.5\n0:40:00\n180\n")

	w := Walks{}
	if err := w.Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(w.Walks) != 1 || w.Walks[0].Id != 1 {
		t.Fatalf("unexpected walks: %+v", w.Walks)
	}
	assertNotices(t, "Added walk #1")
}

func TestAddInvalidThenValid(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "abc\n4.8\n0:00\n1:02:10\n3.5\n310\n")

	w := Walks{}
	if err := w.Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(w.Walks) != 1 || w.Walks[0].Id != 1 {
		t.Fatalf("unexpected walks: %+v", w.Walks)
	}
	if w.Walks[0].Calories != "310" {
		t.Errorf("Calories = %q, want 310", w.Walks[0].Calories)
	}
}

func TestAddIncompleteInputDoesNotAdd(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "4.8\n") // EOF before duration and calories

	w := sampleWalks()
	if err := w.Add(); err == nil {
		t.Error("expected error on incomplete input, got nil")
	}
	if len(w.Walks) != 3 {
		t.Errorf("len = %d, want 3 (nothing added)", len(w.Walks))
	}
	assertNotices(t) // no success notice
	if fileExists(LocalFile) {
		t.Error("nothing should be saved")
	}
}

func TestAddSaveFailureReturnsErrorWithoutNotice(t *testing.T) {
	useTempFiles(t)
	breakLocalFile(t)
	setStdin(t, "4.8\n1:02:10\n310\n")

	w := sampleWalks()
	if err := w.Add(); err == nil {
		t.Fatal("expected save error, got nil")
	}
	assertNotices(t) // must not claim success
}

/* Update */

func TestUpdate(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "5.0\n\n400\n") // change distance and calories, keep duration

	w := sampleWalks()
	if err := w.Update("2"); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got := w.Walks[1]
	if got.Id != 2 || got.Date != "02-10-2026" {
		t.Errorf("Id/Date should be unchanged: %+v", got)
	}
	if got.Distance != "5.0" || got.Duration != "1:02:10" || got.Calories != "400" {
		t.Errorf("unexpected walk after update: %+v", got)
	}

	assertNotices(t, "Updated walk #2")

	if saved := reload(t); saved.Walks[1].Distance != "5.0" {
		t.Errorf("update not persisted: %+v", saved.Walks[1])
	}
}

func TestUpdateNoChanges(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "\n\n\n") // press Enter through every prompt

	w := sampleWalks()
	if err := w.Update("2"); err != nil {
		t.Fatalf("Update: %v", err)
	}

	assertNotices(t, "No changes")

	if fileExists(LocalFile) {
		t.Error("nothing should be saved when nothing changed")
	}
	if fileExists(BackupFile) {
		t.Error("no backup should be written when nothing changed")
	}
	if w.Walks[1] != sampleWalks().Walks[1] {
		t.Errorf("walk should be unchanged: %+v", w.Walks[1])
	}
}

func TestUpdateErrors(t *testing.T) {
	useTempFiles(t)

	w := sampleWalks()
	if err := w.Update("abc"); err == nil {
		t.Error("invalid ID: expected error")
	}
	if err := w.Update("99"); err == nil {
		t.Error("missing ID: expected error")
	}
	if len(w.Walks) != 3 {
		t.Errorf("len = %d, want 3", len(w.Walks))
	}
	assertNotices(t)
}

func TestUpdateIncompleteInputKeepsOriginal(t *testing.T) {
	useTempFiles(t)
	setStdin(t, "9.9\n") // EOF before the remaining prompts

	w := sampleWalks()
	if err := w.Update("2"); err == nil {
		t.Error("expected error on incomplete input, got nil")
	}
	if w.Walks[1] != sampleWalks().Walks[1] {
		t.Errorf("walk should be unchanged: %+v", w.Walks[1])
	}
	assertNotices(t)
}

func TestUpdateSaveFailureReturnsErrorWithoutNotice(t *testing.T) {
	useTempFiles(t)
	breakLocalFile(t)
	setStdin(t, "5.0\n\n400\n")

	w := sampleWalks()
	if err := w.Update("2"); err == nil {
		t.Fatal("expected save error, got nil")
	}
	assertNotices(t)
}

/* Delete */

func TestDelete(t *testing.T) {
	t.Run("confirmed", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "y\n")

		w := sampleWalks()
		if err := w.Delete("2"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if len(w.Walks) != 2 {
			t.Fatalf("len = %d, want 2", len(w.Walks))
		}
		for _, walk := range w.Walks {
			if walk.Id == 2 {
				t.Error("walk 2 should be gone")
			}
		}
		assertNotices(t, "Deleted walk #2")

		if saved := reload(t); len(saved.Walks) != 2 {
			t.Errorf("saved len = %d, want 2", len(saved.Walks))
		}
	})

	t.Run("uppercase Y accepted", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "Y\n")

		w := sampleWalks()
		if err := w.Delete("1"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if len(w.Walks) != 2 {
			t.Errorf("len = %d, want 2", len(w.Walks))
		}
		assertNotices(t, "Deleted walk #1")
	})

	t.Run("declined is a notice, not an error", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "n\n")

		w := sampleWalks()
		if err := w.Delete("2"); err != nil {
			t.Errorf("canceling is not an error, got: %v", err)
		}
		if len(w.Walks) != 3 {
			t.Errorf("len = %d, want 3", len(w.Walks))
		}
		assertNotices(t, "Delete canceled")
		if fileExists(LocalFile) {
			t.Error("nothing should be saved after canceling")
		}
	})

	t.Run("anything but y cancels", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "yes\n")

		w := sampleWalks()
		if err := w.Delete("2"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(w.Walks) != 3 {
			t.Errorf("len = %d, want 3", len(w.Walks))
		}
		assertNotices(t, "Delete canceled")
	})

	t.Run("EOF cancels", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "")

		w := sampleWalks()
		if err := w.Delete("2"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(w.Walks) != 3 {
			t.Errorf("len = %d, want 3", len(w.Walks))
		}
		assertNotices(t, "Delete canceled")
	})

	t.Run("bad or missing ID", func(t *testing.T) {
		useTempFiles(t)

		w := sampleWalks()
		if err := w.Delete("abc"); err == nil {
			t.Error("invalid ID: expected error")
		}
		if err := w.Delete("99"); err == nil {
			t.Error("missing ID: expected error")
		}
		if len(w.Walks) != 3 {
			t.Errorf("len = %d, want 3", len(w.Walks))
		}
		assertNotices(t)
	})

	t.Run("save failure returns error without notice", func(t *testing.T) {
		useTempFiles(t)
		breakLocalFile(t)
		setStdin(t, "y\n")

		w := sampleWalks()
		if err := w.Delete("2"); err == nil {
			t.Fatal("expected save error, got nil")
		}
		assertNotices(t)
	})
}

/* Undo */

func TestUndo(t *testing.T) {
	t.Run("confirmed removes last", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "y\n")

		w := sampleWalks()
		if err := w.Undo(); err != nil {
			t.Fatalf("Undo: %v", err)
		}
		if len(w.Walks) != 2 || w.Walks[1].Id != 2 {
			t.Errorf("unexpected walks after undo: %+v", w.Walks)
		}
		assertNotices(t, "Undone walk #5") // last walk's ID, not its position

		if saved := reload(t); len(saved.Walks) != 2 {
			t.Errorf("saved len = %d, want 2", len(saved.Walks))
		}
	})

	t.Run("declined is a notice, not an error", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "n\n")

		w := sampleWalks()
		if err := w.Undo(); err != nil {
			t.Errorf("canceling is not an error, got: %v", err)
		}
		if len(w.Walks) != 3 {
			t.Errorf("len = %d, want 3", len(w.Walks))
		}
		assertNotices(t, "Undo canceled")
		if fileExists(LocalFile) {
			t.Error("nothing should be saved after canceling")
		}
	})

	t.Run("EOF cancels", func(t *testing.T) {
		useTempFiles(t)
		setStdin(t, "")

		w := sampleWalks()
		if err := w.Undo(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(w.Walks) != 3 {
			t.Errorf("len = %d, want 3", len(w.Walks))
		}
		assertNotices(t, "Undo canceled")
	})

	t.Run("empty list", func(t *testing.T) {
		useTempFiles(t)

		w := Walks{}
		if err := w.Undo(); err != nil {
			t.Errorf("empty list is not an error, got: %v", err)
		}
		assertNotices(t, "No walks to undo")
	})

	t.Run("save failure returns error without notice", func(t *testing.T) {
		useTempFiles(t)
		breakLocalFile(t)
		setStdin(t, "y\n")

		w := sampleWalks()
		if err := w.Undo(); err == nil {
			t.Fatal("expected save error, got nil")
		}
		assertNotices(t)
	})
}

/* Several commands in one run */

func TestNoticesAccumulateInOrder(t *testing.T) {
	useTempFiles(t)

	w := sampleWalks()

	setStdin(t, "4.8\n1:02:10\n310\n")
	if err := w.Add(); err != nil {
		t.Fatal(err)
	}
	setStdin(t, "y\n")
	if err := w.Delete("1"); err != nil {
		t.Fatal(err)
	}
	setStdin(t, "y\n")
	if err := w.Undo(); err != nil {
		t.Fatal(err)
	}

	assertNotices(t, "Added walk #6", "Deleted walk #1", "Undone walk #6")
}

func TestCommandsSurviveReloadRoundTrip(t *testing.T) {
	useTempFiles(t)

	w := Walks{}
	setStdin(t, "4.8\n1:02:10\n310\n")
	if err := w.Add(); err != nil {
		t.Fatal(err)
	}

	// A fresh process would Load and continue with the next ID.
	again := reload(t)
	setStdin(t, "2.0\n0:30:00\n120\n")
	if err := again.Add(); err != nil {
		t.Fatal(err)
	}

	if len(again.Walks) != 2 || again.Walks[1].Id != 2 {
		t.Errorf("unexpected walks after reload: %+v", again.Walks)
	}
}
