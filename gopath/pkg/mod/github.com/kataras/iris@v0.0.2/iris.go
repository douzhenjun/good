package iris

import (
	stdContext "context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/errgroup"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris/core/netutil"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/i18n"
	requestLogger "github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/view"
	"github.com/kataras/pio"

	"github.com/kataras/golog"
	"github.com/kataras/tunnel"
)

// Version is the current version number of the Iris Web Framework.
const Version = "stale"

func init() {
	fmt.Println(`You have installed an invalid version. Install with:
	go get -u github.com/kataras/iris/v12@latest

	If your Open Source project depends on that pre-go1.9 version please open an issue
	at https://github.com/kataras/iris/issues/new and share your repository with us,
	we will upgrade your project's code base to the latest version for free.

	If you have a commercial project that you cannot share publically, please contact with
	@kataras at https://chat.iris-go.com. Assistance will be provided to you and your colleagues
	for free.
	`)

	fmt.Print("Run ")
	pio.WriteRich(os.Stdout, "autofix", pio.Green, pio.Underline)
	fmt.Print("? (Y/n): ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		golog.Fatalf("can not take input from user: %v", err)
	}
	input = strings.ToLower(input)
	if input == "" || input == "y" {
		err := tryFix()
		if err != nil {
			golog.Fatalf("autofix: %v", err)
		}

		golog.Infof("OK. Restart the application manually now.")
		os.Exit(0)
	} else {
		os.Exit(-1)
	}
}

// Byte unit helpers.
const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)

// Application is responsible to manage the state of the application.
// It contains and handles all the necessary parts to create a fast web server.
type Application struct {
	// routing embedded | exposing APIBuilder's and Router's public API.
	*router.APIBuilder
	*router.Router
	router.HTTPErrorHandler // if Router is Downgraded this is nil.
	ContextPool             *context.Pool

	// config contains the configuration fields
	// all fields defaults to something that is working, developers don't have to set it.
	config *Configuration

	// the golog logger instance, defaults to "Info" level messages (all except "Debug")
	logger *golog.Logger

	// I18n contains localization and internationalization support.
	// Use the `Load` or `LoadAssets` to locale language files.
	//
	// See `Context#Tr` method for request-based translations.
	I18n *i18n.I18n

	// Validator is the request body validator, defaults to nil.
	Validator context.Validator

	// view engine
	view view.View
	// used for build
	builded     bool
	defaultMode bool

	mu sync.Mutex
	// Hosts contains a list of all servers (Host Supervisors) that this app is running on.
	//
	// Hosts may be empty only if application ran(`app.Run`) with `iris.Raw` option runner,
	// otherwise it contains a single host (`app.Hosts[0]`).
	//
	// Additional Host Supervisors can be added to that list by calling the `app.NewHost` manually.
	//
	// Hosts field is available after `Run` or `NewHost`.
	Hosts             []*host.Supervisor
	hostConfigurators []host.Configurator
}

// New creates and returns a fresh empty iris *Application instance.
func New() *Application {
	config := DefaultConfiguration()

	app := &Application{
		config:     &config,
		logger:     golog.Default,
		I18n:       i18n.New(),
		APIBuilder: router.NewAPIBuilder(),
		Router:     router.NewRouter(),
	}

	app.ContextPool = context.New(func() interface{} {
		return context.NewContext(app)
	})

	return app
}

// Default returns a new Application instance which on build state registers
// html view engine on "./views" and load locales from "./locales/*/*".
// The return instance recovers on panics and logs the incoming http requests too.
func Default() *Application {
	app := New()
	app.Use(recover.New())
	app.Use(requestLogger.New())
	app.Use(Compression)

	app.defaultMode = true

	return app
}

// WWW creates and returns a "www." subdomain.
// The difference from `app.Subdomain("www")` or `app.Party("www.")` is that the `app.WWW()` method
// wraps the router so all http(s)://mydomain.com will be redirect to http(s)://www.mydomain.com.
// Other subdomains can be registered using the app: `sub := app.Subdomain("mysubdomain")`,
// child subdomains can be registered using the www := app.WWW(); www.Subdomain("wwwchildSubdomain").
func (app *Application) WWW() router.Party {
	return app.SubdomainRedirect(app, app.Subdomain("www"))
}

