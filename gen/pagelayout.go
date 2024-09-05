package gen

var (
	imgBannerHeight   = "10em"
	emptyBannerHeight = "2em"

	defaultLayout = `
<!DOCTYPE html>
<html>
<head>
<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
<meta content="utf-8" http-equiv="encoding">
<title></title>
<style>
@media (min-width: 1080px) {
	.row {
		display: grid;
		grid-template-columns: 1fr 2fr 1fr;
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
	.row .header {
		padding-bottom: 4em;
	}
	.row .navi {
		padding-left: 2em;
		padding-right: 2em;
	}
	.row .main {
		padding-left: 2em;
		padding-right: 2em;
		border-left: 1px solid {{ .Palette.LightGray }};
		border-right: 1px solid {{ .Palette.LightGray }};
  		height: 100%;
		min-height: 50em;
	}
	.row .main-left {
	}
	.row .footer {
		margin-top: 4em;
	}
	{{ .Content.Header.Css }}
	{{ .Content.Navi.Css }}
	{{ .Content.MainLeft.Css }}
	{{ .Content.Main.Css }}
	{{ .Content.Footer.Css }}
</style>
</head>
<body>
	<script>
		{{ .Content.Header.Js }}
		{{ .Content.Navi.Js }}
		{{ .Content.MainLeft.Js }}
		{{ .Content.Main.Js }}
		{{ .Content.Footer.Js }}
	</script>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="header">
			{{ .Content.Header.Html }}
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="navi">
			{{ .Content.Navi.Html }}
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="main-left">
			{{ .Content.MainLeft.Html }}
		</div>
		<div class="main">
			{{ .Content.Main.Html }}
		</div>
		<div class="main-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="footer">
			{{ .Content.Footer.Html }}
		</div>
		<div class="col-right">
		</div>
	</div>
</body>
</html>
`
)
