package jsonfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type fileDB struct {
	file    *os.File
	encoder *json.Encoder
}

type fileContent struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

func newFileDB(file *os.File) *fileDB {
	encoder := json.NewEncoder(&tape{file})
	encoder.SetIndent("", "    ")
	return &fileDB{file, encoder}
}

func (f *fileDB) Save(gauges map[string]float64, counters map[string]int64) error {
	content := fileContent{counters, gauges}

	if err := f.encoder.Encode(content); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	if err := f.file.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %v", err)
	}

	return nil
}

func (f *fileDB) Load() (g map[string]float64, c map[string]int64, err error) {
	stat, err := f.file.Stat()
	if err != nil {
		err = fmt.Errorf("error getting file info from file %s: %v", f.file.Name(), err)
		return
	}

	if stat.Size() == 0 {
		return
	}

	var content fileContent
	if err = json.NewDecoder(f.file).Decode(&content); err != nil {
		err = fmt.Errorf("error decoding file content for metrics restore: %v", err)
		return
	}

	return content.Gauges, content.Counters, nil
}

func (f *fileDB) Close() {
	f.file.Close()
}

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (int, error) {
	if err := t.file.Truncate(0); err != nil {
		return 0, err
	}
	if _, err := t.file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	return t.file.Write(p)
}
