package drawables

type drawableText struct {
	content [][]string
}

func (d *drawableText) Draw() [][]rune {
	var res [][]rune
	for _, line := range d.content {
		row := []rune{}
		for _, word := range line {
			row = append(row, []rune(word)...)
		}
		res = append(res, row)
	}
	return res
}

func NewText(content [][]string) *drawableText {
	return &drawableText{
		content: content,
	}
}
