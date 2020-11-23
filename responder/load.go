package responder

import "mega_bot/models"

func Load() {
	rms, err := models.ReadEnabledResponderMatchers()
	if err != nil {
		logger.Errorf("could not load responder matchers: %s", err.Error())
	}

	activeRespondersMutex.Lock()
	defer activeRespondersMutex.Unlock()

	activeResponders = rms
}
