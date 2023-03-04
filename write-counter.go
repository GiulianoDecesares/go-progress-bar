package progressbar

type WriteCounter struct {
	Total       int64
	Partial     int64
	progressBar *Bar
}

func NewWriteCounter(totalBytes int64, counterName string) *WriteCounter {
	counter := new(WriteCounter)

	counter.Total = totalBytes
	counter.Partial = 0
	counter.progressBar = NewProgressBar(counter.Partial, counter.Total, counterName)

	return counter
}

func (writeCounter *WriteCounter) Write(data []byte) (int, error) {
	dataLength := len(data)

	writeCounter.Partial += int64(dataLength)

	// Update progress bar
	writeCounter.progressBar.Update(writeCounter.Partial)

	if writeCounter.Partial == writeCounter.Total {
		writeCounter.progressBar.Finish()
	}

	return dataLength, nil
}
