package api

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Post("/", s.Mailer)
	s.Router.Post("/send", s.SendMail)
}
