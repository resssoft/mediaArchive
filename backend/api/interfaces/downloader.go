package interfaces

type IDownLoader interface {
	Download(url string) (string, error)
	Response() int
}
