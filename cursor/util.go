package cursor

// // DeltaMarkup returns markup to transform from one cursor to another
// func DeltaMarkup(from, to Cursor) string {
// 	out := ""

// 	// xDiff := to.X - from.X
// 	// yDiff := to.Y - from.Y

// 	// if yDiff == 0 {
// 	// 	if xDiff > 0 {
// 	// 		out += fmt.Sprintf("\033[%dC", xDiff) // move forwards
// 	// 	} else {
// 	// 		out += fmt.Sprintf("\033[%dD", -xDiff) // move backwards
// 	// 	}
// 	// } else {
// 	// 	out += fmt.Sprintf("\033[%d;%dH", to.Y, to.X)
// 	// }

// 	out += fmt.Sprintf("\033[%d;%dH", to.Y+1, to.X+1)

// 	if to.Bg.ColorMode != from.Bg.ColorMode || to.Bg.Code != from.Bg.Code {
// 		out += to.Bg.ToANSI(true)
// 	}

// 	if to.Fg.ColorMode != from.Fg.ColorMode || to.Fg.Code != from.Fg.Code {
// 		out += to.Fg.ToANSI(false)
// 	}

// 	/* removing effects */

// 	if from.Faint && !to.Faint {
// 		out += "\033[22m"
// 	}

// 	if from.Underline && !to.Underline {
// 		out += "\033[24m"
// 	}

// 	/* adding effects */

// 	if !from.Faint && to.Faint {
// 		out += "\033[2m"
// 	}

// 	if !from.Underline && to.Underline {
// 		out += "\033[4m"
// 	}

// 	return out
// }
