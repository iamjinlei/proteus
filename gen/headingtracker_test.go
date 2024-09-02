package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeadingTracker(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "flat list",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(1, "heading1")
				ht.add(1, "heading2")
				ht.add(1, "heading3")

				hds := ht.getHeadings()
				require.Equal(t, 3, len(hds))
				require.Equal(t, 1, hds[0].level)
				require.Equal(t, "heading1", hds[0].name)
				require.Nil(t, hds[0].children)
				require.Equal(t, 1, hds[1].level)
				require.Equal(t, "heading2", hds[1].name)
				require.Nil(t, hds[1].children)
				require.Equal(t, 1, hds[2].level)
				require.Equal(t, "heading3", hds[2].name)
				require.Nil(t, hds[2].children)
			},
		},
		{
			name: "nested list",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(1, "heading1")
				ht.add(1, "heading2")
				ht.add(2, "heading2.1")
				ht.add(2, "heading2.2")
				ht.add(2, "heading2.3")
				ht.add(3, "heading2.3.1")
				ht.add(3, "heading2.3.2")
				ht.add(1, "heading3")
				ht.add(2, "heading3.1")
				ht.add(3, "heading3.1.1")
				ht.add(3, "heading3.1.2")
				ht.add(2, "heading3.2")
				ht.add(3, "heading3.2.1")
				ht.add(2, "heading3.3")

				hds := ht.getHeadings()
				require.Equal(t, 3, len(hds))
				require.Equal(t, 1, hds[0].level)
				require.Equal(t, "heading1", hds[0].name)
				require.Nil(t, hds[0].children)
				require.Equal(t, 1, hds[1].level)
				require.Equal(t, "heading2", hds[1].name)
				require.Equal(t, 3, len(hds[1].children))
				require.Equal(t, 1, hds[2].level)
				require.Equal(t, "heading3", hds[2].name)
				require.Equal(t, 3, len(hds[2].children))

				// 2.x
				require.Equal(t, 2, hds[1].children[0].level)
				require.Equal(t, "heading2.1", hds[1].children[0].name)
				require.Nil(t, hds[1].children[0].children)
				require.Equal(t, 2, hds[1].children[1].level)
				require.Equal(t, "heading2.2", hds[1].children[1].name)
				require.Nil(t, hds[1].children[1].children)
				require.Equal(t, 2, hds[1].children[2].level)
				require.Equal(t, "heading2.3", hds[1].children[2].name)
				require.Equal(t, 2, len(hds[1].children[2].children))

				// 2.3.x
				require.Equal(t, 3, hds[1].children[2].children[0].level)
				require.Equal(t, "heading2.3.1", hds[1].children[2].children[0].name)
				require.Nil(t, hds[1].children[2].children[0].children)
				require.Equal(t, 3, hds[1].children[2].children[1].level)
				require.Equal(t, "heading2.3.2", hds[1].children[2].children[1].name)
				require.Nil(t, hds[1].children[2].children[1].children)

				// 3.x
				require.Equal(t, 2, hds[2].children[0].level)
				require.Equal(t, "heading3.1", hds[2].children[0].name)
				require.Equal(t, 2, len(hds[2].children[0].children))
				require.Equal(t, 2, hds[2].children[1].level)
				require.Equal(t, "heading3.2", hds[2].children[1].name)
				require.Equal(t, 1, len(hds[2].children[1].children))
				require.Equal(t, 2, hds[2].children[2].level)
				require.Equal(t, "heading3.3", hds[2].children[2].name)
				require.Nil(t, hds[2].children[2].children)

				// 3.1.x
				require.Equal(t, 3, hds[2].children[0].children[0].level)
				require.Equal(t, "heading3.1.1", hds[2].children[0].children[0].name)
				require.Nil(t, hds[2].children[0].children[0].children)
				require.Equal(t, 3, hds[2].children[0].children[1].level)
				require.Equal(t, "heading3.1.2", hds[2].children[0].children[1].name)
				require.Nil(t, hds[2].children[0].children[1].children)
				// 3.2.x
				require.Equal(t, 3, hds[2].children[1].children[0].level)
				require.Equal(t, "heading3.2.1", hds[2].children[1].children[0].name)
				require.Nil(t, hds[2].children[1].children[0].children)
			},
		},
		{
			name: "skipped higher headings",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(3, "heading2.1.1")
				ht.add(3, "heading2.1.2")
				ht.add(4, "heading2.1.2.1")

				hds := ht.getHeadings()
				require.Equal(t, 2, len(hds))
				require.Equal(t, 3, hds[0].level)
				require.Equal(t, "heading2.1.1", hds[0].name)
				require.Nil(t, hds[0].children)
				require.Equal(t, 3, hds[1].level)
				require.Equal(t, "heading2.1.2", hds[1].name)
				require.Equal(t, 1, len(hds[1].children))
				require.Equal(t, 4, hds[1].children[0].level)
				require.Equal(t, "heading2.1.2.1", hds[1].children[0].name)
				require.Nil(t, hds[1].children[0].children)
			},
		},
		{
			name: "missing higher headings",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(3, "heading3.1.1")
				ht.add(3, "heading3.1.2")
				ht.add(2, "heading3.2")
				ht.add(2, "heading3.3")
				ht.add(1, "heading4")

				hds := ht.getHeadings()
				require.Equal(t, 2, len(hds))
				require.Equal(t, 1, hds[0].level)
				require.Equal(t, "", hds[0].name)
				require.Equal(t, 3, len(hds[0].children))
				require.Equal(t, 1, hds[1].level)
				require.Equal(t, "heading4", hds[1].name)
				require.Nil(t, hds[1].children)

				require.Equal(t, 2, hds[0].children[0].level)
				require.Equal(t, "", hds[0].children[0].name)
				require.Equal(t, 2, len(hds[0].children[0].children))
				require.Equal(t, 2, hds[0].children[1].level)
				require.Equal(t, "heading3.2", hds[0].children[1].name)
				require.Nil(t, hds[0].children[1].children)
				require.Equal(t, 2, hds[0].children[2].level)
				require.Equal(t, "heading3.3", hds[0].children[2].name)
				require.Nil(t, hds[0].children[2].children)

				require.Equal(t, 3, hds[0].children[0].children[0].level)
				require.Equal(t, "heading3.1.1", hds[0].children[0].children[0].name)
				require.Nil(t, hds[0].children[0].children[0].children)
				require.Equal(t, 3, hds[0].children[0].children[1].level)
				require.Equal(t, "heading3.1.2", hds[0].children[0].children[1].name)
				require.Nil(t, hds[0].children[0].children[1].children)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
