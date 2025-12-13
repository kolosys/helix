package helix

// MIME type constants for HTTP Content-Type headers.
// Base types (without charset) are used for content-type detection/matching.
// CharsetUTF8 variants are used for setting response headers.
const (
	// Text types - base (for matching)
	MIMETextPlain      = "text/plain"
	MIMETextHTML       = "text/html"
	MIMETextCSS        = "text/css"
	MIMETextCSV        = "text/csv"
	MIMETextJavaScript = "text/javascript"
	MIMETextXML        = "text/xml"

	// Text types - with charset (for responses)
	MIMETextPlainCharsetUTF8      = "text/plain; charset=utf-8"
	MIMETextHTMLCharsetUTF8       = "text/html; charset=utf-8"
	MIMETextCSSCharsetUTF8        = "text/css; charset=utf-8"
	MIMETextCSVCharsetUTF8        = "text/csv; charset=utf-8"
	MIMETextJavaScriptCharsetUTF8 = "text/javascript; charset=utf-8"
	MIMETextXMLCharsetUTF8        = "text/xml; charset=utf-8"

	// Application types - base (for matching)
	MIMEApplicationJSON       = "application/json"
	MIMEApplicationXML        = "application/xml"
	MIMEApplicationJavaScript = "application/javascript"
	MIMEApplicationXHTMLXML   = "application/xhtml+xml"

	// Application types - with charset (for responses)
	MIMEApplicationJSONCharsetUTF8       = "application/json; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8        = "application/xml; charset=utf-8"
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"

	// Application types - no charset needed
	MIMEApplicationProblemJSON = "application/problem+json"
	MIMEApplicationForm        = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf    = "application/x-protobuf"
	MIMEApplicationMsgPack     = "application/msgpack"
	MIMEApplicationOctetStream = "application/octet-stream"
	MIMEApplicationPDF         = "application/pdf"
	MIMEApplicationZip         = "application/zip"
	MIMEApplicationGzip        = "application/gzip"
	MIMEMultipartForm          = "multipart/form-data"

	// Image types
	MIMEImagePNG  = "image/png"
	MIMEImageSVG  = "image/svg+xml"
	MIMEImageJPEG = "image/jpeg"
	MIMEImageGIF  = "image/gif"
	MIMEImageWebP = "image/webp"
	MIMEImageICO  = "image/x-icon"
	MIMEImageAVIF = "image/avif"

	// Audio types
	MIMEAudioMPEG = "audio/mpeg"
	MIMEAudioWAV  = "audio/wav"
	MIMEAudioOGG  = "audio/ogg"

	// Video types
	MIMEVideoMP4  = "video/mp4"
	MIMEVideoWebM = "video/webm"
	MIMEVideoOGG  = "video/ogg"
)
