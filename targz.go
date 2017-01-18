package targz

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"os"
	"os/user"
	"strconv"

	"github.com/sisatech/targz/structs"
)

func ArchiveFile(path, destination string) error {
	// Get file path & check exists ...
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	// os.Stat for size ...
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Write header ...
	header := new(structs.POSIX_HEADER)

	// Header Name ...
	copy(header.Name[:], []byte(stat.Name()))

	u, err := user.Current()
	if err != nil {
		return err
	}

	// User ID and Group ID ...
	uidInt, _ := strconv.Atoi(u.Uid)
	gidInt, _ := strconv.Atoi(u.Gid)
	uid := "000" + strconv.FormatUint(uint64(uidInt), 8)
	gid := "000" + strconv.FormatUint(uint64(gidInt), 8)

	copy(header.Uid[:], []byte(uid))
	copy(header.Gid[:], []byte(gid))

	// Size ...
	octSize := "0" + strconv.FormatUint(uint64(stat.Size()), 8)
	copy(header.Size[:], []byte(octSize))

	// Mode ...
	var mode [8]byte
	mode[0] = byte(0x30)
	mode[1] = byte(0x30)
	mode[2] = byte(0x30)
	mode[3] = byte(0x30)
	mode[4] = byte(0x36)
	mode[5] = byte(0x34)
	mode[6] = byte(0x34)
	mode[7] = byte(0x00)

	header.Mode = mode

	// Time ...
	time := "00000000000"
	copy(header.Mtime[:], []byte(time))

	// Typeflag ...
	header.Typeflag = 0x30

	// Magic ...
	magic := []byte("ustar")
	byteMagic := make([]byte, len(magic)+1)

	for i := 0; i < len(byteMagic); i++ {
		if i == len(byteMagic)-1 {
			byteMagic[i] = byte(0x20)
		} else {
			byteMagic[i] = magic[i]
		}
	}

	copy(header.Magic[:], byteMagic)

	// Version ...
	var version [2]byte
	version[0] = 0x20
	version[1] = 0x00
	header.Version = version

	// Uname ...
	uname := []byte(u.Username)
	copy(header.Uname[:], uname)

	// Gname ...
	group, err := user.LookupGroupId(u.Gid)
	if err != nil {
		return err
	}
	gname := []byte(group.Name)
	copy(header.Gname[:], gname)

	// Calculate Checksum ...
	headerBytes := &bytes.Buffer{}
	err = binary.Write(headerBytes, binary.LittleEndian, header)
	if err != nil {
		return err
	}
	bh := make([]byte, 512)
	err = binary.Read(headerBytes, binary.LittleEndian, &bh)

	var csVal uint32
	csVal = 0
	for i := 0; i < len(bh); i++ {
		csVal += uint32(bh[i])
	}

	// Treat the 8 Checksum Field bytes as ASCII spaces (dec 32)
	csVal += (8 * 32)

	csOctal := "0" + strconv.FormatUint(uint64(csVal), 8)
	csBytes := []byte(csOctal)
	copy(header.Chksum[:], csBytes)
	header.Chksum[len(header.Chksum)-1] = 0x20

	// binary.Write(buf, binary.LittleEndian, header)

	// Create Archive File
	f, err := os.Create(path + ".tar")
	if err != nil {
		return err
	}

	// Create Reader on targeted file ...
	r, err := os.Open(path)
	if err != nil {
		return err
	}

	// Write Header to Archive File ...
	// f.WriteAt([]byte(fmt.Sprintf("%v", header)), 0)
	binary.Write(f, binary.LittleEndian, header)

	// Offset by headersize & write to archive
	f.Seek(512, 0)
	io.Copy(f, r)

	// Offset by file size & write end 2 empty blocks
	f.Seek(stat.Size()+512, 0)
	var endBlocks []byte
	for i := 0; i < (15 * 512); i++ {
		endBlocks = append(endBlocks, 0x00)
	}
	f.Write(endBlocks)

	err = gzipit(path+".tar", destination+".tar.gz")
	if err != nil {
		return err
	}

	defer f.Close()
	defer r.Close()
	return nil
}

func gzipit(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	// target = strings.TrimSuffix(target, ".img")
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = target
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)

	os.Remove(source)

	return err

}
