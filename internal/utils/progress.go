package utils

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Progress struct {
	total   int64
	current int64
	mutex   sync.Mutex
}

func NewProgress(total int64) *Progress {
	return &Progress{total: total}
}

func (p *Progress) Update(bytes int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.current += bytes
	p.display()
}

func (p *Progress) display() {
	percentage := float64(p.current) / float64(p.total) * 100
	fmt.Printf("\rProgress: [%-50s] %.2f%%", p.progressBar(), percentage)
	if p.current >= p.total {
		fmt.Println("\nDone!")
	}
}

func (p *Progress) progressBar() string {
	barLength := 50
	progressLength := int(float64(barLength) * float64(p.current) / float64(p.total))
	return string(repeat('#', progressLength)) + string(repeat('-', barLength-progressLength))
}

func repeat(char rune, n int) []rune {
	return []rune{char}
}

func CopyFileWithProgress(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	progress := NewProgress(fileInfo.Size())
	buffer := make([]byte, 1024*1024)
	for {
		n, err := sourceFile.Read(buffer)
		if n > 0 {
			if _, err := destFile.Write(buffer[:n]); err != nil {
				return err
			}
			progress.Update(int64(n))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
