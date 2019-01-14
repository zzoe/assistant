package format

var (
	Normal = NewStyle(&Style{
		Border: []Border{
			Border{
				Type:  "top",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "bottom",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "left",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "right",
				Color: "#A9D08E",
				Style: 1,
			},
		},
		Font: Font{
			Bold:   false,
			Italic: false,
			Family: "宋休",
			Size:   8,
			Color:  "#000000",
		},
		Alignment: Alignment{
			Horizontal:  "center",
			Vertical:    "center",
			WrapText:    true,
			ShrinkToFit: true,
		},
	})

	YellowFill = NewStyle(&Style{
		Fill: Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
		Border: []Border{
			Border{
				Type:  "top",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "bottom",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "left",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "right",
				Color: "#A9D08E",
				Style: 1,
			},
		},
		Font: Font{
			Bold:   false,
			Italic: false,
			Family: "宋休",
			Size:   8,
			Color:  "#000000",
		},
		Alignment: Alignment{
			Horizontal:  "center",
			Vertical:    "center",
			WrapText:    true,
			ShrinkToFit: true,
		},
	})

	RedFont = NewStyle(&Style{
		Border: []Border{
			Border{
				Type:  "top",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "bottom",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "left",
				Color: "#A9D08E",
				Style: 1,
			},
			Border{
				Type:  "right",
				Color: "#A9D08E",
				Style: 1,
			},
		},
		Font: Font{
			Bold:   false,
			Italic: false,
			Family: "宋休",
			Size:   8,
			Color:  "#FF0000",
		},
		Alignment: Alignment{
			Horizontal:  "center",
			Vertical:    "center",
			WrapText:    true,
			ShrinkToFit: true,
		},
	})
)
