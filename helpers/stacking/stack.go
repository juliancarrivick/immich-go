package stacking

import (
	"path"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/simulot/immich-go/helpers/fshelper"
	"github.com/simulot/immich-go/helpers/gen"
	"github.com/simulot/immich-go/immich"
)

type Key struct {
	date     time.Time // time rounded at 5 min
	baseName string    // stack group
}

type Stack struct {
	CoverID string
	IDs     []string
	Date    time.Time
	Names   []string
}

type StackBuilder struct {
	dateRange immich.DateRange // Set capture date range
	stacks    map[Key]Stack
}

func NewStackBuilder() *StackBuilder {
	sb := StackBuilder{
		stacks: map[Key]Stack{},
	}
	sb.dateRange.Set("1850-01-04,2030-01-01")

	return &sb

}

func (sb *StackBuilder) ProcessAsset(ID string, fileName string, captureDate time.Time) {
	if !sb.dateRange.InRange(captureDate) {
		return
	}
	ext := path.Ext(fileName)
	base := strings.TrimSuffix(path.Base(fileName), ext)
	ext = strings.ToLower(ext)

	idx := strings.Index(base, "_BURST")
	if idx > 1 {
		base = base[:idx]
	}

	k := Key{
		date:     captureDate.Round(time.Minute),
		baseName: base,
	}
	s, ok := sb.stacks[k]
	if !ok {
		s.CoverID = ID
		s.Date = captureDate
	}
	s.IDs = append(s.IDs, ID)
	s.Names = append(s.Names, path.Base(fileName))
	if (idx == -1 && slices.Contains([]string{".jpeg", ".jpg", ".jpe"}, ext)) || strings.Contains(fileName, "COVER.") {
		s.CoverID = ID
	}
	sb.stacks[k] = s
}

func (sb *StackBuilder) Stacks() []Stack {
	keys := gen.MapFilterKeys(sb.stacks, func(i Stack) bool {
		return len(i.IDs) > 1
	})

	var stacks []Stack
	for _, k := range keys {
		s := sb.stacks[k]

		// Exclude live photos
		hasPhoto := 0
		hasVideo := 0

		for _, n := range s.Names {
			mime, err := fshelper.MimeFromExt(path.Ext(n))
			if err != nil {
				continue
			}
			switch mime[0] {
			case "video":
				hasVideo++
			case "photo":
				hasPhoto++
			}
		}

		if hasPhoto == 1 && hasVideo == 1 {
			// oh, a live photo!
			continue
		}

		ids := gen.Filter(s.IDs, func(id string) bool {
			return id != s.CoverID
		})
		s.IDs = ids
		stacks = append(stacks, s)

	}
	sort.Slice(stacks, func(i, j int) bool {
		c := stacks[i].Date.Compare(stacks[j].Date)
		switch c {
		case -1:
			return true
		case +1:
			return false
		}
		c = strings.Compare(stacks[i].Names[0], stacks[j].Names[0])
		switch c {
		case -1:
			return true
		}
		return false
	})
	return stacks
}
