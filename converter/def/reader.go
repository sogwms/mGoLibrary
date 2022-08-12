package def

type Reader interface {
	GetData() (fd *FileData, err error)
}
