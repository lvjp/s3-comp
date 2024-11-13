package client

type Client struct {
	config Config
}

func New(cfg Config) (*Client, error) {
	c := &Client{
		config: cfg,
	}

	c.config.SetDefaults()

	return c, nil
}