// SubdomainRedirect registers a router wrapper which
// redirects(StatusMovedPermanently) a (sub)domain to another subdomain or to the root domain as fast as possible,
// before the router's try to execute route's handler(s).
//
// It receives two arguments, they are the from and to/target locations,
// 'from' can be a wildcard subdomain as well (app.WildcardSubdomain())
// 'to' is not allowed to be a wildcard for obvious reasons,
// 'from' can be the root domain(app) when the 'to' is not the root domain and visa-versa.
//
// Usage:
// www := app.Subdomain("www") <- same as app.Party("www.")
// app.SubdomainRedirect(app, www)
// This will redirect all http(s)://mydomain.com/%anypath% to http(s)://www.mydomain.com/%anypath%.
//
// One or more subdomain redirects can be used to the same app instance.
//
// If you need more information about this implementation then you have to navigate through
// the `core/router#NewSubdomainRedirectWrapper` function instead.
//
// Example: https://github.com/kataras/iris/tree/master/_examples/routing/subdomains/redirect
func (app *Application) SubdomainRedirect(from, to router.Party) router.Party {
	sd := router.NewSubdomainRedirectWrapper(app.ConfigurationReadOnly().GetVHost, from.GetRelPath(), to.GetRelPath())
	app.Router.WrapRouter(sd)
	return to
}

// Configure can called when modifications to the framework instance needed.
// It accepts the framework instance
// and returns an error which if it's not nil it's printed to the logger.
// See configuration.go for more.
//
// Returns itself in order to be used like `app:= New().Configure(...)`
func (app *Application) Configure(configurators ...Configurator) *Application {
	for _, cfg := range configurators {
		if cfg != nil {
			cfg(app)
		}
	}

	return app
}

// ConfigurationReadOnly returns an object which doesn't allow field writing.
func (app *Application) ConfigurationReadOnly() context.ConfigurationReadOnly {
	return app.config
}

// Logger returns the golog logger instance(pointer) that is being used inside the "app".
//
// Available levels:
// - "disable"
// - "fatal"
// - "error"
// - "warn"
// - "info"
// - "debug"
// Usage: app.Logger().SetLevel("error")
// Or set the level through Configurartion's LogLevel or WithLogLevel functional option.
// Defaults to "info" level.
//
// Callers can use the application's logger which is
// the same `golog.Default` logger,
// to print custom logs too.
// Usage:
// app.Logger().Error/Errorf("...")
// app.Logger().Warn/Warnf("...")
// app.Logger().Info/Infof("...")
// app.Logger().Debug/Debugf("...")
//
// Setting one or more outputs: app.Logger().SetOutput(io.Writer...)
// Adding one or more outputs : app.Logger().AddOutput(io.Writer...)
//
// Adding custom levels requires import of the `github.com/kataras/golog` package:
//	First we create our level to a golog.Level
//	in order to be used in the Log functions.
//	var SuccessLevel golog.Level = 6
//	Register our level, just three fields.
//	golog.Levels[SuccessLevel] = &golog.LevelMetadata{
//		Name:    "success",
//		RawText: "[SUCC]",
//		// ColorfulText (Green Color[SUCC])
//		ColorfulText: "\x1b[32m[SUCC]\x1b[0m",
//	}
// Usage:
// app.Logger().SetLevel("success")
// app.Logger().Logf(SuccessLevel, "a custom leveled log message")
func (app *Application) Logger() *golog.Logger {
	return app.logger
}

// I18nReadOnly returns the i18n's read-only features.
// See `I18n` method for more.
func (app *Application) I18nReadOnly() context.I18nReadOnly {
	return app.I18n
}

