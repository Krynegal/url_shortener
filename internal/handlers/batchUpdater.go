package handlers

import "github.com/Krynegal/url_shortener.git/internal/storage"

type job struct {
	userID string
	urlID  int
}

type batchUpdater struct {
	deleteURLCh chan job
	buffer      []int
}

func NewBatchUpdater(batchSize int) *batchUpdater {
	return &batchUpdater{
		deleteURLCh: make(chan job, batchSize),
		buffer:      make([]int, 0, batchSize),
	}
}

func (b *batchUpdater) clearBuffer() {
	b.buffer = b.buffer[:cap(b.buffer)]
}

func (b *batchUpdater) deleteQueuedURLs(st storage.Storager) {
	b.clearBuffer()

	db, ok := st.(*storage.DB)
	if !ok {
		return
	}

	go func() {
		for j := range b.deleteURLCh {
			b.buffer = append(b.buffer, j.urlID)
			if len(b.buffer) >= 3 {
				err := db.Delete(j.userID, b.buffer)
				if err != nil {
					return
				}
			}
		}
	}()
}
