package torrent

import "github.com/anacrolix/torrent/metainfo"

// Provides access to regions of torrent data that correspond to its files.
type File struct {
	t      Torrent
	path   string
	offset int64
	length int64
	fi     metainfo.FileInfo
}

// Data for this file begins this far into the torrent.
func (f *File) Offset() int64 {
	return f.offset
}

func (f File) FileInfo() metainfo.FileInfo {
	return f.fi
}

func (f File) Path() string {
	return f.path
}

func (f *File) Length() int64 {
	return f.length
}

type FilePieceState struct {
	Bytes int64 // Bytes within the piece that are part of this File.
	PieceState
}

// Returns the state of pieces in this file.
func (f *File) State() (ret []FilePieceState) {
	f.t.cl.mu.Lock()
	defer f.t.cl.mu.Unlock()
	pieceSize := int64(f.t.usualPieceSize())
	off := f.offset % pieceSize
	remaining := f.length
	for i := int(f.offset / pieceSize); ; i++ {
		if remaining == 0 {
			break
		}
		len1 := pieceSize - off
		if len1 > remaining {
			len1 = remaining
		}
		ret = append(ret, FilePieceState{len1, f.t.pieceState(i)})
		off = 0
		remaining -= len1
	}
	return
}

func (f *File) PrioritizeRegion(off, len int64) {
	if off < 0 || off >= f.length {
		return
	}
	if off+len > f.length {
		len = f.length - off
	}
	off += f.offset
	f.t.SetRegionPriority(off, len)
}
