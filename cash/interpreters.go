package cash

import "github.com/shahinrahimi/teletradebot/models"

func (c *Cash) GetInterpreter(ID int64) (*models.Interpreter, bool) {
	mu.RLock()
	defer mu.RUnlock()
	d, exist := c.interpreters[ID]
	return d, exist
}

func (c *Cash) SetInterpreter(i *models.Interpreter, ID int64) {
	mu.Lock()
	defer mu.Unlock()
	c.interpreters[ID] = i
}

func (c *Cash) RemoveInterpreter(ID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(c.interpreters, ID)
}
