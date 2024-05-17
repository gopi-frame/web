package cors

// Option option
type Option func(cors *CORS)

// AllowCredentials allow credentials
func AllowCredentials(value bool) Option {
	return func(cors *CORS) {
		cors.AllowCredentials = true
	}
}

// AllowHeaders allow headers
func AllowHeaders(headers ...string) Option {
	return func(cors *CORS) {
		cors.AllowHeaders = headers
	}
}

// AllowOrigins allow origins
func AllowOrigins(origins ...string) Option {
	return func(cors *CORS) {
		cors.AllowOrigins = origins
	}
}

// AllowMethods all methods
func AllowMethods(methods ...string) Option {
	return func(cors *CORS) {
		cors.AllowMethods = methods
	}
}

// ExposeHeaders expose headers
func ExposeHeaders(headers ...string) Option {
	return func(cors *CORS) {
		cors.ExposeHeaders = headers
	}
}
