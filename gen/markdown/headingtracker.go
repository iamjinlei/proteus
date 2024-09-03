package markdown

type headingTracker struct {
	queue [][]*Heading
}

func newHeadingTracker() *headingTracker {
	return &headingTracker{}
}

func (t *headingTracker) add(level int, id, name string) {
	h := &Heading{
		Level: level,
		ID:    id,
		Name:  name,
	}

	if len(t.queue) == 0 {
		for i := 1; i < level; i++ {
			t.queue = append(t.queue,
				[]*Heading{
					&Heading{
						Level: i,
					},
				},
			)
		}
		// The first heading inserted.
		t.queue = append(t.queue, []*Heading{h})
		return
	}

	currList := t.queue[len(t.queue)-1]
	lastHeading := currList[len(currList)-1]
	if level == lastHeading.Level {
		// Append to the current heading list.
		t.queue[len(t.queue)-1] = append(currList, h)
	} else if level > lastHeading.Level {
		// Go to lower levels, fill empty middle levels if needed.
		for i := lastHeading.Level + 1; i < level; i++ {
			t.queue = append(t.queue,
				[]*Heading{
					&Heading{
						Level: i,
					},
				},
			)
		}
		t.queue = append(t.queue, []*Heading{h})
	} else {
		// Go back to higher levels, attach lower level list to parent.
		for i := lastHeading.Level; i > level; i-- {
			list := t.queue[len(t.queue)-1]
			t.queue = t.queue[:len(t.queue)-1]
			parentList := t.queue[len(t.queue)-1]
			parentList[len(parentList)-1].Children = list
		}
		currList := t.queue[len(t.queue)-1]
		t.queue[len(t.queue)-1] = append(currList, h)
	}
}

func (t *headingTracker) getHeadings() []*Heading {
	// Trim empty headings if the top heading level is > 1
	idx := 0
	for idx < len(t.queue) {
		if len(t.queue[idx]) > 1 || t.queue[idx][0].Name != "" {
			break
		}
		idx++
	}
	if idx > 0 {
		t.queue = t.queue[idx:]
	}

	// Collapse remaining heading nests
	for i := len(t.queue) - 1; i > 0; i-- {
		parentList := t.queue[i-1]
		parentList[len(parentList)-1].Children = t.queue[i]
	}

	return t.queue[0]
}
