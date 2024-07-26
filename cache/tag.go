package cache

type Tag struct {
	buckets []*Bucket
}

func (t *Tag) Add(b *Bucket) {
	if t.buckets == nil {
		t.buckets = make([]*Bucket, 0, 4)
	}
	for _, b1 := range t.buckets {
		if b1 == b {
			return
		}
	}
	t.buckets = append(t.buckets, b)
}

func (t *Tag) Delete() {
	for _, b := range t.buckets {
		b.Delete()
	}
}