// Validate validates a value and returns nil if passed or
// the failure reason if does not.
func (app *Application) Validate(v interface{}) error {
	if app.Validator == nil {
		return nil
	}

	// val := reflect.ValueOf(v)
	// if val.Kind() == reflect.Ptr && !val.IsNil() {
	// 	val = val.Elem()
	// }

	// if val.Kind() == reflect.Struct && val.Type() != timeType {
	// 	return app.Validator.Struct(v)
	// }

	// no need to check the kind, underline lib does it but in the future this may change (look above).
	err := app.Validator.Struct(v)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "validator: ") {
			return err
		}
	}

	return nil
}

// RegisterView should be used to register view engines mapping to a root directory
// and the template file(s) extension.
func (app *Application) RegisterView(viewEngine view.Engine) {
	app.view.Register(viewEngine)
}

// View executes and writes the result of a template file to the writer.
//
// First parameter is the writer to write the parsed template.
// Second parameter is the relative, to templates directory, template filename, including extension.
// Third parameter is the layout, can be empty string.
// Forth parameter is the bindable data to the template, can be nil.
//
// Use context.View to render templates to the client instead.
// Returns an error on failure, otherwise nil.
func (app *Application) View(writer io.Writer, filename string, layout string, bindingData interface{}) error {
	if app.view.Len() == 0 {
		err := errors.New("view engine is missing, use `RegisterView`")
		app.logger.Error(err)
		return err
	}

	err := app.view.ExecuteWriter(writer, filename, layout, bindingData)
	if err != nil {
		app.logger.Error(err)
	}
	return err
}

// ConfigureHost accepts one or more `host#Configuration`, these configurators functions
// can access the host created by `app.Run` or `app.Listen`,
// they're being executed when application is ready to being served to the public.
//
// It's an alternative way to interact with a host that is automatically created by
// `app.Run`.
//
// These "configurators" can work side-by-side with the `iris#Addr, iris#Server, iris#TLS, iris#AutoTLS, iris#Listener`
// final arguments("hostConfigs") too.
//
// Note that these application's host "configurators" will be shared with the rest of
// the hosts that this app will may create (using `app.NewHost`), meaning that
// `app.NewHost` will execute these "configurators" everytime that is being called as well.
//
// These "configurators" should be registered before the `app.Run` or `host.Serve/Listen` functions.
func (app *Application) ConfigureHost(configurators ...host.Configurator) *Application {
	app.mu.Lock()
	app.hostConfigurators = append(app.hostConfigurators, configurators...)
	app.mu.Unlock()
	return app
}

// NewHost accepts a standard *http.Server object,
// completes the necessary missing parts of that "srv"
// and returns a new, ready-to-use, host (supervisor).
func (app *Application) NewHost(srv *http.Server) *host.Supervisor {
	app.mu.Lock()
	defer app.mu.Unlock()

	// set the server's handler to the framework's router
	if srv.Handler == nil {
		srv.Handler = app.Router
	}

	// check if different ErrorLog provided, if not bind it with the framework's logger
	if srv.ErrorLog == nil {
		srv.ErrorLog = log.New(app.logger.Printer.Output, "[HTTP Server] ", 0)
	}

	if addr := srv.Addr; addr == "" {
		addr = ":8080"
		if len(app.Hosts) > 0 {
			if v := app.Hosts[0].Server.Addr; v != "" {
				addr = v
			}
		}

		srv.Addr = addr
	}

	// app.logger.Debugf("Host: addr is %s", srv.Addr)

	// create the new host supervisor
	// bind the constructed server and return it
	su := host.New(srv)

	if app.config.vhost == "" { // vhost now is useful for router subdomain on wildcard subdomains,
		// in order to correct decide what to do on:
		// mydomain.com -> invalid
		// localhost -> invalid
		// sub.mydomain.com -> valid
		// sub.localhost -> valid
		// we need the host (without port if 80 or 443) in order to validate these, so:
		app.config.vhost = netutil.ResolveVHost(srv.Addr)
	}

	// app.logger.Debugf("Host: virtual host is %s", app.config.vhost)

	// the below schedules some tasks that will run among the server

	if !app.config.DisableStartupLog {
		// show the available info to exit from app.
		su.RegisterOnServe(host.WriteStartupLogOnServe(app.logger.Printer.Output)) // app.logger.Writer -> Info
		// app.logger.Debugf("Host: register startup notifier")
	}

	if !app.config.DisableInterruptHandler {
		// when CTRL/CMD+C pressed.
		shutdownTimeout := 10 * time.Second
		host.RegisterOnInterrupt(host.ShutdownOnInterrupt(su, shutdownTimeout))
		// app.logger.Debugf("Host: register server shutdown on interrupt(CTRL+C/CMD+C)")
	}

	su.IgnoredErrors = append(su.IgnoredErrors, app.config.IgnoreServerErrors...)
	if len(su.IgnoredErrors) > 0 {
		app.logger.Debugf("Host: server will ignore the following errors: %s", su.IgnoredErrors)
	}

	su.Configure(app.hostConfigurators...)

	app.Hosts = append(app.Hosts, su)

	return su
}

