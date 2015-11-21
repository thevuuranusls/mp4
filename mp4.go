
package mp4

import (
	"./atom"
	"os"
	"io"
	"log"
)

type File struct {
}

func (self *File) AddAvcc(avcc *Avcc) {
}

func (self *File) AddMp4a(mp4a *Mp4a) {
}

func (self *File) GetAvcc() (avcc []*Avcc) {
	return
}

func (self *File) GetMp4a() (mp4a []*Mp4a) {
	return
}

func (self *File) Sync() {
}

func (self *File) Close() {
}

func Open(filename string) (file *File, err error) {
	var osfile *os.File
	if osfile, err = os.Open(filename); err != nil {
		return
	}

	var finfo os.FileInfo
	if finfo, err = osfile.Stat(); err != nil {
		return
	}
	log.Println("filesize", finfo.Size())

	lr := &io.LimitedReader{R: osfile, N: finfo.Size()}

	var outfile *os.File
	if outfile, err = os.Create(filename+".out.mp4"); err != nil {
		return
	}

	for lr.N > 0 {
		var ar *io.LimitedReader

		var cc4 string
		if ar, cc4, err = atom.ReadAtomHeader(lr, ""); err != nil {
			return
		}

		if cc4 == "moov" {
			curPos, _ := outfile.Seek(0, 1)
			origSize := ar.N+8
			var moov *atom.Movie
			if moov, err = atom.ReadMovie(ar); err != nil {
				return
			}
			if err = atom.WriteMovie(outfile, moov); err != nil {
				return
			}
			curPosAfterRead, _ := outfile.Seek(0, 1)
			bytesWritten := curPosAfterRead - curPos

			log.Println("regen moov", "tracks nr", len(moov.Tracks),
				"origSize", origSize, "bytesWritten", bytesWritten,
			)

			padSize := origSize - bytesWritten - 8
			aw, _ := atom.WriteAtomHeader(outfile, "free")
			atom.WriteDummy(outfile, int(padSize))
			aw.Close()
		} else {
			var aw *atom.Writer
			if aw, err = atom.WriteAtomHeader(outfile, cc4); err != nil {
				return
			}
			log.Println("copy", cc4)
			if _, err = io.CopyN(aw, ar, ar.N); err != nil {
				return
			}
			if err = aw.Close(); err != nil {
				return
			}
		}

		//log.Println("atom", cc4, "left", lr.N)
		//atom.ReadDummy(ar, int(ar.N))
	}

	if err = outfile.Close(); err != nil {
		return
	}

	return
}

func Create(filename string) (file *File, err error) {
	return
}
