package csp

const (
	Self = "'self'"
	None = "'none'"
	Any  = "*"
)

type Config struct {

	// default-src is a fallback for all other fetch directives
	DefaultSources []string // default-src

	// img-src specifies valid sources of images and favicons.
	ImageSources []string // img-src

	// media-src specifies valid sources for loading media using the <audio>, <video> and <track> elements.
	MediaSources []string // media-src

	// connect-src restricts the URLs which can be loaded using script interfaces
	ConnectSources []string // connect-src

	// font-src specifies valid sources for fonts loaded using @font-face.
	FontSources []string // font-src

	// object-src specifies valid sources for the <object> and <embed> elements.
	ObjectSources []string // object-src

	// manifest-src specifies valid sources of application manifest files.
	ManifestSources []string // manifest-src

	// script-src specifies valid sources for JavaScript and WebAssembly resources.
	ScriptSources []string // script-src

	// base-uri restricts the URLs which can be used in a document's <base> element.
	BaseURI []string // base-uri

	// frame-ancestors specifies valid parents that may embed a page using <frame>, <iframe>, <object>, or <embed>.
	FrameAncestors []string // frame-ancestors

	// form-action restricts the URLs which can be used as the target of a form submissions from a given context.
	FormcAtion []string // form-action

	//	style-src specifies valid sources for stylesheets.
	StyleSources []string // style-src

	// child-src defines the valid sources for web workers and nested browsing contexts loaded using elements such as <frame> and <iframe>.
	ChildSources []string
}

func DefaultConfig() Config {
	return Config{
		DefaultSources: []string{Self},
		ObjectSources:  []string{None},
		BaseURI:        []string{None},
		FrameAncestors: []string{None}, // a more flexible replacement for the X-Frame-Options header
	}
}