// Shutdown gracefully terminates all the application's server hosts and any tunnels.
// Returns an error on the first failure, otherwise nil.
func (app *Application) Shutdown(ctx stdContext.Context) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	for i, su := range app.Hosts {
		app.logger.Debugf("Host[%d]: Shutdown now", i)
		if err := su.Shutdown(ctx); err != nil {
			app.logger.Debugf("Host[%d]: Error while trying to shutdown", i)
			return err
		}
	}

	for _, t := range app.config.Tunneling.Tunnels {
		if t.Name == "" {
			continue
		}

		if err := app.config.Tunneling.StopTunnel(t); err != nil {
			return err
		}
	}

	return nil
}

// Build sets up, once, the framework.
// It builds the default router with its default macros
// and the template functions that are very-closed to iris.
//
// If error occurred while building the Application, the returns type of error will be an *errgroup.Group
// which let the callers to inspect the errors and cause, usage:
//
// import "github.com/kataras/iris/core/errgroup"
//
// errgroup.Walk(app.Build(), func(typ interface{}, err error) {
// 	app.Logger().Errorf("%s: %s", typ, err)
// })
func (app *Application) Build() error {
	if app.builded {
		return nil
	}
	// start := time.Now()
	app.builded = true // even if fails.

	// check if a prior app.Logger().SetLevel called and if not
	// then set the defined configuration's log level.
	if app.logger.Level == golog.InfoLevel /* the default level */ {
		app.logger.SetLevel(app.config.LogLevel)
	}

	rp := errgroup.New("Application Builder")
	rp.Err(app.APIBuilder.GetReporter())

	if app.defaultMode { // the app.I18n and app.View will be not available until Build.
		if !app.I18n.Loaded() {
			for _, s := range []string{"./locales/*/*", "./locales/*", "./translations"} {
				if _, err := os.Stat(s); os.IsNotExist(err) {
					continue
				}

				if err := app.I18n.Load(s); err != nil {
					continue
				}

				app.I18n.SetDefault("en-US")
				break
			}
		}

		if app.view.Len() == 0 {
			for _, s := range []string{"./views", "./templates", "./web/views"} {
				if _, err := os.Stat(s); os.IsNotExist(err) {
					continue
				}

				app.RegisterView(HTML(s, ".html"))
				break
			}
		}
	}

	if app.I18n.Loaded() {
		// {{ tr "lang" "key" arg1 arg2 }}
		app.view.AddFunc("tr", app.I18n.Tr)
		app.Router.WrapRouter(app.I18n.Wrapper())
	}

	if n := app.view.Len(); n > 0 {
		tr := "engines"
		if n == 1 {
			tr = tr[0 : len(tr)-1]
		}

		app.logger.Debugf("Application: %d registered view %s", n, tr)
		// view engine
		// here is where we declare the closed-relative framework functions.
		// Each engine has their defaults, i.e yield,render,render_r,partial, params...
		rv := router.NewRoutePathReverser(app.APIBuilder)
		app.view.AddFunc("urlpath", rv.Path)
		// app.view.AddFunc("url", rv.URL)
		if err := app.view.Load(); err != nil {
			rp.Group("View Builder").Err(err)
		}
	}

	if !app.Router.Downgraded() {
		// router
		if _, err := injectLiveReload(app.ContextPool, app.Router); err != nil {
			rp.Errf("LiveReload: init: failed: %v", err)
		}

		if app.config.ForceLowercaseRouting {
			app.Router.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
				r.URL.Path = strings.ToLower(r.URL.Path)
				next(w, r)
			})
		}

		// create the request handler, the default routing handler
		routerHandler := router.NewDefaultHandler(app.config, app.logger)
		err := app.Router.BuildRouter(app.ContextPool, routerHandler, app.APIBuilder, false)
		if err != nil {
			rp.Err(err)
		}
		app.HTTPErrorHandler = routerHandler
		// re-build of the router from outside can be done with
		// app.RefreshRouter()
	}

	// if end := time.Since(start); end.Seconds() > 5 {
	// app.logger.Debugf("Application: build took %s", time.Since(start))

	return errgroup.Check(rp)
}

