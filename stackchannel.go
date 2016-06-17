package crawler


// see http://gowithconfidence.tumblr.com/post/31426832143/stacked-channels
type RequestChan chan []*UrlRequest

func NewRequestChan() RequestChan {
	// The pop channel is stacked, so only a buffer of 1 is required
	return make(chan []*UrlRequest, 1)
}

// do the pop action by <-RequestChan
// The stack function ensures the specified URLs are added to the channel
// with minimal blocking (since the channel is stacked, it is virtually equivalent
// to an infinitely buffered channel).
func (pc RequestChan) push(cmd ...*UrlRequest) {
	toStack := cmd
	for {
		select {
		case pc <- toStack:
			return
		case old := <-pc:
		// Content of the channel got emptied and is now in old, so append whatever
		// is in toStack to it, so that it can either be inserted in the channel,
		// or appended to some other content that got through in the meantime.
			toStack = append(old, toStack...)
		}
	}
}
//
type ResponseChan chan []*UrlResponse
func NewResponseChan() ResponseChan {
	return make(chan []*UrlResponse, 1)
}
func (pc ResponseChan) push(cmd ...*UrlResponse) {
	toStack := cmd
	for {
		select {
		case pc <- toStack:
			return
		case old := <-pc:
			toStack = append(old, toStack...)
		}
	}
}

