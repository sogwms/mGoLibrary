package def

type Writer interface {
	Write(data *FileData) error
}