// Runner is just an interface which accepts the framework instance
// and returns an error.
//
// It can be used to register a custom runner with `Run` in order
// to set the framework's server listen action.
//
// Currently `Runner` is being used to declare the builtin server listeners.
//
// See `Run` for more.
type Runner func(*Application) error

// Listener can be used as an argument for the `Run` method.
// It can start a server with a custom net.Listener via server's `Serve`.
//
// Second argument is optional, it accepts one or more
// `func(*host.Configurator)` that are being executed
// on that specific host that this function will create to start the server.
// Via host configurators you can configure the back-end host supervisor,
// i.e to add events for shutdown, serve or error.
// An example of this use case can be found at:
// https://github.com/kataras/iris/blob/master/_examples/http-server/notify-on-shutdown/main.go
// Look at the `ConfigureHost` too.
//
// See `Run` for more.
func Listener(l net.Listener, hostConfigs ...host.Configurator) Runner {
	return func(app *Application) error {
		app.config.vhost = netutil.ResolveVHost(l.Addr().String())
		return app.NewHost(&http.Server{Addr: l.Addr().String()}).
			Configure(hostConfigs...).
			Serve(l)
	}
}

// Server can be used as an argument for the `Run` method.
// It can start a server with a *http.Server.
//
// Second argument is optional, it accepts one or more
// `func(*host.Configurator)` that are being executed
// on that specific host that this function will create to start the server.
// Via host configurators you can configure the back-end host supervisor,
// i.e to add events for shutdown, serve or error.
// An example of this use case can be found at:
// https://github.com/kataras/iris/blob/master/_examples/http-server/notify-on-shutdown/main.go
// Look at the `ConfigureHost` too.
//
// See `Run` for more.
func Server(srv *http.Server, hostConfigs ...host.Configurator) Runner {
	return func(app *Application) error {
		return app.NewHost(srv).
			Configure(hostConfigs...).
			ListenAndServe()
	}
}

// Addr can be used as an argument for the `Run` method.
// It accepts a host address which is used to build a server
// and a listener which listens on that host and port.
//
// Addr should have the form of [host]:port, i.e localhost:8080 or :8080.
//
// Second argument is optional, it accepts one or more
// `func(*host.Configurator)` that are being executed
// on that specific host that this function will create to start the server.
// Via host configurators you can configure the back-end host supervisor,
// i.e to add events for shutdown, serve or error.
// An example of this use case can be found at:
// https://github.com/kataras/iris/blob/master/_examples/http-server/notify-on-shutdown/main.go
// Look at the `ConfigureHost` too.
//
// See `Run` for more.
func Addr(addr string, hostConfigs ...host.Configurator) Runner {
	return func(app *Application) error {
		return app.NewHost(&http.Server{Addr: addr}).
			Configure(hostConfigs...).
			ListenAndServe()
	}
}

