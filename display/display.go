package display

type Options struct {
	showBytes  bool
	showWords  bool
	showLines  bool
	showHeader bool
}

func NewOptions(bytes, words, lines, header bool) Options {
	return Options{
		showBytes:  bytes,
		showWords:  words,
		showLines:  lines,
		showHeader: header,
	}
}

func (d Options) ShouldShowBytes() bool {
	if !d.showBytes && !d.showWords && !d.showLines {
		return true
	}

	return d.showBytes
}

func (d Options) ShouldShowWords() bool {
	if !d.showBytes && !d.showWords && !d.showLines {
		return true
	}

	return d.showWords
}
func (d Options) ShouldShowLines() bool {
	if !d.showBytes && !d.showWords && !d.showLines {
		return true
	}

	return d.showLines
}

func (d Options) ShouldShowHeader() bool {
	return d.showHeader
}
