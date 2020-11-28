package responder

import "mega_bot/models"

func Load() error {
	rms, err := models.ReadEnabledResponderMatchers()
	if err != nil {
		logger.Errorf("could not load responder matchers: %s", err.Error())
		return err
	}

	rmMap := make(map[string]*models.ResponderMatcher)
	for _, rm := range *rms {
		newRM := rm
		rmMap[rm.ID] = &newRM
	}

	activeRespondersMutex.Lock()
	defer activeRespondersMutex.Unlock()

	activeResponders = rmMap
	return nil
}
