package gen

import "bytes"

var (
	centerColWidth    = "60em"
	imgBannerHeight   = "10em"
	emptyBannerHeight = "2em"

	defaultLayout = []byte(fillColors(`
<!DOCTYPE html>
<html>
<head>
<title></title>
<style>
	.row {
		display: grid;
		grid-template-columns: 1fr $$__CENTER_COL_WIDTH__$$ 1fr;
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
		border-left: 1px solid #LightGray;
		border-right: 1px solid #LightGray;
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
			$$__HEADER__$$
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="navi">
			$$__NAVI__$$
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="main">
			$$__MAIN__$$
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="footer">
			$$__FOOTER__$$
		</div>
		<div class="col-right">
		</div>
	</div>
</body>
</html>
`))
)

func fillPageTemplate(
	header []byte,
	navi []byte,
	main []byte,
	footer []byte,
) []byte {
	p := bytes.Replace(defaultLayout, []byte("$$__HEADER__$$"), header, 1)
	p = bytes.Replace(p, []byte("$$__NAVI__$$"), navi, 1)
	p = bytes.Replace(p, []byte("$$__MAIN__$$"), main, 1)
	p = bytes.Replace(p, []byte("$$__FOOTER__$$"), footer, 1)
	p = bytes.Replace(p, []byte("$$__CENTER_COL_WIDTH__$$"), []byte(centerColWidth), 1)
	return p
}
