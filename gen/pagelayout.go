package gen

import "bytes"

var (
	defaultLayout = []byte(fillColors(`
<!DOCTYPE html>
<html>
<head>
<title></title>
<style>
	.row {
		display: grid;
		grid-template-columns: 1fr 60em 1fr;
		font-size: 120%;
	}
	.row .top-center {
	}
	.row .mid-center {
		padding-top: 2em;
		padding-left: 2em;
		padding-right: 2em;
		border-left: 1px solid #LightGray;
		border-right: 1px solid #LightGray;
		border-bottom: 1px solid #LightGray;
  		height: 100%;
		min-height: 50em;
	}
	.row .bottom-center {
	}
</style>
</head>
<body>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="top-center">
			$$__HEADER__$$
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="mid-center">
			$$__CONTENT__$$
		</div>
		<div class="col-right">
		</div>
	</div>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="bottom-center">
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
	content []byte,
	footer []byte,
) []byte {
	p := bytes.Replace(defaultLayout, []byte("$$__HEADER__$$"), header, 1)
	p = bytes.Replace(p, []byte("$$__CONTENT__$$"), content, 1)
	p = bytes.Replace(p, []byte("$$__FOOTER__$$"), footer, 1)
	return p
}
