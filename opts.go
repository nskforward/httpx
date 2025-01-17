package httpx

type SetOpt func(r *Router)

// WithSlashRedirection redirects to url without trailing slashes if shouldRedirect is true
func WithSlashRedirection(shouldRedirect bool) SetOpt {
	return func(r *Router) {
		r.slashRedirect = shouldRedirect
	}
}

func WithLogBeforeRequest(f LogFunc) SetOpt {
	return func(r *Router) {
		r.beforeRequestLog = f
	}
}

func WithLogAfterResponse(f LogFunc) SetOpt {
	return func(r *Router) {
		r.afterResponseLog = f
	}
}
