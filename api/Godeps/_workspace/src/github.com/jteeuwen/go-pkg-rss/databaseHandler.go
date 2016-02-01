package feeder

type databaseHandler struct {
	db          *database
	itemhandler ItemHandler
	chanhandler ChannelHandler
}

func (d *databaseHandler) ProcessItems(f *Feed, ch *Channel, items []*Item) {
	var newitems []*Item
	for _, item := range items {
		if d.db.request <- item.Key(); !<-d.db.response {
			newitems = append(newitems, item)
		}
	}
	if len(newitems) > 0 && d.itemhandler != nil {
		d.itemhandler.ProcessItems(f, ch, newitems)
	}
}

func (d *databaseHandler) ProcessChannels(f *Feed, ch []*Channel) {
	var newchannels []*Channel
	for _, channel := range ch {
		if d.db.request <- channel.Key(); !<-d.db.response {
			newchannels = append(newchannels, channel)
		}
	}
	if len(newchannels) > 0 && d.chanhandler != nil {
		d.chanhandler.ProcessChannels(f, newchannels)
	}
}

func NewDatabaseItemHandler(db *database, itemhandler ItemHandler) ItemHandler {
	database := new(databaseHandler)
	database.db = db
	database.itemhandler = itemhandler
	return database
}

func NewDatabaseChannelHandler(db *database, chanhandler ChannelHandler) ChannelHandler {
	database := new(databaseHandler)
	database.db = db
	database.chanhandler = chanhandler
	return database
}
