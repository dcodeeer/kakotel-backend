package kuro

func (s *Server) SetBeforeUpgrade(callback BeforeUpgradeFunc) {
	s.beforeUpgrade = callback
}
