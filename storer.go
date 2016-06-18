package eago

type Storer interface {
	Store(resp *UrlResponse)
}

type DefaultStore struct {
}

func (d *DefaultStore) Store(resp *UrlResponse) {

	return
}
