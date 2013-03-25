package templates

const documentTemplate = `<!DOCTYPE HTML>
<html lang="{{.LanguageTag}}">
<head>
	<title>{{.Title}}</title>
</head>

<body>

<article>

<header>
<h1>
{{.Title}}
</h1>
</header>

<section>
{{.Description}}
</section>

<section>
{{.Content}}
</section>

</article>

</body>
</html>`
