package gen

var (
	centerColWidth    = "60em"
	imgBannerHeight   = "10em"
	emptyBannerHeight = "2em"

	defaultLayout = `
<!DOCTYPE html>
<html>
<head>
<title></title>
<style>
	.row {
		display: grid;
		grid-template-columns: 1fr {{ .Dimensions.CenterColWidth }} 1fr;
		font-size: 120%;
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
	.row .footer {
		margin-top: 4em;
	}
</style>
</head>
<body>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="header">
			{{ .Content.Header }}
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="navi">
			{{ .Content.Navi }}
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="main">
			{{ .Content.Main }}
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="footer">
			{{ .Content.Footer }}
		</div>
		<div class="col-right">
		</div>
	</div>
</body>
</html>
`
)
