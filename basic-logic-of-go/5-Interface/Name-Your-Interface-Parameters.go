package __Interface

type Copier interface {
	Copy(string, string) int
}
type Copier1 interface {
	Copy(sourceFile string, destinationFile string) (bytesCopied int)
}
