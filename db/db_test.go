package db

import (
	"os"
	"testing"
)

func TestTagCRUD(t *testing.T) {
	tmp := os.TempDir() + "/test_sshconfig.db"
	os.Remove(tmp)
	os.Setenv("HOME", os.TempDir())
	os.WriteFile(tmp, []byte{}, 0600)

	dbPath := getDBPath()
	os.Remove(dbPath)

	host := "testhost"
	tags := []string{"a", "b", "c"}

	if err := AddTags(host, tags); err != nil {
		t.Fatal(err)
	}
	got, err := GetTags(host)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 tags, got %d", len(got))
	}

	if err := RemoveTags(host, []string{"b"}); err != nil {
		t.Fatal(err)
	}
	got, _ = GetTags(host)
	if len(got) != 2 {
		t.Fatalf("want 2 tags after remove, got %d", len(got))
	}

	if err := SetTags(host, []string{"x", "y"}); err != nil {
		t.Fatal(err)
	}
	got, _ = GetTags(host)
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("set tags failed, got %v", got)
	}
}

func TestTagEdgeCases(t *testing.T) {
	host := "edgehost"
	// 空标签
	if err := AddTags(host, []string{""}); err != nil {
		t.Fatal(err)
	}
	got, _ := GetTags(host)
	if len(got) != 0 {
		t.Fatalf("empty tag should not be added")
	}
	// 重复标签
	AddTags(host, []string{"dup", "dup"})
	got, _ = GetTags(host)
	if len(got) != 2 {
		t.Fatalf("should allow duplicate tags, got %v", got)
	}
	// 删除不存在的标签
	if err := RemoveTags(host, []string{"notfound"}); err != nil {
		t.Fatal(err)
	}
	// 覆盖为空
	if err := SetTags(host, []string{}); err != nil {
		t.Fatal(err)
	}
	got, _ = GetTags(host)
	if len(got) != 0 {
		t.Fatalf("all tags should be removed, got %v", got)
	}
}
