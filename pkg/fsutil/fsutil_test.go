package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "subdir", "test.txt")
	data := []byte("hello world")

	if err := WriteFile(path, data, 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("expected %q, got %q", data, got)
	}
}

func TestEnsureDir(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "a", "b", "c")
	if err := EnsureDir(path); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestHashFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")
	os.WriteFile(path, []byte("abc"), 0644)

	hash, err := HashFile(path)
	if err != nil {
		t.Fatalf("HashFile failed: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}

	// Same content => same hash
	hash2, _ := HashFile(path)
	if hash != hash2 {
		t.Error("same file should have same hash")
	}
}

func TestHashDir(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("b"), 0644)
	os.Mkdir(filepath.Join(tmp, "sub"), 0755)
	os.WriteFile(filepath.Join(tmp, "sub", "c.txt"), []byte("c"), 0644)

	hashes, err := HashDir(tmp)
	if err != nil {
		t.Fatalf("HashDir failed: %v", err)
	}
	if len(hashes) != 3 {
		t.Errorf("expected 3 hashes, got %d", len(hashes))
	}
}

func TestIsBinary(t *testing.T) {
	tmp := t.TempDir()

	txt := filepath.Join(tmp, "text.txt")
	os.WriteFile(txt, []byte("hello"), 0644)
	bin, err := IsBinary(txt)
	if err != nil {
		t.Fatalf("IsBinary failed: %v", err)
	}
	if bin {
		t.Error("text file should not be binary")
	}

	binFile := filepath.Join(tmp, "binary.bin")
	os.WriteFile(binFile, []byte{0x00, 0x01, 0x02}, 0644)
	bin, err = IsBinary(binFile)
	if err != nil {
		t.Fatalf("IsBinary failed: %v", err)
	}
	if !bin {
		t.Error("file with null byte should be binary")
	}
}

func TestCopyDir(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	os.WriteFile(filepath.Join(src, "file.txt"), []byte("data"), 0644)
	os.Mkdir(filepath.Join(src, "sub"), 0755)
	os.WriteFile(filepath.Join(src, "sub", "nested.txt"), []byte("nested"), 0644)

	if err := CopyDir(src, dst, nil); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, "file.txt")); err != nil {
		t.Error("file.txt not copied")
	}
	if _, err := os.Stat(filepath.Join(dst, "sub", "nested.txt")); err != nil {
		t.Error("nested.txt not copied")
	}
}

func TestFindFiles(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmp, "b.md"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(tmp, "c.txt"), []byte("c"), 0644)

	matches, err := FindFiles(tmp, []string{"*.txt"})
	if err != nil {
		t.Fatalf("FindFiles failed: %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}
