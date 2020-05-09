package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ericaro/frontmatter"
)

const (
	editor        = "mvim"
	fileExtension = "md"
)

func main() {
	// If the -build flag is set, we're not creating a new page, we're rebuilding the index and tag pages
	boolPtr := flag.Bool("build", false, "builds the index and tag pages")
	flag.Parse()
	if *boolPtr {
		pages := loadPages()

		tags := buildTagPages(pages)
		buildIndexPage(pages, tags)

		os.Exit(0)
	}

	// Every non-dash argument is considered a part of the title. If there are no arguments, we have no title
	// Can't have a page without a title
	if len(os.Args[1:]) < 1 {
		fmt.Println("Must have a title")
		os.Exit(1)
	}

	title := strings.Title(strings.Join(os.Args[1:], " "))

	filePath := createNewPage(title)

	// Write the filepath to the console. This makes it easy to know which file we just created
	fmt.Println(filePath)

	// And rebuild the index and tag pages
	pages := loadPages()

	tags := buildTagPages(pages)
	buildIndexPage(pages, tags)

	os.Exit(0)
}

func buildIndexPage(pages []*Page, tags []string) {
	content := "A collection of things\n\n"

	// Loop over the pages in reverse, which puts them in reverse-chronological order
	for _, page := range pages {
		if page.IsContentPage() {
			content += fmt.Sprintf("* %s\n", page.Link())
		}
	}

	content += fmt.Sprintf("\n")

	// Loop over the tags in order and create links to those pages
	sort.Strings(tags)
	for _, tag := range tags {
		content += fmt.Sprintf(
			"[%s](%s), ",
			tag,
			fmt.Sprintf("./%s", tag),
		)
	}

	content += fmt.Sprintf("\n")
	content += fmt.Sprintf("\n")

	content += timestamp()

	// And write the file to disk
	err := ioutil.WriteFile("./docs/index.md", []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// buildTagPages creates the tag pages, with links to posts tagged with those values
func buildTagPages(pages []*Page) []string {
	tags := make(map[string][]*Page)

	// Sort the pages into tag buckets
	for _, page := range pages {
		for _, tag := range page.Tags() {
			tag = strings.TrimSpace(tag)

			if tag != "" {
				tags[tag] = append(tags[tag], page)
			}
		}
	}

	// Then enumerate over the tags in alphabetical order and create a page for each of them
	tagArr := make([]string, len(tags))
	i := 0

	for tag := range tags {
		tagArr[i] = tag
		i++
	}

	for _, tag := range tagArr {
		content := fmt.Sprintf("%s\n\n", tag)

		for _, page := range tags[tag] {
			if page.IsContentPage() {
				content += fmt.Sprintf("* %s\n", page.Link())
			}
		}

		content += fmt.Sprintf("\n")

		content += timestamp()

		// And write the file to disk
		err := ioutil.WriteFile(fmt.Sprintf("./docs/%s.md", tag), []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	return tagArr
}

func createNewPage(title string) string {
	date := time.Now()
	pathDate := date.Format("2006-01-02T15-04-05") // a custom format that plays nicely with GitHub Pages filename restrictions

	filePath := fmt.Sprintf("./docs/%s-%s.%s", pathDate, strings.ReplaceAll(strings.ToLower(title), " ", "-"), fileExtension)

	// Front matter lives at the top of the generated file and contains bits of info about the file
	// This is loosely based on the format Hugo uses
	frontMatter := fmt.Sprintf(
		"---\ndate: %s\ntitle: %s\ntags: %s\n---\n\n",
		date.Format(time.RFC3339),
		title,
		"",
	)

	content := frontMatter + fmt.Sprintf("# %s\n\n\n", title)

	// Write out the stub file, explode if we can't do that
	err := ioutil.WriteFile(fmt.Sprintf("%s", filePath), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// And open the file for editing, exploding if we can't do that
	cmd := exec.Command(editor, filePath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return filePath
}

// loadPages reads the page files from disk (in reverse chronological order) and
// creates Page instances from them
func loadPages() []*Page {
	pages := []*Page{}

	filePaths, _ := filepath.Glob("./docs/*.md")

	for i := len(filePaths) - 1; i >= 0; i-- {
		file := filePaths[i]
		page := readPage(file)

		pages = append(pages, page)
	}

	return pages
}

// readPage reads the contents of the page and unmarshals it into the Page struct,
// making the frontmatter programmatically accessible
func readPage(filePath string) *Page {
	page := new(Page)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = frontmatter.Unmarshal(([]byte)(data), page)
	if err != nil {
		log.Fatal(err)
	}

	page.FilePath = filePath

	return page
}

func timestamp() string {
	return fmt.Sprintf("<sup><sub>generated %s</sub></sup>\n", time.Now().Format("2 Jan 2006 15:04:05"))
}

/* -------------------- Types -------------------- */

// Page represents a TIL page
type Page struct {
	Content  string `fm:"content" yaml:"-"`
	Date     string `yaml:"date"`
	FilePath string `yaml:"filepath"`
	TagsStr  string `yaml:"tags"`
	Title    string `yaml:"title"`
}

// CreatedAt returns a time instance representing when the page was created
func (page *Page) CreatedAt() time.Time {
	date, err := time.Parse(time.RFC3339, page.Date)
	if err != nil {
		log.Fatal(err)
	}

	return date
}

// Link returns a link string suitable for embedding in a Markdown page
func (page *Page) Link() string {
	return fmt.Sprintf(
		" <code>%s</code> [%s](%s)",
		page.PrettyDate(),
		page.Title,
		strings.Replace(page.FilePath, "docs/", "", -1))
}

// IsContentPage returns true if the page is a valid entry page, false if it is not
func (page *Page) IsContentPage() bool {
	return page.Title != ""
}

// PrettyDate returns a human-friendly representation of the CreatedAt
func (page *Page) PrettyDate() string {
	return page.CreatedAt().Format("Jan 02, 2006")
}

// Tags returns a slice of tags assigned to this page
func (page *Page) Tags() []string {
	return strings.Split(page.TagsStr, ",")
}
