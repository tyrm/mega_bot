package web

var HeadFrameworkCSSTemplate = []templateHeadLink{
	{
		HRef: "/static/vendor/bootstrap-4.5.3-cerulean.min.css",
		Rel: "stylesheet",
		Integrity: "sha384-mUmCmzW7iDM1B23t+9wTAPdZzabcxrgWTtz7iHs+6QOPBl6rMS3fHSNqONSgMaJI",
		CrossOrigin: "anonymous",
	},
	{
		HRef: "/static/vendor/fontawesome-free-5.15.1-web/css/all.min.css",
		Rel: "stylesheet",
		Integrity: "sha384-vp86vTRFVJgpjF9jiIGPEEqYqlDwgyBgEF109VFjmqGmIY/Y4HV4d3Gp2irVfcrp",
		CrossOrigin: "anonymous",
	},
}

var HeadFaviconsTemplate = []templateHeadLink{
	{
		HRef: "/static/img/favicons/apple-touch-icon.png",
		Rel: "apple-touch-icon",
		Sizes: "180x180",
	},
	{
		HRef: "/static/img/favicons/favicon.svg",
		Rel: "icon",
		Type: "image/svg+xml",
	},
}

var HeadCSSTemplate = []templateHeadLink{
	{
		HRef: "/static/css/default.css",
		Rel: "stylesheet",
	},
}