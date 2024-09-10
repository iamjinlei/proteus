package gen

var (
	imgBannerHeight = "10em"

	defaultLayout = `
<!DOCTYPE html>
<html>
<head>
{{ if .CanonicalDomain }}
<link rel="canonical" href="{{ .CanonicalDomain }}/{{ .RelPath }}"/>
{{ end }}
<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
<meta content="utf-8" http-equiv="encoding">
<title></title>
<style>
@media (min-width: 1080px) {
	.row {
		display: grid;
		grid-template-columns: 1fr 960px 1fr;
		font-size: 1.2em;
	}
}
@media (max-width: 1079px) {
	.row {
		display: grid;
		grid-template-columns: 1fr 4fr 1fr;
		font-size: 60%;
	}
}
	.col-middle {
		padding-left: 1em;
		padding-right: 1em;
	}
	.row.header.nonempty {
		padding-bottom: 1em;
	}
	.row.header.empty {
		min-height: 6em;
	}
	.row.nav.nonempty {
	}
	.row.nav.empty {
	}
	.row.main {
  		height: 100%;
		min-height: 50em;
		margin-bottom: 5em;
		padding: 1em;
	}
	.row.main .col-middle {
		border: 1px solid {{ .Palette.LightGray }};
	}
	.row.footer {
	}
	{{ .Content.Header.Css }}
	{{ .Content.Nav.Css }}
	{{ .Content.MainLeft.Css }}
	{{ .Content.MainRight.Css }}
	{{ .Content.Main.Css }}
	{{ .Content.Footer.Css }}
</style>
</head>
<body>
	<script>
		{{ .Content.Header.Js }}
		{{ .Content.Nav.Js }}
		{{ .Content.MainLeft.Js }}
		{{ .Content.MainRight.Js }}
		{{ .Content.Main.Js }}
		{{ .Content.Footer.Js }}
	</script>

	{{ if .Content.Header.Html }}
	<div class="row header nonempty">
		<div class="col-left">
		</div>
		<div class="col-middle">
			{{ .Content.Header.Html }}
		</div>
		<div class="col-right">
		</div>
	</div>
	{{ else }}
	<div class="row header empty">
	</div>
	{{ end }}

	{{ if .Content.Nav.Html }}
	<div class="row nav nonempty">
		<div class="col-left">
		</div>
		<div class="col-middle">
			{{ .Content.Nav.Html }}
		</div>
		<div class="col-right">
		</div>
	</div>
	{{ else }}
	<div class="row nav empty">
	</div>
	{{ end }}

	<div class="row main">
		<div class="col-left">
			{{ .Content.MainLeft.Html }}
		</div>
		<div class="col-middle">
			{{ .Content.Main.Html }}
		</div>
		<div class="col-right">
			{{ .Content.MainRight.Html }}
		</div>
	</div>

	<div class="row footer">
		<div class="col-left">
		</div>
		<div class="col-middle">
			{{ .Content.Footer.Html }}
		</div>
		<div class="col-right">
		</div>
	</div>
</body>
</html>
`
)
