package structs

type POSIX_HEADER struct {
	Name       [100]byte // 0
	Mode       [8]byte   // 100
	Uid        [8]byte   // 108
	Gid        [8]byte   // 116
	Size       [12]byte  // 124
	Mtime      [12]byte  // 136
	Chksum     [8]byte   // 148
	Typeflag   byte      // 156
	Linkname   [100]byte // 157
	Magic      [6]byte   // 257
	Version    [2]byte   // 263
	Uname      [32]byte  // 265
	Gname      [32]byte  // 297
	Devmajor   [8]byte   // 329
	Devmintor  [8]byte   // 337
	Prefix     [155]byte // 345
	Endpadding [12]byte  // 500
	// 512
}

type OLDGNU_HEADER struct {
	Unused_pad1 [345]byte // 0 NOT USED
	Atime       [12]byte  // 345
	Ctime       [12]byte  // 357
	Offset      [12]byte  // 369 // Offset of this volume if multivolume archive
	Longnames   [4]byte   // 381 NOT USED
	Unused_pad2 byte      // 385 NOT USED
	Sparses     [4]SPARSE // 386
	IsExtended  bool      // 482 // IS THE FOLLOWING FILE SPARSE - SPARSE HEADER FOLLOWS IF TRUE
	Realsize    [12]byte  // 483 // Sparse file: Real size
	// 495

}

type SPARSE struct {
	Offset   [12]byte // 0
	Numbytes [12]byte // 12
	// 24

}

type SPARSE_HEADER struct {
	Sparses    [21]SPARSE // 0
	IsExtended byte       // 504
	// 505

}
