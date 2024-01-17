package internal

func WithFilters(filters map[string][]string) ListImageOptions {
	return func(o *listImageOptions) {
		o.filters = filters
	}
}

func WithFilter(filter string) ListImageOptions {
	return func(o *listImageOptions) {
		o.filter = filter
	}
}

func WithAll(all bool) ListImageOptions {
	return func(o *listImageOptions) {
		o.all = all
	}
}

func WithDigests(digests bool) ListImageOptions {
	return func(o *listImageOptions) {
		o.digests = digests
	}
}
