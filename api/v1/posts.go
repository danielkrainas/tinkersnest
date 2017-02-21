package v1

type Post struct {
	Name    string     `json:"name" yaml:"name"`
	Title   string     `json:"title" yaml:"title"`
	Publish bool       `json:"publish" yaml:"publish"`
	Author  *Author    `json:"author" yaml:"author"`
	Created int64      `json:"created" yaml:"created"`
	Content []*Content `json:"content" yaml:"content"`
	Tags    []string   `json:"tags" yaml:"tags"`
}

type Author struct {
	Name string `json:"name" yaml:"name"`
	User string `json:"user" yaml:"user"`
}

type ContentType string

var (
	ContentMarkdown ContentType = "markdown"
	ContentText     ContentType = "text"
	ContentHtml     ContentType = "html"
)

type Content struct {
	Type string `json:"type" yaml:"type"`
	Data []byte `json:"data" yaml:"data"`
	Rel  string `json:"rel" yaml:"rel"`
}
