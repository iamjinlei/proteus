package markdown

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
				ht.add(1, "", "heading1")
				ht.add(1, "", "heading2")
				ht.add(1, "", "heading3")

				hds := ht.getHeadings()
				require.Equal(t, 3, len(hds))
				require.Equal(t, 1, hds[0].Level)
				require.Equal(t, "heading1", hds[0].Name)
				require.Nil(t, hds[0].Children)
				require.Equal(t, 1, hds[1].Level)
				require.Equal(t, "heading2", hds[1].Name)
				require.Nil(t, hds[1].Children)
				require.Equal(t, 1, hds[2].Level)
				require.Equal(t, "heading3", hds[2].Name)
				require.Nil(t, hds[2].Children)
			},
		},
		{
			name: "nested list",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(1, "", "heading1")
				ht.add(1, "", "heading2")
				ht.add(2, "", "heading2.1")
				ht.add(2, "", "heading2.2")
				ht.add(2, "", "heading2.3")
				ht.add(3, "", "heading2.3.1")
				ht.add(3, "", "heading2.3.2")
				ht.add(1, "", "heading3")
				ht.add(2, "", "heading3.1")
				ht.add(3, "", "heading3.1.1")
				ht.add(3, "", "heading3.1.2")
				ht.add(2, "", "heading3.2")
				ht.add(3, "", "heading3.2.1")
				ht.add(2, "", "heading3.3")

				hds := ht.getHeadings()
				require.Equal(t, 3, len(hds))
				require.Equal(t, 1, hds[0].Level)
				require.Equal(t, "heading1", hds[0].Name)
				require.Nil(t, hds[0].Children)
				require.Equal(t, 1, hds[1].Level)
				require.Equal(t, "heading2", hds[1].Name)
				require.Equal(t, 3, len(hds[1].Children))
				require.Equal(t, 1, hds[2].Level)
				require.Equal(t, "heading3", hds[2].Name)
				require.Equal(t, 3, len(hds[2].Children))

				// 2.x
				require.Equal(t, 2, hds[1].Children[0].Level)
				require.Equal(t, "heading2.1", hds[1].Children[0].Name)
				require.Nil(t, hds[1].Children[0].Children)
				require.Equal(t, 2, hds[1].Children[1].Level)
				require.Equal(t, "heading2.2", hds[1].Children[1].Name)
				require.Nil(t, hds[1].Children[1].Children)
				require.Equal(t, 2, hds[1].Children[2].Level)
				require.Equal(t, "heading2.3", hds[1].Children[2].Name)
				require.Equal(t, 2, len(hds[1].Children[2].Children))

				// 2.3.x
				require.Equal(t, 3, hds[1].Children[2].Children[0].Level)
				require.Equal(t, "heading2.3.1", hds[1].Children[2].Children[0].Name)
				require.Nil(t, hds[1].Children[2].Children[0].Children)
				require.Equal(t, 3, hds[1].Children[2].Children[1].Level)
				require.Equal(t, "heading2.3.2", hds[1].Children[2].Children[1].Name)
				require.Nil(t, hds[1].Children[2].Children[1].Children)

				// 3.x
				require.Equal(t, 2, hds[2].Children[0].Level)
				require.Equal(t, "heading3.1", hds[2].Children[0].Name)
				require.Equal(t, 2, len(hds[2].Children[0].Children))
				require.Equal(t, 2, hds[2].Children[1].Level)
				require.Equal(t, "heading3.2", hds[2].Children[1].Name)
				require.Equal(t, 1, len(hds[2].Children[1].Children))
				require.Equal(t, 2, hds[2].Children[2].Level)
				require.Equal(t, "heading3.3", hds[2].Children[2].Name)
				require.Nil(t, hds[2].Children[2].Children)

				// 3.1.x
				require.Equal(t, 3, hds[2].Children[0].Children[0].Level)
				require.Equal(t, "heading3.1.1", hds[2].Children[0].Children[0].Name)
				require.Nil(t, hds[2].Children[0].Children[0].Children)
				require.Equal(t, 3, hds[2].Children[0].Children[1].Level)
				require.Equal(t, "heading3.1.2", hds[2].Children[0].Children[1].Name)
				require.Nil(t, hds[2].Children[0].Children[1].Children)
				// 3.2.x
				require.Equal(t, 3, hds[2].Children[1].Children[0].Level)
				require.Equal(t, "heading3.2.1", hds[2].Children[1].Children[0].Name)
				require.Nil(t, hds[2].Children[1].Children[0].Children)
			},
		},
		{
			name: "skipped higher headings",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(3, "", "heading2.1.1")
				ht.add(3, "", "heading2.1.2")
				ht.add(4, "", "heading2.1.2.1")

				hds := ht.getHeadings()
				require.Equal(t, 2, len(hds))
				require.Equal(t, 3, hds[0].Level)
				require.Equal(t, "heading2.1.1", hds[0].Name)
				require.Nil(t, hds[0].Children)
				require.Equal(t, 3, hds[1].Level)
				require.Equal(t, "heading2.1.2", hds[1].Name)
				require.Equal(t, 1, len(hds[1].Children))
				require.Equal(t, 4, hds[1].Children[0].Level)
				require.Equal(t, "heading2.1.2.1", hds[1].Children[0].Name)
				require.Nil(t, hds[1].Children[0].Children)
			},
		},
		{
			name: "missing higher headings",
			run: func(t *testing.T) {
				ht := newHeadingTracker()
				ht.add(3, "", "heading3.1.1")
				ht.add(3, "", "heading3.1.2")
				ht.add(2, "", "heading3.2")
				ht.add(2, "", "heading3.3")
				ht.add(1, "", "heading4")

				hds := ht.getHeadings()
				require.Equal(t, 2, len(hds))
				require.Equal(t, 1, hds[0].Level)
				require.Equal(t, "", hds[0].Name)
				require.Equal(t, 3, len(hds[0].Children))
				require.Equal(t, 1, hds[1].Level)
				require.Equal(t, "heading4", hds[1].Name)
				require.Nil(t, hds[1].Children)

				require.Equal(t, 2, hds[0].Children[0].Level)
				require.Equal(t, "", hds[0].Children[0].Name)
				require.Equal(t, 2, len(hds[0].Children[0].Children))
				require.Equal(t, 2, hds[0].Children[1].Level)
				require.Equal(t, "heading3.2", hds[0].Children[1].Name)
				require.Nil(t, hds[0].Children[1].Children)
				require.Equal(t, 2, hds[0].Children[2].Level)
				require.Equal(t, "heading3.3", hds[0].Children[2].Name)
				require.Nil(t, hds[0].Children[2].Children)

				require.Equal(t, 3, hds[0].Children[0].Children[0].Level)
				require.Equal(t, "heading3.1.1", hds[0].Children[0].Children[0].Name)
				require.Nil(t, hds[0].Children[0].Children[0].Children)
				require.Equal(t, 3, hds[0].Children[0].Children[1].Level)
				require.Equal(t, "heading3.1.2", hds[0].Children[0].Children[1].Name)
				require.Nil(t, hds[0].Children[0].Children[1].Children)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
