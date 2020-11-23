package responder

import (
	"mega_bot/discord"
	"mega_bot/models"
	"strings"
)

func worker(id int, c *chan *models.ResponderRequest) {
	logger.Infof("responder worker %d running", id)
	for req := range *c {
		logger.Debugf("[%d] request %#v", id, req)

		activeRespondersMutex.RLock()

		for _, matcher := range activeResponders {
			if matcher.MatcherRE.FindStringIndex(req.Message) != nil {
				logger.Debugf("[%d] %s matched: %s", id, matcher.MatcherString, req.Message)

				atMe := false
				if strings.Index(req.Message, req.MeString) != -1 {
					atMe = true
				}

				if atMe || req.DirectMessage || matcher.AlwaysRespond {
					switch req.Service {
					case "discord":
						err := discord.SendMessage(matcher.Response, req.ResponseTarget)
						if err != nil {
							logger.Errorf("[%d] couldn't send discord response to %s", id, req.ResponseTarget)
						}
					default:
						logger.Warningf("uknown service type: %s", req.Service)
					}
				}

				activeRespondersMutex.RUnlock()
				break
			}
		}

	}
	logger.Infof("responder worker %d stopping", id)
}
