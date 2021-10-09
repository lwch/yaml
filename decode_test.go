package yaml

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	var ret struct {
		Includes struct {
			Next struct {
				Title string `yaml:"title"`
			} `yaml:"next"`
		} `yaml:"includes"`
	}
	err := Decode("test/main.yaml", &ret)
	if err != nil {
		t.Fatal(err)
	}
	if ret.Includes.Next.Title != "next" {
		t.Fatalf("unexpected title: %s", ret.Includes.Next.Title)
	}
}

func TestRender(t *testing.T) {
	str, err := Render("test/main.yaml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}
