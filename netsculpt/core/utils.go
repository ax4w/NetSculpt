package core

func isAllOne(l []int) bool {
	for _, v := range l {
		if v != 1 {
			return false
		}
	}
	return true
}
