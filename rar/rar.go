package rar

import (
	"fmt"
	"github.com/nwaples/rardecode/v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Reader struct {
	r        *rardecode.ReadCloser
	password string
}

func OpenReader(source string, password string) (Reader, error) {
	var opts []rardecode.Option
	if len(password) >= 0 {
		opts = append(opts, rardecode.Password(password))
	}

	r, err := rardecode.OpenReader(source, opts...)
	if err != nil {
		return Reader{}, err
	}

	rar := Reader{
		r: r,
	}
	return rar, nil
}

func (r Reader) Extract(output string) {
	for {
		f, err := r.r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error opening rar encoded file: %s\n", err)
			continue
		}
		path := fmt.Sprintf("%s/%s", output, f.Name)

		if f.IsDir {
			_ = os.MkdirAll(filepath.Dir(path), 0)
			continue
		}

		buf := make([]byte, 1024)
		n, err := r.r.Read(buf)
		if err != nil {
			log.Printf("error reading rar encoded file: %s\n", err)
			continue
		}
		buf = buf[:n]

		if err = os.WriteFile(path, buf, 0); err != nil {
			fmt.Printf("error decoding rar encoded file: %s\n", err)
		}
	}
	_ = r.r.Close()
}
