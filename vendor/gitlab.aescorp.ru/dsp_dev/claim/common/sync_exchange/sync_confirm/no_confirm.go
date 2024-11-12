package sync_confirm

type NoConfirmer struct{}

func NewNoConfirmer(path string) (h Confirmer, err error) {
	block.Lock()
	defer block.Unlock()
	if confirmer == nil {
		var sc NoConfirmer
		confirmer = &sc
	}

	return confirmer, err
}

func (s *NoConfirmer) getIsInited() bool {
	block.RLock()
	defer block.RUnlock()
	return confirmer != nil
}

func (s *NoConfirmer) DeInitConfirm() error {
	block.Lock()
	defer block.Unlock()

	if confirmer == nil {
		return nil
	}

	confirmer = nil
	return nil
}

func (s *NoConfirmer) NewConfirmation(netID string, wait bool) error {
	return nil
}

func (s *NoConfirmer) GetConfirmation(netID string) (*Confirmation, error) {
	return nil, nil
}

func (s *NoConfirmer) MakeConfirmation(netID string, b bool) error {
	return nil
}

func (s *NoConfirmer) SentConfirmation(netID string, b bool) error {
	return nil
}

func (s *NoConfirmer) RecvConfirmation(netID string, b bool) error {
	return nil
}
