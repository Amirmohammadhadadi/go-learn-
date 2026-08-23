package __Functions

//Explicit Returns
//Even though a function has named return values, we can still explicitly return values if we want to.
func getCoords() (x, y int) {
	return x, y // this is explicit
}
//Using this explicit pattern we can even overwrite the return values:
func getCoords() (x, y int) {
   return 5, 6 // this is explicit, x and y are NOT returned
}
//Otherwise, if we want to return the values defined in the function signature we can just use a naked return (blank return):
func getCoords() (x, y int) {
   return // implicitly returns x and y
}
