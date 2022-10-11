package websvcclient

type Client struct{}

func (c *Client) GetUserRepo(username *string) string {
	return "IgorTkachuk/cartridge_accounting"
}
