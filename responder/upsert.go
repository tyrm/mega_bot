package responder

import "mega_bot/models"

func Upsert(rm *models.ResponderMatcher) {
	activeRespondersMutex.Lock()
	defer activeRespondersMutex.Unlock()

	activeResponders[rm.ID] = rm
}