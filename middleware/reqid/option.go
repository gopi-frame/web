package reqid

// Option option
type Option func(*RequestID)

// Header set header name
func Header(header string) Option {
	return func(ri *RequestID) {
		if header != "" {
			ri.Header = header
		}
	}
}

// Generator set id generator
func Generator(generator func() string) Option {
	return func(ri *RequestID) {
		if generator != nil {
			ri.Generator = generator
		}
	}
}
