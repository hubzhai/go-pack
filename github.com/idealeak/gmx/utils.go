package gmx

func checkWhiteList(ip string) bool {
	if len(Cfg.WhiteAddr) > 0 {
		for _, value := range Cfg.WhiteAddr {
			if value == ip {
				return true
			}
		}
		return false
	}
	return true
}

func checkAuth(token, ip string) bool {
	if token == "" {
		return false
	}

	s := GetSession(token)
	if s == nil {
		return false
	}

	if ipv, ok := s.Get(SESSION_AUTH_IP); ok {
		if authIp, ok := ipv.(string); ok {
			if authIp != ip {
				return false
			}
		}
	}
	return true
}
