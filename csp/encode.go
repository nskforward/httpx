package csp

import "strings"

func Encode(cfg Config) string {
	arr := make([]string, 0, 32)

	addValue(&arr, "default-src", cfg.DefaultSources)
	addValue(&arr, "script-src", cfg.ScriptSources)
	addValue(&arr, "style-src", cfg.StyleSources)
	addValue(&arr, "img-src", cfg.ImageSources)
	addValue(&arr, "media-src", cfg.MediaSources)
	addValue(&arr, "connect-src", cfg.ConnectSources)
	addValue(&arr, "font-src", cfg.FontSources)
	addValue(&arr, "object-src", cfg.ObjectSources)
	addValue(&arr, "manifest-src", cfg.ManifestSources)
	addValue(&arr, "base-uri", cfg.BaseURI)
	addValue(&arr, "frame-ancestors", cfg.FrameAncestors)
	addValue(&arr, "form-action", cfg.FormcAtion)
	addValue(&arr, "child-src", cfg.ChildSources)

	return strings.Join(arr, "; ")
}

func addValue(dst *[]string, directive string, values []string) {
	if len(values) > 0 {
		res := append([]string{directive}, values...)
		*dst = append(*dst, strings.Join(res, " "))
	}
}
