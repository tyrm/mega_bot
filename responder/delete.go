package responder

func Delete(id string) {
	activeRespondersMutex.Lock()
	defer activeRespondersMutex.Unlock()

	_, ok := activeResponders[id];
	if ok {
		delete(activeResponders, id)
	}

}