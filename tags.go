package stork

// IgnoreTags are tags that will be removed from the document before analysing it.
// This list contains tags such as metadata elements that don't make sense in
// the context of the extracted content.
//
// Note:
// Fresh and clean new metadata will be added afterwards if using the Html() method
// in order to re-create a full and complete html document.
//
// https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/Content_categories#Metadata_content
//
// Other tags can be added to this list if they should be removed from the extracted document.
var IgnoreTags = map[string]bool{
	"base":     true,
	"command":  true,
	"link":     true,
	"meta":     true,
	"noscript": true,
	"script":   true,
	"style":    true,
	"title":    true,
}

// voidTags is a list of void elements. Void elements
// are those that can't have any contents.
var voidTags = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

// blockTags are elements that always start on a new line and takes up the full width available
// They are used to determine what is a structural tag in order to extract the main content of the page.
//
// https://www.w3schools.com/html/html_blocks.asp
var blockTags = map[string]bool{
	"address":    true,
	"article":    true,
	"aside":      true,
	"blockquote": true,
	"canvas":     true,
	"dd":         true,
	"div":        true,
	"dl":         true,
	"dt":         true,
	"fieldset":   true,
	"figcaption": true,
	"figure":     true,
	"footer":     true,
	"form":       true,
	"h1":         true,
	"h2":         true,
	"h3":         true,
	"h4":         true,
	"h5":         true,
	"h6":         true,
	"header":     true,
	"hr":         true,
	"li":         true,
	"main":       true,
	"nav":        true,
	"noscript":   true,
	"ol":         true,
	"p":          true,
	"pre":        true,
	"section":    true,
	"table":      true,
	"tfoot":      true,
	"ul":         true,
	"video":      true,
}

// inlineTags don't start on a new line and only take up as much width as necessary.
//
// https://www.w3schools.com/html/html_blocks.asp
var inlineTags = map[string]bool{
	"a":        true,
	"abbr":     true,
	"acronym":  true,
	"b":        true,
	"bdo":      true,
	"big":      true,
	"br":       true,
	"button":   true,
	"cite":     true,
	"code":     true,
	"dfn":      true,
	"em":       true,
	"i":        true,
	"img":      true,
	"input":    true,
	"kbd":      true,
	"label":    true,
	"map":      true,
	"object":   true,
	"output":   true,
	"q":        true,
	"samp":     true,
	"script":   true,
	"select":   true,
	"small":    true,
	"span":     true,
	"strong":   true,
	"sub":      true,
	"sup":      true,
	"textarea": true,
	"time":     true,
	"tt":       true,
	"var":      true,
}
