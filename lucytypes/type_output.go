package lucytypes

type Field interface {
	Output()
}

type OutputData struct {
	Fields []Field
}
