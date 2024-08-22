package gen

var (
	placeHolder   = []byte("$$__CONTENT__$$")
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
	.row [class*="col"] {
	}
	.row .col-right {
	}
	.row .col-center {
		padding-left: 2em;
		padding-right: 2em;
		border-left: 1px solid #LightGray;
		border-right: 1px solid #LightGray;
  		height: 100%;
	}
</style>
</head>
<body>
	<div class="row">
		<div class="col-left">
		</div>
		<div class="col-center">
			$$__CONTENT__$$
		</div>
		<div class="col-right">
		</div>
	</div>
</body>
</html>
`))
)
