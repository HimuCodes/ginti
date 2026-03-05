package counter

type DisplayOptions struct {
	ShowBytes  bool
	ShowWords  bool
	ShowLines  bool
	ShowHeader bool
}

func (d DisplayOptions) ShouldShowBytes() bool {
	if !d.ShowBytes && !d.ShowWords && !d.ShowLines {
		return true
	}

	return d.ShowBytes
}

func (d DisplayOptions) ShouldShowWords() bool {
	if !d.ShowBytes && !d.ShowWords && !d.ShowLines {
		return true
	}

	return d.ShowWords
}
func (d DisplayOptions) ShouldShowLines() bool {
	if !d.ShowBytes && !d.ShowWords && !d.ShowLines {
		return true
	}

	return d.ShowLines
}