var (
	// TLSNoRedirect is a `host.Configurator` which can be passed as last argument
	// to the `TLS` runner function. It disables the automatic
	// registration of redirection from "http://" to "https://" requests.
	// Applies only to the `TLS` runner.
	// See `AutoTLSNoRedirect` to register a custom fallback server for `AutoTLS` runner.
	TLSNoRedirect = func(su *host.Supervisor) { su.NoRedirect() }
	// AutoTLSNoRedirect is a `host.Configurator`.
	// It registers a fallback HTTP/1.1 server for the `AutoTLS` one.
	// The function accepts the letsencrypt wrapper and it
	// should return a valid instance of http.Server which its handler should be the result
	// of the "acmeHandler" wrapper.
	// Usage:
	//	 getServer := func(acme func(http.Handler) http.Handler) *http.Server {
	//	     srv := &http.Server{Handler: acme(yourCustomHandler), ...otherOptions}
	//	     go srv.ListenAndServe()
	//	     return srv
	//   }
	//   app.Run(iris.AutoTLS(":443", "example.com example2.com", "mail@example.com", getServer))
	//
	// Note that if Server.Handler is nil then the server is automatically ran
	// by the framework and the handler set to automatic redirection, it's still
	// a valid option when the caller wants just to customize the server's fields (except Addr).
	// With this host configurator the caller can customize the server
	// that letsencrypt relies to perform the challenge.
	// LetsEncrypt Certification Manager relies on http://%s:80/.well-known/acme-challenge/<TOKEN>.
	AutoTLSNoRedirect = func(getFallbackServer func(acmeHandler func(fallback http.Handler) http.Handler) *http.Server) host.Configurator {
		return func(su *host.Supervisor) {
			su.NoRedirect()
			su.Fallback = getFallbackServer
		}
	}
)

// TLS can be used as an argument for the `Run` method.
// It will start the Application's secure server.
//
// Use it like you used to use the http.ListenAndServeTLS function.
//
// Addr should have the form of [host]:port, i.e localhost:443 or :443.
// "certFileOrContents" & "keyFileOrContents" should be filenames with their extensions
// or raw contents of the certificate and the private key.
//
// Last argument is optional, it accepts one or more
// `func(*host.Configurator)` that are being executed
// on that specific host that this function will create to start the server.
// Via host configurators you can configure the back-end host supervisor,
// i.e to add events for shutdown, serve or error.
// An example of this use case can be found at:
// https://github.com/kataras/iris/blob/master/_examples/http-server/notify-on-shutdown/main.go
// Look at the `ConfigureHost` too.
//
// See `Run` for more.
func TLS(addr string, certFileOrContents, keyFileOrContents string, hostConfigs ...host.Configurator) Runner {
	return func(app *Application) error {
		return app.NewHost(&http.Server{Addr: addr}).
			Configure(hostConfigs...).
			ListenAndServeTLS(certFileOrContents, keyFileOrContents)
	}
}

// AutoTLS can be used as an argument for the `Run` method.
// It will start the Application's secure server using
// certifications created on the fly by the "autocert" golang/x package,
// so localhost may not be working, use it at "production" machine.
//
// Addr should have the form of [host]:port, i.e mydomain.com:443.
//
// The whitelisted domains are separated by whitespace in "domain" argument,
// i.e "iris-go.com", can be different than "addr".
// If empty, all hosts are currently allowed. This is not recommended,
// as it opens a potential attack where clients connect to a server
// by IP address and pretend to be asking for an incorrect host name.
// Manager will attempt to obtain a certificate for that host, incorrectly,
// eventually reaching the CA's rate limit for certificate requests
// and making it impossible to obtain actual certificates.
//
// For an "e-mail" use a non-public one, letsencrypt needs that for your own security.
//
// Note: `AutoTLS` will start a new server for you
// which will redirect all http versions to their https, including subdomains as well.
//
// Last argument is optional, it accepts one or more
// `func(*host.Configurator)` that are being executed
// on that specific host that this function will create to start the server.
// Via host configurators you can configure the back-end host supervisor,
// i.e to add events for shutdown, serve or error.
// An example of this use case can be found at:
// https://github.com/kataras/iris/blob/master/_examples/http-server/notify-on-shutdown/main.go
// Look at the `ConfigureHost` too.
//
// Usage:
// app.Run(iris.AutoTLS("iris-go.com:443", "iris-go.com www.iris-go.com", "mail@example.com"))
//
// See `Run` and `core/host/Supervisor#ListenAndServeAutoTLS` for more.
func AutoTLS(
	addr string,
	domain string, email string,
	hostConfigs ...host.Configurator) Runner {
	return func(app *Application) error {
		return app.NewHost(&http.Server{Addr: addr}).
			Configure(hostConfigs...).
			ListenAndServeAutoTLS(domain, email, "letscache")
	}
}

