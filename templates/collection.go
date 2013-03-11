package templates

const CollectionTemplate = `<!DOCTYPE HTML>
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

<section>
<ul>
{{range .Entries}}
<li>
	<a href="{{.Path}}">{{.Title}}</a>
	<p>{{.Description}}</p>
</li>
{{end}}
</ul>
</section>

</article>

</body>
</html>`
