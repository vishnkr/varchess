package main 
func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}
func Abs(x int)int{
	if x<0{
		return -x
	}
	return x
}
func Min(x, y int) int {
    if x > y {
        return y
    }
    return x
}