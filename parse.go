package dicom

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"

	"github.com/suyashkumar/dicom/pkg/charset"
	"github.com/suyashkumar/dicom/pkg/dicomio"
	"github.com/suyashkumar/dicom/pkg/frame"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/suyashkumar/dicom/pkg/uid"
)

const (
	magicWord = "DICM"
)

var (
	ErrorMagicWord              = errors.New("error, DICM magic word not found in correct location")
	ErrorMetaElementGroupLength = errors.New("MetaElementGroupLength tag not found where expected")
)

type Parser interface {
	// Parse DICOM data into a Dataset
	Parse() (Dataset, error)
}

type parser struct {
	reader  dicomio.Reader
	dataset Dataset
	// file is optional, might be populated if reading from an underlying file
	file         *os.File
	frameChannel chan *frame.Frame
	opts         options
}

// NewParser returns a new Parser that points to the provided io.Reader, with bytesToRead bytes left to read. The
// frameChannel is an optional channel (can be nil) upon which DICOM image frames will be sent as they are parsed (if
// provided).
func NewParser(in io.Reader, bytesToRead int64, frameChannel chan *frame.Frame, opts ...Option) (Parser, error) {
	reader, err := dicomio.NewReader(bufio.NewReader(in), binary.LittleEndian, bytesToRead)
	if err != nil {
		return nil, err
	}

	p := parser{
		reader:       reader,
		frameChannel: frameChannel,
	}

	// Apply input options
	for _, opt := range opts {
		opt(&p.opts)
	}

	if !p.opts.assumeNoHeaderAndOffset {
		elems, err := p.readHeader()
		if err != nil {
			return nil, err
		}

		p.dataset = Dataset{Elements: elems}
	}

	return &p, nil
}

// Option represents a parser configuration option that can be passed into a call to NewParser
type Option func(o *options)

// AssumeNoHeaderAndOffset is an Option that configures this parser to assume that the input DICOM has no valid
// header and no DICM magic word after the 128 byte offset.
// This is to support some odd non-conformant dicoms seen in the wild.
var AssumeNoHeaderAndOffset = func(o *options) {
	o.assumeNoHeaderAndOffset = true
}

type options struct {
	assumeNoHeaderAndOffset bool
}

// readHeader reads the DICOM magic header and group two metadata elements.
func (p *parser) readHeader() ([]*Element, error) {
	// Must read as LittleEndian explicit VR
	err := p.reader.Skip(128) // skip preamble
	if err != nil {
		log.Println("skip er")
		return nil, err
	}

	// Check DICOM magic word
	if s, err := p.reader.ReadString(4); err != nil || s != magicWord {
		return nil, ErrorMagicWord
	}

	// Read the length of the metadata elements: (0002,0000) MetaElementGroupLength
	maybeMetaLen, err := readElement(p.reader, nil, nil)
	if err != nil {
		log.Println("read element err")
		return nil, err
	}

	if maybeMetaLen.Tag != tag.FileMetaInformationGroupLength || maybeMetaLen.Value.ValueType() != Ints {
		return nil, ErrorMetaElementGroupLength
	}

	metaLen := maybeMetaLen.Value.GetValue().([]int)[0]

	metaElems := []*Element{maybeMetaLen} // TODO: maybe set capacity to a reasonable initial size

	// Read the metadata elements
	err = p.reader.PushLimit(int64(metaLen))
	if err != nil {
		return nil, err
	}
	defer p.reader.PopLimit()
	for !p.reader.IsLimitExhausted() {
		elem, err := readElement(p.reader, nil, nil)
		if err != nil {
			// TODO: see if we can skip over malformed elements somehow
			log.Println("read element err")

			return nil, err
		}
		// log.Printf("Metadata Element: %s\n", elem)
		metaElems = append(metaElems, elem)
	}
	return metaElems, nil
}

func (p *parser) Parse() (Dataset, error) {
	// Determine and set the transfer syntax based on the metadata elements parsed so far.
	ts, err := p.dataset.FindElementByTag(tag.TransferSyntaxUID)
	if err == nil {
		bo, implicit, err := uid.ParseTransferSyntaxUID(MustGetStrings(ts.Value)[0])
		if err != nil {
			log.Println("WARN: could not parse transfer syntax uid in metadata, proceeding with little endian implicit")
		}
		p.reader.SetTransferSyntax(bo, implicit)
	} else {
		log.Println("WARN: could not parse transfer syntax uid in metadata, proceeding with little endian implicit")
	}
	for !p.reader.IsLimitExhausted() {
		// TODO: avoid silent looping
		elem, err := readElement(p.reader, &p.dataset, p.frameChannel)
		if err != nil {
			// TODO: tolerate some kinds of errors and continue parsing
			return Dataset{}, err
		}

		// log.Println("Read tag: ", elem.Tag)

		// TODO: add dicom options to only keep track of certain tags

		if elem.Tag == tag.SpecificCharacterSet {
			encodingNames := MustGetStrings(elem.Value)
			cs, err := charset.ParseSpecificCharacterSet(encodingNames)
			if err != nil {
				// unable to parse character set, hard error
				// TODO: add option continue, even if unable to parse
				return p.dataset, err
			}
			p.reader.SetCodingSystem(cs)
		}

		p.dataset.Elements = append(p.dataset.Elements, elem)

	}

	if p.frameChannel != nil {
		close(p.frameChannel)
	}
	return p.dataset, nil
}
