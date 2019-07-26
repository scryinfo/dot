// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
)

//FindNewer find newer
func FindNewer(lives []*dot.TypeLives, liveId dot.LiveId) dot.Newer {
	for _, it := range lives {
		for _, lid := range it.Lives {
			if lid.LiveId == liveId {
				return it.Meta.NewDot
			}
		}
	}

	for _, it := range lives {
		if it.Meta.TypeId == dot.TypeId(liveId) {
			return it.Meta.NewDot
		}
	}
	return nil
}
