package flacvorbis

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	httpclient "github.com/ddliu/go-httpclient"
	flac "github.com/go-flac/go-flac"
)

func TestNewVorbisComment(t *testing.T) {
	cmt := New()
	if !strings.HasPrefix(cmt.Vendor, "flacvorbis") {
		t.Errorf("Unexpected vendor: %s\n", cmt.Vendor)
		t.Fail()
	}
	if len(cmt.Comments) != 0 {
		t.Error("Unexpected comments: ", cmt.Comments)
		t.Fail()
	}
}
func TestVorbisFromExistingFlac(t *testing.T) {
	zipres, err := httpclient.Begin().Get("http://helpguide.sony.net/high-res/sample1/v1/data/Sample_BeeMoved_96kHz24bit.flac.zip")
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	zipdata, err := zipres.ReadAll()
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	zipfile, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		t.Errorf("Error while decompressing test file: %s", err.Error())
		t.FailNow()
	}
	if zipfile.File[0].Name != "Sample_BeeMoved_96kHz24bit.flac" {
		t.Errorf("Unexpected test file content: %s", zipfile.File[0].Name)
		t.FailNow()
	}

	flachandle, err := zipfile.File[0].Open()
	if err != nil {
		t.Errorf("Failed to decompress test file: %s", err)
		t.FailNow()
	}

	f, err := flac.ParseBytes(flachandle)
	if err != nil {
		t.Errorf("Failed to parse flac file: %s", err)
		t.FailNow()
	}

	var cmt *MetaDataBlockVorbisComment
	for _, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			cmt, err = ParseFromMetaDataBlock(*meta)
			if err != nil {
				t.Errorf("Error while parsing metadata image: %s\n", err.Error())
				t.Fail()
			}
		}
	}

	if err := cmt.Add(FIELD_GENRE, "Bee Pop"); err != nil {
		t.Error(err)
		t.Fail()
	}

	check := func(cmt *MetaDataBlockVorbisComment) {
		if cmt.Vendor != "reference libFLAC 1.2.1 win64 20080709" {
			t.Errorf("Unexpected vendor string: %s\n", cmt.Vendor)
			t.Fail()
		}
		if res, err := cmt.Get(FIELD_ALBUM); err != nil {
			t.Error(err)
			t.Fail()
		} else if len(res) != 1 || res[0] != "Bee Moved" {
			t.Error("Unexpected album name: ", res)
			t.Fail()
		}

		if res, err := cmt.Get(FIELD_ARTIST); err != nil {
			t.Error(err)
			t.Fail()
		} else if len(res) != 1 || res[0] != "Blue Monday FM" {
			t.Error("Unexpected artist name: ", res)
			t.Fail()
		}

		if res, err := cmt.Get(FIELD_TITLE); err != nil {
			t.Error(err)
			t.Fail()
		} else if len(res) != 1 || res[0] != "Bee Moved" {
			t.Error("Unexpected title name: ", res)
			t.Fail()
		}

		if res, err := cmt.Get(FIELD_GENRE); err != nil {
			t.Error(err)
			t.Fail()
		} else if len(res) != 1 || res[0] != "Bee Pop" {
			t.Error("Unexpected title name: ", res)
			t.Fail()
		}
	}
	check(cmt)
	new, err := ParseFromMetaDataBlock(cmt.Marshal())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	check(new)

}
