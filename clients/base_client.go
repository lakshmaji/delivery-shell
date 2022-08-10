package clients

// Handles responsibility of writing to given Writer
type BaseWriter interface {
	Write(interface{})
	WriteError(interface{})
}
