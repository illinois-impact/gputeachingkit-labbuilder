package pandoc

type config struct {
	IsFullTookit bool
}

const (
	// eval and full are now the same when released
	BitbucketURLEval = "https://bitbucket.org/hwuligans/gputeachingkit-labs/src/master/"
	BitbucketURLFull = "https://bitbucket.org/hwuligans/gputeachingkit-labs/src/master/"

	MarkdownFormat = "markdown+hard_line_breaks+pandoc_title_block+lists_without_preceding_blankline+" +
		"compact_definition_lists+simple_tables+table_captions"
)

var (
	DefaultFilter = []string{
		"--highlight-style",
		"pygments",
		"--self-contained",
	}
	Config = config{
		IsFullTookit: true,
	}
)

func init() {
	//DefaultFilter = []string{
	//	"--filter",
	//	"pandoc-crossref",
	//	"--filter",
	//	"pandoc-citeproc",
	//	"--filter",
	//	"pandoc-citeproc-preamble",
	//}
}
