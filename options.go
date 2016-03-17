package goprox

// Options is a configuration container to setup the GoProx middleware.
type Options struct {
	// Destination is the URL for the host where requests will be sent to
	Destination string
	// Path will restrict the proxy requests to the specified path.
	Path string
	// Debugging flag adds additional output to debug server side issues
	Debug bool
	// Cache enables GET requests caching
	Cache bool
	// EnableCORS enables cors
	EnableCORS bool
}
