package v1

type Post struct {
	Name    string     `json:"name" yaml:"name"`
	Title   string     `json:"title" yaml:"title"`
	Publish bool       `json:"publish" yaml:"publish"`
	Created int64      `json:"created" yaml:"created"`
	Content []*Content `json:"content" yaml:"content"`
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
