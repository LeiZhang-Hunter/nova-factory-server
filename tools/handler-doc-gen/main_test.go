package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWalkGoFilesFollowsSymlinkDirs(t *testing.T) {
	root := t.TempDir()
	appDir := filepath.Join(root, "app")
	targetDir := filepath.Join(root, "shop-target")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "handler.go"), []byte("package shop\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "notes.txt"), []byte("skip"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(targetDir, filepath.Join(appDir, "shop")); err != nil {
		t.Fatal(err)
	}

	var got []string
	if err := walkGoFiles(appDir, func(path string) error {
		got = append(got, path)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	want := []string{filepath.Join(appDir, "shop", "handler.go")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("walkGoFiles() = %#v, want %#v", got, want)
	}
}
