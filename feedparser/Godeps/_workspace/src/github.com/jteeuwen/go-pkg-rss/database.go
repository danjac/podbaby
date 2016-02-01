/*
Credits go to github.com/SlyMarbo/rss for inspiring this solution.
*/
package feeder

type database struct {
	request  chan string
	response chan bool
	known    map[string]struct{}
}

func (d *database) Run() {
	d.known = make(map[string]struct{})
	var s string

	for {
		s = <-d.request
		if _, ok := d.known[s]; ok {
			d.response <- true
		} else {
			d.response <- false
			d.known[s] = struct{}{}
		}
	}
}

func NewDatabase() *database {
	database := new(database)
	database.request = make(chan string)
	database.response = make(chan bool)
	go database.Run()
	return database
}
