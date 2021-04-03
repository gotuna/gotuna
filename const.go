package gotdd

//func xOptionsWithDefaults() Options {
//	keyPairs := os.Getenv("APP_KEY")
//
//	if options.Session == nil {
//		options.Session = NewSession(sessions.NewCookieStore([]byte(keyPairs)))
//	}
//
//	if options.Router == nil {
//		options.Router = mux.NewRouter()
//	}
//
//	if options.Locale == nil {
//		options.Locale = NewLocale(Translations)
//	}
//
//	if options.FS == nil {
//		options.FS = static.EmbededStatic
//	}
//
//	if options.Logger == nil {
//		options.Logger = log.New(os.Stdout, "", 0)
//	}
//
//	// path prefix for static files
//	// e.g. "/public" or "http://cdn.example.com/assets"
//	options.StaticPrefix = strings.TrimRight(options.StaticPrefix, "/")
//
//	return options
//}
