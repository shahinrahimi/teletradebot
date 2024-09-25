package cash

import "github.com/shahinrahimi/teletradebot/models"

func (c *Cash) GetDescribers() []*models.Describer {
	mu.RLock()
	defer mu.RUnlock()
	var ds []*models.Describer
	for _, d := range c.describers {
		ds = append(ds, d)
	}
	return ds
}

func (c *Cash) GetDescriber(ID int64) (*models.Describer, bool) {
	mu.RLock()
	defer mu.RUnlock()
	d, exist := c.describers[ID]
	return d, exist
}

func (c *Cash) SetDescriber(i *models.Describer, ID int64) {
	mu.Lock()
	defer mu.Unlock()
	c.describers[ID] = i
}

func (c *Cash) RemoveDescriber(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(c.describers, ID)
}
