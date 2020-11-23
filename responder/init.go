package responder

import (
	"github.com/juju/loggo"
	"mega_bot/config"
	"mega_bot/models"
	"sync"
)

var (
	activeResponders []*models.ResponderMatcher
	activeRespondersMutex sync.RWMutex
	logger *loggo.Logger
)

func Init(conf *config.Config, c *chan *models.ResponderRequest) error {
	// Init Logging
	newLogger := loggo.GetLogger("responder")
	logger = &newLogger

	Load()

	for w := 1; w <= conf.ResponderWorkers; w++ {
		go worker(w, c)
	}

	return nil
}

