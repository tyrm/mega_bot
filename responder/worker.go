package responder

import (
	"mega_bot/discord"
	"mega_bot/models"
	"strings"
)

func worker(id int, c *chan *models.ResponderRequest) {
	logger.Infof("responder worker %d running", id)
	for req := range *c {
		logger.Tracef("[%d] request %#v", id, req)

		activeRespondersMutex.RLock()
		for _, matcher := range activeResponders {
			if matcher.MatcherRE.FindStringIndex(req.Message) != nil {
				logger.Debugf("[%d] matched message from %s(%s). responder: %s", id, req.Service, req.ResponseTarget, matcher.ID)

				atMe := false
				if strings.Index(req.Message, req.MeString) != -1 {
					atMe = true
				}

				if atMe || req.DirectMessage || matcher.AlwaysRespond {
					logger.Debugf("[%d] ending response to %s(%s). responder: %s", id, req.Service, req.ResponseTarget, matcher.ID)
					switch req.Service {
					case "discord":
						err := discord.SendMessage(matcher.Response, req.ResponseTarget)
						if err != nil {
							logger.Errorf("[%d] couldn't send discord response to %s(%s)", id, req.Service, req.ResponseTarget)
						}
					default:
						logger.Warningf("[%d] unknown provider: %s", id, req.Service)
					}
				}

				break
			}
		}

	}
	activeRespondersMutex.RUnlock()

	logger.Infof("responder worker %d stopping", id)
}