// Raw can be used as an argument for the `Run` method.
// It accepts any (listen) function that returns an error,
// this function should be block and return an error
// only when the server exited or a fatal error caused.
//
// With this option you're not limited to the servers
// that iris can run by-default.
//
// See `Run` for more.
func Raw(f func() error) Runner {
	return func(app *Application) error {
		app.logger.Debugf("HTTP Server will start from unknown, external function")
		return f()
	}
}

// ErrServerClosed is returned by the Server's Serve, ServeTLS, ListenAndServe,
// and ListenAndServeTLS methods after a call to Shutdown or Close.
//
// A shortcut for the `http#ErrServerClosed`.
var ErrServerClosed = http.ErrServerClosed

// Listen builds the application and starts the server
// on the TCP network address "host:port" which
// handles requests on incoming connections.
//
// Listen always returns a non-nil error.
// Ignore specific errors by using an `iris.WithoutServerError(iris.ErrServerClosed)`
// as a second input argument.
//
// Listen is a shortcut of `app.Run(iris.Addr(hostPort, withOrWithout...))`.
// See `Run` for details.
func (app *Application) Listen(hostPort string, withOrWithout ...Configurator) error {
	return app.Run(Addr(hostPort), withOrWithout...)
}

// Run builds the framework and starts the desired `Runner` with or without configuration edits.
//
// Run should be called only once per Application instance, it blocks like http.Server.
//
// If more than one server needed to run on the same iris instance
// then create a new host and run it manually by `go NewHost(*http.Server).Serve/ListenAndServe` etc...
// or use an already created host:
// h := NewHost(*http.Server)
// Run(Raw(h.ListenAndServe), WithCharset("utf-8"), WithRemoteAddrHeader("CF-Connecting-IP"))
//
// The Application can go online with any type of server or iris's host with the help of
// the following runners:
// `Listener`, `Server`, `Addr`, `TLS`, `AutoTLS` and `Raw`.
func (app *Application) Run(serve Runner, withOrWithout ...Configurator) error {
	app.Configure(withOrWithout...)

	if err := app.Build(); err != nil {
		app.logger.Error(err)
		return err
	}

	app.ConfigureHost(func(host *Supervisor) {
		host.SocketSharding = app.config.SocketSharding
	})

	app.tryStartTunneling()

	if len(app.Hosts) > 0 {
		app.logger.Debugf("Application: running using %d host(s)", len(app.Hosts)+1 /* +1 the current */)
	}

	// this will block until an error(unless supervisor's DeferFlow called from a Task).
	err := serve(app)
	if err != nil {
		app.logger.Error(err)
	}

	return err
}

// https://ngrok.com/docs
func (app *Application) tryStartTunneling() {
	if len(app.config.Tunneling.Tunnels) == 0 {
		return
	}

	app.ConfigureHost(func(su *host.Supervisor) {
		su.RegisterOnServe(func(h host.TaskHost) {
			publicAddrs, err := tunnel.Start(app.config.Tunneling)
			if err != nil {
				app.logger.Errorf("Host: tunneling error: %v", err)
				return
			}

			publicAddr := publicAddrs[0]
			// to make subdomains resolution still based on this new remote, public addresses.
			app.config.vhost = publicAddr[strings.Index(publicAddr, "://")+3:]

			directLog := []byte(fmt.Sprintf("??? Public Address: %s\n", publicAddr))
			app.logger.Printer.Write(directLog) // nolint:errcheck
		})
	})
}
