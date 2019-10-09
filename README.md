# flacvorbis

[![Documentation](https://godoc.org/github.com/go-flac/flacvorbis?status.svg)](https://godoc.org/github.com/go-flac/flacvorbis)
[![Build Status](https://travis-ci.org/go-flac/flacvorbis.svg?branch=master)](https://travis-ci.org/go-flac/flacvorbis)
[![Coverage Status](https://coveralls.io/repos/github/go-flac/flacvorbis/badge.svg?branch=master)](https://coveralls.io/github/go-flac/flacvorbis?branch=master)

FLAC vorbis comment metablock manipulation for [go-flac](https://www.github.com/go-flac/go-flac)

## Examples

The following example adds a title to the FLAC metadata. Only add a new entry if you are sure there is none existing, otherwise will mislead some audio players.

```golang
package example

import (
    "github.com/go-flac/flacvorbis"
    "github.com/go-flac/go-flac"
)

func addFLACTitle(fileName string, title []byte) {
	f, err := flac.ParseFile(fileName)
	if err != nil {
		panic(err)
	}
	cmts := flacvorbis.New()
	cmts.Add(flacvorbis.FIELD_TITLE, title)
	cmtsmeta := cmts.Marshal()
	f.Meta = append(f.Meta, &cmtsmeta)
	f.Save(fileName)
}
```

The following example checks for an existing vorbis comment block and updates it with a title to the FLAC metadata.

```golang
package example

import (
    "github.com/go-flac/flacvorbis"
    "github.com/go-flac/go-flac"
)

func updateFLACTitle(fileName string, title []byte) {
	f, err := flac.ParseFile(fileName)
	if err != nil {
		panic(err)
	}
	cmts, idx := extractFLACComment(f)
	if err != nil {
		panic err
	}
	cmts.Add(flacvorbis.FIELD_TITLE, title)
	cmtsmeta := cmts.Marshal()
	// Replace existing Metadata with the modified one
	f.Meta[idx] = &cmtsmeta
	f.Save(fileName)
}
```

The following example extracts existing tags from a FLAC file. It returnt the last vorbis comment block and also the corresponding index of the metadata, which could be used for updating later on.
```golang
package example

import (
    "github.com/go-flac/flacvorbis"
    "github.com/go-flac/go-flac"
)

func extractFLACComment(fileName string) (*flacvorbis.MetadataBlockVorbisComment, int) {
	f, err := flac.ParseFile(fileName)
	if err != nil {
		panic(err)
	}
    
	var cmt *flacvorbis.MetadataBlockVorbisComment
	var cmtIdx int
	for idx, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			cmt, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			cmtIdx = idx
			if err != nil {
				panic(err)
			}
		}
    	}
	return cmt, cmtIdx
}
```
