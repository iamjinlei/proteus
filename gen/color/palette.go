package color

var (
	DefaultPalette = Palette{
		Red:       Red,
		Green:     Green,
		Blue:      Blue,
		Yellow:    Yellow,
		Orange:    Orange,
		Pink:      Pink,
		LightGray: LightGray,
	}
)

type Palette struct {
	Red       Color
	Green     Color
	Blue      Color
	Yellow    Color
	Orange    Color
	Pink      Color
	LightGray Color
}
