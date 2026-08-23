package __variable

import "strconv"

func iMain() {
	x1, _ := strconv.ParseInt("-10", 0, 32)
	x2, _ := strconv.ParseUint("10", 0, 32)
}
