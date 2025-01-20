package rustmaps

func (g *Generator) SetApiKey(apiKey string) {
	g.config.APIKey = apiKey
	g.rmcli.SetApiKey(apiKey)
}

func (g *Generator) SetTier(tier string) {
	g.config.Tier = tier
}
