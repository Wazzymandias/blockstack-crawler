package storage

type Storage interface {
	Exists(key interface{}) bool
	Find(path string) (interface{}, bool)
	DirExists(path string) bool
	IsDirEmpty(path string) bool
	DirSize(path string) int
}
