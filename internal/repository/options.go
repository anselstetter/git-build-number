package repository

type commitOptions struct {
	author  Author
	setHead bool
	headers []Header
}

type commitOption func(opts *commitOptions)

func newCommitOptions(option ...commitOption) commitOptions {
	opts := commitOptions{
		setHead: false,
		author: Author{
			Name:  "Build Number",
			Email: "Not Set",
		},
		headers: []Header{},
	}
	for _, fn := range option {
		fn(&opts)
	}
	return opts
}

func WithAuthor(author Author) commitOption {
	return func(opts *commitOptions) {
		opts.author = author
	}
}

func WithHead() commitOption {
	return func(opts *commitOptions) {
		opts.setHead = true
	}
}

func WithHeaders(headers []Header) commitOption {
	return func(opts *commitOptions) {
		opts.headers = headers
	}
}

type refsOptions struct {
	prefix *string
}

type refsOption func(opts *refsOptions)

func newRefsOptions(option ...refsOption) refsOptions {
	opts := refsOptions{}
	for _, fn := range option {
		fn(&opts)
	}
	return opts
}

func WithPrefix(prefix string) refsOption {
	return func(opts *refsOptions) {
		opts.prefix = &prefix
	}
}

type commitsOptions struct {
	headerKey *string
}

type commitsOption func(opts *commitsOptions)

func newCommitsOptions(option ...commitsOption) commitsOptions {
	opts := commitsOptions{}
	for _, fn := range option {
		fn(&opts)
	}
	return opts
}

func WithHeaderKey(key string) commitsOption {
	return func(opts *commitsOptions) {
		opts.headerKey = &key
	}
}
