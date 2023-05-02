package slogx

import (
	"io"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	SizeMB = 1024 * 1024

	RotateTimeLayoutEveryYear    = "2006"
	RotateTimeLayoutEveryMonth   = "200601"
	RotateTimeLayoutEveryDay     = "20060102"
	RotateTimeLayoutEveryHour    = "2006010215"
	RotateTimeLayoutEveryMinutes = "200601021504"
	RotateTimeLayoutEverySecond  = "20060102150405"
)

type FileWriterOption interface {
	apply(*FileWriterOptions)
}

type FileWriterOptionFn struct {
	f func(*FileWriterOptions)
}

func NewFileWriterOptionFn(f func(*FileWriterOptions)) *FileWriterOptionFn {
	return &FileWriterOptionFn{f: f}
}

func (fn FileWriterOptionFn) apply(opts *FileWriterOptions) {
	fn.f(opts)
}

var _ io.Writer = (*FileWriter)(nil)

var defaultFileWriterOptions = FileWriterOptions{
	Path:             "logs",
	FileName:         "app.log",
	MaxTime:          0,
	MaxCount:         0,
	RotateTimeLayout: RotateTimeLayoutEveryHour,
	RotateSize:       0,
	Compress:         false,
}

type FileWriterOptions struct {
	Path             string
	FileName         string
	MaxTime          time.Duration
	MaxCount         int
	RotateTimeLayout string
	// RotateSize is the maximum size in MB of the log file before it gets rotated.
	// It defaults to 100 MB.
	RotateSize int
	Compress   bool
}

func (opts FileWriterOptions) NewFileWriter() (*FileWriter, error) {
	info, err := os.Stat(opts.Path)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to stat path")
	}
	if !info.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	return &FileWriter{
		opts:              opts,
		mutex:             sync.Mutex{},
		currentFileWriter: nil,
		mainFileName:      path.Join(opts.Path, opts.FileName),
		currentTimeSuffix: "",
	}, nil
}

type FileWriter struct {
	opts              FileWriterOptions
	mutex             sync.Mutex
	currentFileWriter *os.File
	currentTimeSuffix string
	mainFileName      string
}

func NewFileWriter(opts ...FileWriterOption) (*FileWriter, error) {
	fwOpts := defaultFileWriterOptions
	for _, opt := range opts {
		opt.apply(&fwOpts)
	}
	return fwOpts.NewFileWriter()
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	writer, err := w.getCurrentWriter()
	if err != nil {
		return 0, err
	}

	return writer.Write(p)
}

func (w *FileWriter) getCurrentWriter() (io.Writer, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.currentFileWriter == nil || !w.hasMainFile() {
		fw, err := os.OpenFile(w.mainFileName,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		w.currentFileWriter = fw
		w.currentTimeSuffix = w.getRotateTimeSuffix()
		return w.currentFileWriter, nil
	}

	timeSuffix := w.getRotateTimeSuffix()
	sizeSuffix := w.getRotateSizeSuffix()
	if timeSuffix == w.currentTimeSuffix && sizeSuffix == "" {
		return w.currentFileWriter, nil
	}

	filename := w.mainFileName + timeSuffix + sizeSuffix
	err := os.Rename(w.mainFileName, filename)
	if err != nil {
		return nil, err
	}
	err = w.currentFileWriter.Close()
	if err != nil {
		return nil, err
	}

	fw, err := os.OpenFile(w.mainFileName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	w.currentFileWriter = fw
	w.currentTimeSuffix = timeSuffix
	return w.currentFileWriter, err
}

func (w *FileWriter) getRotateTimeSuffix() string {
	if w.opts.RotateTimeLayout == "" {
		return ""
	}
	return "." + time.Now().Format(w.opts.RotateTimeLayout)
}

func (w *FileWriter) getRotateSizeSuffix() string {
	if w.opts.RotateSize == 0 {
		return ""
	}
	info, _ := w.currentFileWriter.Stat()
	if info.Size()*SizeMB >= int64(w.opts.RotateSize) {
		return "." + strconv.Itoa(int(time.Now().Unix()))
	}

	return ""
}

func (w *FileWriter) hasMainFile() bool {
	_, err := os.Stat(w.mainFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return true
}

func WithPath(path string) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.Path = path
	})
}

func WithFileName(fileName string) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.FileName = fileName
	})
}

func WithMaxTime(maxTime time.Duration) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.MaxTime = maxTime
	})
}

func WithMaxCount(maxCount int) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.MaxCount = maxCount
	})
}

func WithRotateTimeLayout(rotateTimeLayout string) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.RotateTimeLayout = rotateTimeLayout
	})
}

func WithRotateSize(rotateSize int) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.RotateSize = rotateSize
	})
}

func WithCompress(compress bool) FileWriterOption {
	return NewFileWriterOptionFn(func(opts *FileWriterOptions) {
		opts.Compress = compress
	})
}
