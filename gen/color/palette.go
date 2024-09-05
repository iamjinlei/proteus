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
		DarkGray:  DarkGray,

		HighlighterRed:    HighlighterUltraRed,
		HighlighterGreen:  HighlighterFrenchLime,
		HighlighterBlue:   HighlighterMayaBlue,
		HighlighterYellow: HighlighterLaserLemon,
		HighlighterOrange: HighlighterMacCheese,
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
	DarkGray  Color

	HighlighterRed    Color
	HighlighterGreen  Color
	HighlighterBlue   Color
	HighlighterYellow Color
	HighlighterOrange Color
}
